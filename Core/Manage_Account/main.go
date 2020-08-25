package main

import (
	pb "google.golang.org/grpc/examples/App/Proto/SignIn"
	grpc "google.golang.org/grpc"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"context"
	"net"
	"fmt"
	"sync"
)

type selectedUser struct {
	User_Id int
	Username string 
	Password string
	Authorized int
}
type SignInServer struct{
	pb.UnimplementedSignInServer
}

func AccessDB() (*sql.DB){
	dbDriver := "mysql"
 
    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	
	return db
}


func (signInServer *SignInServer) UserAuthor(ctx context.Context,request *pb.AuthorRequest) (*pb.AuthorRespone,error){
	DB := AccessDB()
	fmt.Println(request.GetUsername())
	selectedUsers, err := DB.Query("SELECT User_Id,Username,Password,Authorized FROM User WHERE Username='" + request.GetUsername() +"';")
	
	isExisted := -1
	user_Id := -1
	authorized := -1

	if err != nil {
		log.Fatal(err)
	} else {
		for selectedUsers.Next(){
			var selectedUser selectedUser
			err = selectedUsers.Scan(&selectedUser.User_Id,&selectedUser.Username,&selectedUser.Password,&selectedUser.Authorized)
			//fmt.Println(selectedUser)
			if err != nil {
				log.Fatal(err)
			}

			if selectedUser.Password == request.GetPassword() {
				isExisted = 1
				user_Id = selectedUser.User_Id
				authorized = selectedUser.Authorized
				break;
			} else {
				isExisted = 1
			}
		}
	}

	return &pb.AuthorRespone{IsExisted:int64(isExisted),User_Id:int64(user_Id),Authorized:int64(authorized)},nil
}

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterSignInServer(grpcServer,new(SignInServer))

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