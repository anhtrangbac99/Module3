package main

import (
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageBill"
	grpc "google.golang.org/grpc"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
	"sync"
	"net"
	handler "git.zapa.cloud/fresher/kietcdx/Module3/App/Core/Merchant/Manage_Bill/Handler"

)

func main(){

	grpcServer := grpc.NewServer()
	pb.RegisterManageBillServer(grpcServer,new(handler.ManageBillServer))

	fmt.Println("Listening to port 6001")
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func(){
		defer Wg.Done()
		
		lis,err := net.Listen("tcp",":6001")
		if err != nil {
			log.Fatal(err)
		}

		if err:=grpcServer.Serve(lis);err != nil {
			log.Fatal(err)
		}
	}()
	Wg.Wait()

}