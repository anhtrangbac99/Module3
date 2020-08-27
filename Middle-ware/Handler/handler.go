package handler

import (
	pb "google.golang.org/grpc/examples/App/Proto/Middle-ware"
	"context"
	grpc "google.golang.org/grpc"
	"fmt"
	"log"

	"io"
	pbSignIn "google.golang.org/grpc/examples/App/Proto/SignIn"
	pbManageBill "google.golang.org/grpc/examples/App/Proto/ManageBill"
	pbManageSupport "google.golang.org/grpc/examples/App/Proto/ManageSupport"
	"database/sql"
	"github.com/jinzhu/copier"
	"github.com/golang/glog"
)
const merchantManageAcountCoreHost = "127.0.0.1:6000"
const merchantManageBillCoreHost = "127.0.0.1:6001"
const merchantManageSupportCoreHost = "127.0.0.1:6002"

type merchantHandler struct {
	MerchantManageAccountCoreService pbSignIn.SignInClient
	MerchantManageBillCoreService pbManageBill.ManageBillClient
	MerchantManageSupportCoreService pbManageSupport.ManageSupportClient
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

	connManageSupport,err := grpc.Dial(merchantManageSupportCoreHost,grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	clientSignIn := pbSignIn.NewSignInClient(connManageAccount)

	clientManageBill := pbManageBill.NewManageBillClient(connManageBill)

	clientManageSupport := pbManageSupport.NewManageSupportClient(connManageSupport)
	return &merchantHandler{
		MerchantManageAccountCoreService: clientSignIn,
		MerchantManageBillCoreService: clientManageBill,
		MerchantManageSupportCoreService: clientManageSupport,
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
	glog.Info("UserAuthor request:",request)
	client := sv.MerchantManageAccountCoreService

	reqToCore := pbSignIn.AuthorRequest{}

	if err:= copier.Copy(&reqToCore,request);err!= nil {
		log.Fatal(err)
	}
	//fmt.Println(reqToCore)
	result,err := client.UserAuthor(context.Background(),&pbSignIn.AuthorRequest{Username: request.GetUsername() , Password: request.GetPassword()})

	if err != nil {
		log.Fatal(err)
	}
	//respone,err := json.Marshal(&result)

	//writer.Write(respone)
	respone := pb.AuthorRespone{IsExisted:result.GetIsExisted(),User_Id:result.GetUser_Id(),Authorized:result.GetAuthorized()}
	return &respone,nil
}

func (sv *merchantHandler) CreateBill(ctx context.Context, request *pb.CreateBillRequest) (*pb.CreateBillRespone,error){
	// connect,err := grpc.Dial("localhost:6001",grpc.WithInsecure(),grpc.WithBlock())
	// defer connect.Close()

	// if err != nil {
	// 	log.Fatal(err)
	// }
	glog.Info("Create Bill request: ", request)

	client := sv.MerchantManageBillCoreService
	
	reqToCore := pbManageBill.CreateBillRequest{}

	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	result,_ := client.CreateBill(context.Background(),&reqToCore)
	fmt.Println(result.GetIsSaved())

	respone := pb.CreateBillRespone{IsSaved:result.GetIsSaved()}
	return &respone,nil
}

func  (sv *merchantHandler) SearchBill(ctx context.Context,request *pb.SearchBillRequest) (*pb.ListSearchBillRespone,error){

	glog.Info("SearchBill request: ", request)

	reqToCore := pbManageBill.SearchBillRequest{}

	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	client := sv.MerchantManageBillCoreService

	stream,_ := client.SearchBill(context.Background(),&reqToCore)
	
	listRespone:= []*pb.SearchBillRespone{}
	for {
		result,err := stream.Recv()
		respone := pb.SearchBillRespone{}
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if err := copier.Copy(&respone,result);err != nil {
			log.Fatal(err)
		}
		listRespone = append(listRespone,&respone)
		fmt.Println(respone)
	}

	return &pb.ListSearchBillRespone{SearchBillRespones:listRespone},nil
}

func (sv *merchantHandler) GetListItem(ctx context.Context,request *pb.ListItemRequest) (*pb.ListItemRespone,error){
	glog.Info("Get List Item request: " ,request)

	client := sv.MerchantManageSupportCoreService

	stream,_ := client.GetListItem(context.Background(),&pbManageSupport.ListItemRequest{})

	ListItem := []*pb.Item{}
	for {
		result,err := stream.Recv()
		item := pb.Item{}
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if err := copier.Copy(&item,result);err != nil {
			log.Fatal(err)
		}
		ListItem = append(ListItem,&item)
		fmt.Println(item)
	}

	return &pb.ListItemRespone{Item:ListItem},nil

}

func (sv *merchantHandler) GetCustomer(ctx context.Context,request *pb.CustomerRequest) (*pb.CustomerRespone,error){
	glog.Info("Get Cusomter request",request)

	reqToCore := pbManageSupport.CustomerRequest{}

	if err := copier.Copy(&reqToCore,request);err!=nil{
		log.Fatal(err)
	}

	client := sv.MerchantManageSupportCoreService

	result,_ := client.GetCustomer(context.Background(),&reqToCore) 

	respone := pb.CustomerRespone{}

	if err := copier.Copy(&respone,result);err!=nil{
		log.Fatal(err)
	}
	fmt.Println()
	return &respone,nil
}