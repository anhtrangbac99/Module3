package main

import (
	handler "git.zapa.cloud/fresher/kietcdx/Module3/App/Core/Merchant/Manage_Support/Handler"
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageSupport"
	//"github.com/jinzhu/copier"
	"sync"
	"log"
	"fmt"
	"net"
	grpc "google.golang.org/grpc"
)

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterManageSupportServer(grpcServer,new(handler.ManageSupportServer))

	fmt.Println("Listening to port 6002")
	Wg := sync.WaitGroup{}

	Wg.Add(1)
	go func(){
		defer Wg.Done()
		
		lis,err := net.Listen("tcp",":6002")
		if err != nil {
			log.Fatal(err)
		}

		if err:=grpcServer.Serve(lis);err != nil {
			log.Fatal(err)
		}
	}()
	Wg.Wait()

}
