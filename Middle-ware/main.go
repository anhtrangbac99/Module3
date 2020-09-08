package main

import (
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Middle-ware"
	"net"
	"log"
	handler "git.zapa.cloud/fresher/kietcdx/Module3/App/Middle-ware/Handler"
	"google.golang.org/grpc/reflection"
	grpc "google.golang.org/grpc"
	"context"
	"reflect"
	"github.com/go-redis/redis/v8" 
	"time"
)

const (
	port = ":1234"
)

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,handlerGrpc grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/api.merchant.MerchantMiddlewareService/CheckUserToken"{
		log.Println("--> unary interceptor: ", info.FullMethod)
		val := reflect.ValueOf(req).Elem()
		request := val.FieldByName("UserToken").Interface().(string)
		redisClient := redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			Password: "",
			DB: 0,
		})
		DB := handler.AccessDB()
	
		value,err := redisClient.Get(ctx,request).Result()
	
		if err == redis.Nil {
			return &pb.UserTokenRespone{IsExisted:int64(-1)},nil
		} else {
			if err != nil {
				log.Fatal(err)
			}
			
			isExpire := redisClient.Expire(ctx,request,15*time.Minute)
			log.Println("Extend expiration time of token ",isExpire)
			author,err := DB.Query(`SELECT Authorized FROM User WHERE User_Id="` + value + `";`)
	
			if err != nil {
				log.Fatal(err)
			}
	
			var Author handler.Authorized
	
			for author.Next(){
	
				err := author.Scan(&Author.Authorized)
				if err != nil {
					log.Fatal(err)
				}
			}
			return &pb.UserTokenRespone{IsExisted:int64(1),Authorized:int64(Author.Authorized)},nil
		}
	}
    return handlerGrpc(ctx, req)
}

func main()  {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("server run at: %v", port)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterMerchantMiddlewareServiceServer(s,handler.NewMerchantMiddlewareHanlder())
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
