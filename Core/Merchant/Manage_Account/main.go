package main

import (
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/SignIn"
	grpc "google.golang.org/grpc"
	_ "github.com/go-sql-driver/mysql"
	handler "git.zapa.cloud/fresher/kietcdx/Module3/App/Core/Merchant/Manage_Account/Handler"
	"log"
	"net"
	"fmt"
	"sync"
)

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterSignInServer(grpcServer,new(handler.SignInServer))

	fmt.Println("Listening to port 6000")
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func(){
		defer Wg.Done()
		
		lis,err := net.Listen("tcp",":6000")
		
		if err != nil {
			log.Fatal(err)
		}

		if err:=grpcServer.Serve(lis);err != nil {
			log.Fatal(err)
		}
	}()
	Wg.Wait()
}