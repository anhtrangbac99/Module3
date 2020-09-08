package handler

import (
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Middle-ware"
	"context"
	grpc "google.golang.org/grpc"
	"fmt"
	"log"
	"io"
	pbSignIn "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/SignIn"
	pbManageBill "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageBill"
	pbManageSupport "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageSupport"
	"database/sql"
	"github.com/jinzhu/copier"
	"github.com/golang/glog"
	_ "github.com/go-sql-driver/mysql"

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

type Authorized struct {
	Authorized int `json:"Authorized"`
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
	fmt.Println(result.GetUser_Id())
	respone := pb.AuthorRespone{IsExisted:result.GetIsExisted(),User_Id:result.GetUser_Id(),Authorized:result.GetAuthorized()}
	return &respone,nil
}

func (sv *merchantHandler) Search(ctx context.Context,request *pb.SearchRequest) (*pb.SearchRespone,error){
	glog.Info("Search request",request)


	reqToCore := pbManageBill.SearchBillRequest{}

	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	//glog.Info("Request to core",reqToCore.GetBillId())
	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	client := sv.MerchantManageBillCoreService

	stream,_ := client.SearchBill(context.Background(),&reqToCore)
	
	listRespone:= []*pb.SearchBillRespone{}
	
	existed := false
	indexExisted := -1
	for {
		result,err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		existed = false
		for index,value := range listRespone{
			if value.GetBillId()==result.GetBillId(){
				existed = true
				indexExisted = index
				break
			}
		}

		if existed {
			listRespone[indexExisted] = &pb.SearchBillRespone{BillId:listRespone[indexExisted].GetBillId(),BillStatus:listRespone[indexExisted].GetBillStatus(),CustomerId:listRespone[indexExisted].GetCustomerId(),CustomerPhone:listRespone[indexExisted].GetCustomerPhone(),BillDesc:listRespone[indexExisted].GetBillDesc(),CustomerName:listRespone[indexExisted].GetCustomerName(),Item:
			append(listRespone[indexExisted].GetItem(),&pb.ListItem{ItemId:result.GetItemId(),ItemName:result.GetItemName(),Amount:result.GetAmount()})}
		} else {
			respone := pb.SearchBillRespone{BillId:result.GetBillId(),BillStatus:result.GetBillStatus(),CustomerId:result.GetCustomerId(),CustomerPhone:result.GetCustomerPhone(),BillDesc:result.GetBillDesc(),CustomerName:result.GetCustomerName(),Item:[]*pb.ListItem{&pb.ListItem{ItemId:result.GetItemId(),ItemName:result.GetItemName(),Amount:result.GetAmount()}}}
		

			if err := copier.Copy(&respone,result);err != nil {
				log.Fatal(err)
			}
			listRespone = append(listRespone,&respone)
		} 
	}

	glog.Info("Search Respone ",listRespone)
	return &pb.SearchRespone{SearchRespones:listRespone},nil

	//return &pb.SearchRespone{},nil
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

func (sv *merchantHandler) ListItem(ctx context.Context,request *pb.ListItemRequest) (*pb.ListItemRespone,error){
	glog.Info("Get List Item request: " ,request)

	client := sv.MerchantManageSupportCoreService

	stream,_ := client.ListItem(context.Background(),&pbManageSupport.ListItemRequest{})

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


func (sv *merchantHandler) CheckUserToken(ctx context.Context, request *pb.UserTokenRequest) (*pb.UserTokenRespone,error){
	glog.Info("Check Token request",request)

	reqToCore := pbManageSupport.UserTokenRequest{}

	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	result,err := sv.MerchantManageSupportCoreService.CheckUserToken(context.Background(),&reqToCore)

	if err != nil {
		log.Fatal(err)
	}

	respone := pb.UserTokenRespone{}

	if err:=copier.Copy(&respone,result);err != nil{
		log.Fatal(err)
	}

	return &respone,nil
}

func (sv *merchantHandler) BillDetail(ctx context.Context,request *pb.BillDetailRequest) (*pb.BillDetailRespone,error){
	glog.Info("Bill Detail request",request)

	reqToCore := pbManageSupport.BillDetailRequest{}

	if err:=copier.Copy(&reqToCore,request);err != nil{
		log.Fatal(err)
	}

	result,err := sv.MerchantManageSupportCoreService.BillDetail(context.Background(),&reqToCore)

	if err != nil {
		log.Fatal(err)
	}

	respone := pb.BillDetailRespone{}

	if err:=copier.Copy(&respone,result);err != nil{
		log.Fatal(err)
	}

	return &respone,nil


}