package handler

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"context"
	pb "google.golang.org/grpc/examples/App/Proto/ManageSupport"
	"github.com/go-redis/redis/v8" 
	"strconv" 
)

type ManageSupportServer struct{
	pb.UnimplementedManageSupportServer
}

type Item struct {
	ItemName string 
	ItemId int
}

type Customer struct {
	CustomerName string
	CustomerId int
	CustomerPhone string
}

type Authorized struct {
	Authorized int
}

type Bill struct {
	ItemName string
	Amount int
	Price int
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

func (sv ManageSupportServer) CheckUserToken(ctx context.Context, request *pb.UserTokenRequest) (*pb.UserTokenRespone,error){
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	DB := AccessDB()

	value,err := redisClient.Get(ctx,request.GetUserToken()).Result()

	if err == redis.Nil {
		return &pb.UserTokenRespone{IsExisted:int64(-1)},nil
	} else {
		if err != nil {
			log.Fatal(err)
		}

		author,err := DB.Query(`SELECT Authorized FROM User WHERE User_Id="` + value + `";`)

		if err != nil {
			log.Fatal(err)
		}

		var Author Authorized

		for author.Next(){

			err := author.Scan(&Author.Authorized)
			if err != nil {
				log.Fatal(err)
			}
		}
		return &pb.UserTokenRespone{IsExisted:int64(1),Authorized:int64(Author.Authorized)},nil
	}
}

func (sv *ManageSupportServer) BillDetail(ctx context.Context,request *pb.BillDetailRequest) (*pb.BillDetailRespone,error){
	DB := AccessDB()

	Query := `SELECT i.Item_Name,b.Amount,i.Price
	FROM Bill as b
	LEFT JOIN Item as i ON b.Item_Id=i.Item_Id
	WHERE b.Bill_Id="
	` + strconv.Itoa(int(request.GetBillId())) + `";`

	billDB,err := DB.Query(Query)

	if err != nil {
		log.Fatal(err)
	}

	var bill Bill
	for billDB.Next(){
		err = billDB.Scan(&bill.ItemName,&bill.Amount,&bill.Price)
		
		if err != nil{
			log.Fatal(err)
		}
	}

	return &pb.BillDetailRespone{ItemName:bill.ItemName,Amount:int64(bill.Amount),Price:int64(bill.Price),Total:int64(bill.Price*bill.Amount)},nil
}