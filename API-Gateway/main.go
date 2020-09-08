package main

import (
	"flag"
	"net/http"
	"strings"
	"fmt"
	"google.golang.org/grpc"
	"path"
	gwMerchant "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Middle-ware"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	//"github.com/go-redis/redis/v8" 

)

const (
	merchantEndpoint = "127.0.0.1:1234"
)
type User struct{
	Username string 
	Password string
}

var (
	merchantEnpoint = flag.String("merchant", merchantEndpoint, "endpoint of merchant")
	swaggerDir      = flag.String("swagger", "swagger", "path to the directory which contains swagger definitions")
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {

	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(*swaggerDir, p)
	http.ServeFile(w, r, p)
}


func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	//register merchantApi
	err := gwMerchant.RegisterMerchantMiddlewareServiceHandlerFromEndpoint(ctx, mux, *merchantEnpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE","PATCH","OPTIONS"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	fmt.Println(r)
	return
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			//r.Method == "OPTIONS" && 
			if r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

func Run(address string, opts ...runtime.ServeMuxOption) error {
	fmt.Println("Listenning to port "+address)
	glog.Info("Run")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//mux := runtime.NewServeMux()
	mux := http.NewServeMux()
	//router := mux.NewRouter()
	mux.HandleFunc("/swagger/", serveSwagger)

	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	return http.ListenAndServe(address, allowCORS(mux))
}

func main() {

	if err := Run(":8082"); err != nil {
		glog.Fatal(err)
	}
}
