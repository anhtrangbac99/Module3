package handler

import (
	pb "google.golang.org/grpc/examples/App/Proto/Middle-ware"
	"context"
	grpc "google.golang.org/grpc"
	"fmt"
	"log"
	//"glog"
	"io"
	pbSignIn "google.golang.org/grpc/examples/App/Proto/SignIn"
	pbManageBill "google.golang.org/grpc/examples/App/Proto/ManageBill"
	"database/sql"
	"github.com/jinzhu/copier"

)
const merchantManageAcountCoreHost = "127.0.0.1:6000"
const merchantManageBillCoreHost = "127.0.0.1:6001"
type merchantHandler struct {
	MerchantManageAccountCoreService pbSignIn.SignInClient
	MerchantManageBillCoreService pbManageBill.ManageBillClient
	pb.UnimplementedMerchantMiddlewareServiceServer
}

func NewMerchantMiddlewareHanlder() pb.MerchantMiddlewareServiceServer{
	connManageAccount,err := grpc.Dial(merchantManageAcountCoreHost,grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	connManageBill,err := grpc.Dial(merchantManageBillCoreHost,grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	clientSignIn := pbSignIn.NewSignInClient(connManageAccount)

	clientManageBill := pbManageBill.NewManageBillClient(connManageBill)

	return &merchantHandler{
		MerchantManageAccountCoreService: clientSignIn,
		MerchantManageBillCoreService: clientManageBill,
	}
}


func AccessDB() (*sql.DB){
	dbDriver := "mysql"

    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	
	return db
}

func (sv *merchantHandler) UserAuthor(ctx context.Context, request *pb.AuthorRequest) (*pb.AuthorRespone,error){
	client := sv.MerchantManageAccountCoreService

	reqToCore := pbSignIn.AuthorRequest{}

	if err:= copier.Copy(&reqToCore,request);err!= nil {
		log.Fatal(err)
	}
	fmt.Printf("request: %+v",request)
	//fmt.Println(reqToCore)
	result,err := client.UserAuthor(context.Background(),&pbSignIn.AuthorRequest{Username: request.GetUsername() , Password: request.GetPassword()})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.GetUser_Id())
	fmt.Println(result)
	//respone,err := json.Marshal(&result)

	//writer.Write(respone)
	respone := pb.AuthorRespone{IsExisted:result.GetIsExisted(),User_Id:result.GetUser_Id(),Authorized:result.GetAuthorized()}
	return &respone,nil
}

func (sv *merchantHandler) CreateBill(ctx context.Context, request *pb.CreateBillRequest) (*pb.CreateBillRespone,error){
	connect,err := grpc.Dial("localhost:6001",grpc.WithInsecure(),grpc.WithBlock())
	defer connect.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := pbManageBill.NewManageBillClient(connect)

	result,err := client.CreateBill(context.Background(),&pbManageBill.CreateBillRequest{ItemId:int64(1),Amount:int64(2),CustomerId:int64(2),BillDesc:"None"})
	fmt.Println(result.GetIsSaved())
	respone := pb.CreateBillRespone{IsSaved:result.GetIsSaved()}
	return &respone,nil
}

func  (sv *merchantHandler) SearchBill(request *pb.SearchBillRequest,stream_api pb.MerchantMiddlewareService_SearchBillServer) error{
	connect,err := grpc.Dial("localhost:6001",grpc.WithInsecure(),grpc.WithBlock())
	defer connect.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := pbManageBill.NewManageBillClient(connect)
	temp := pbManageBill.SearchBillRequest{BillDesc:"hhhhh"}

	stream,_ := client.SearchBill(context.Background(),&temp)
	
	for {
		result,err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}

	return nil
}