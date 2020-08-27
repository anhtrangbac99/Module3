package handler

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"context"
	pb "google.golang.org/grpc/examples/App/Proto/ManageSupport"
	
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



func (sv *ManageSupportServer) GetListItem(request *pb.ListItemRequest,stream pb.ManageSupport_GetListItemServer) error {
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
