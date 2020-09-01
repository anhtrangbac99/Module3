package main

import (
	"net"
	"log"
	handler "google.golang.org/grpc/examples/App/Middle-ware/Handler"
	"google.golang.org/grpc/reflection"
	grpc "google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/App/Proto/Middle-ware"

)

const (
	port = ":1234"
)

func main()  {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("server run at: %v", port)
	}
	s := grpc.NewServer()
	pb.RegisterMerchantMiddlewareServiceServer(s,handler.NewMerchantMiddlewareHanlder())
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
