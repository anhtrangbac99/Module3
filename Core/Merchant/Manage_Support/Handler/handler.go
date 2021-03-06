package handler

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"context"
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageSupport"
	//"github.com/go-redis/redis/v8" 
	"strconv" 
)

type ManageSupportServer struct{
	pb.UnimplementedManageSupportServer
}

type Item struct {
	ItemName string `json:"ItemName"`
	ItemId int `json:"ItemId"`
}

type Customer struct {
	CustomerName string `json:"CustomerName"`
	CustomerId int `json:"CustomerId"`
	CustomerPhone string `json:"CustomerPhone"`
}

type Authorized struct {
	Authorized int `json:"Authorized"`
}

type Bill struct {
	ItemName string `json:"ItemName"`
	Amount int `json:"Amount"`
	Price int `json:"Price"`
}
func AccessDB() (*sql.DB){
	dbDriver := "mysql"

    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	
	return db
}

func LogError(err error){
	if err !=nil {
		log.Fatal(err)
	}
}


func (sv *ManageSupportServer) ListItem(request *pb.ListItemRequest,stream pb.ManageSupport_ListItemServer) error {
	DB := AccessDB()

	Query :=`SELECT Item_Name,Item_Id FROM Item;`

	listItem,err := DB.Query(Query)

	LogError(err)
	var ListItem []Item
	for listItem.Next(){
		var Item Item

		err := listItem.Scan(&Item.ItemName,&Item.ItemId)

		LogError(err)

		ListItem = append(ListItem,Item)
	}

	for _,value := range ListItem{
		if err:=stream.Send(&pb.ListItemRespone{ItemName:value.ItemName,ItemId:int64(value.ItemId)});err!=nil{
			LogError(err)
		}
	}

	return nil
}

func (sv *ManageSupportServer) GetCustomer(ctx context.Context, request *pb.CustomerRequest) (*pb.CustomerRespone,error){
	DB := AccessDB()

	Query := `SELECT Customer_Name,Customer_Id,Customer_Phone FROM Customer WHERE Customer_Phone="` + request.GetCustomerPhone() + `";`

	customer,err := DB.Query(Query)

	LogError(err)
	var Customer Customer

	for customer.Next(){

		err := customer.Scan(&Customer.CustomerName,&Customer.CustomerId,&Customer.CustomerPhone)

		LogError(err)
	}

	return &pb.CustomerRespone{CustomerName:Customer.CustomerName,CustomerId:int64(Customer.CustomerId),CustomerPhone:Customer.CustomerPhone},nil 
}

func (sv *ManageSupportServer) BillDetail(ctx context.Context,request *pb.BillDetailRequest) (*pb.BillDetailRespone,error){
	DB := AccessDB()

	Query := `SELECT i.Item_Name,bd.Amount,i.Price
	FROM BillDetail as bd
	LEFT JOIN Item as i ON bd.Item_Id=i.Item_Id
	WHERE bd.Bill_Id="
	` + strconv.Itoa(int(request.GetBillId())) + `";`

	billDB,err := DB.Query(Query)

	if err != nil {
		log.Fatal(err)
	}

	bills := []*pb.Item{}
	total := 0
	for billDB.Next(){
		var bill Bill
		err = billDB.Scan(&bill.ItemName,&bill.Amount,&bill.Price)
		
		if err != nil{
			log.Fatal(err)
		}
		total = total + bill.Amount*bill.Price
		bills = append(bills,&pb.Item{ItemName:bill.ItemName,Amount:int64(bill.Amount),Price:int64(bill.Price)})
	}


	return &pb.BillDetailRespone{ListItem:bills,Total:int64(total)},nil
}