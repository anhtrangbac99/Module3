package main

import (
	pb "google.golang.org/grpc/examples/App/Proto/ManageBill"
	elastic "github.com/olivere/elastic/v7"
	grpc "google.golang.org/grpc"
	_ "github.com/go-sql-driver/mysql"
	"context"
	"strconv"
	"log"
	"fmt"
	"sync"
	"net"
	"encoding/json"
	utils "google.golang.org/grpc/examples/App/Core/Manage_Bill/Handler"
)
type ManageBillServer struct {
	pb.UnimplementedManageBillServer
}


func (ms *ManageBillServer) SearchBill(request *pb.SearchBillRequest,stream pb.ManageBill_SearchBillServer) error{
	esclient,err := elastic.NewClient()
	utils.LogError(err)
	esclient2,err:= utils.GetESClient()
	utils.LogError(err)
	
	bq := elastic.NewBoolQuery()

	if request.GetBillId() != 0{
		bq.Must(elastic.NewMatchQuery("Bill_Id",request.GetBillId()))
	}

	if request.GetBillStatus() != 0{
		bq.Must(elastic.NewMatchQuery("Bill_Status",request.GetBillStatus()))
	}

	if request.GetItemId() != 0{
		bq.Must(elastic.NewMatchQuery("Item_Id",request.GetItemId()))
	}

	if request.GetAmount() != 0{
		bq.Must(elastic.NewMatchQuery("Amount",request.GetAmount()))
	}

	if request.GetCustomerId() != 0{
		bq.Must(elastic.NewMatchQuery("Customer_Id",request.GetCustomerId()))
	}

	if request.GetCustomerPhone() != ""{
		bq.Must(elastic.NewMatchQuery("Customer_Phone",request.GetCustomerPhone()))
	}

	if request.GetBillDesc() != ""{
		bq.Must(elastic.NewMatchQuery("Bill_Desc",request.GetBillDesc()))
	}

	if request.GetItemName() != ""{
		bq.Must(elastic.NewMatchQuery("Item_Name",request.GetItemName()))
	}

	if request.GetCustomerName() != ""{
		bq.Must(elastic.NewMatchQuery("Customer_Name",request.GetCustomerName()))
	}

	
	res, err := esclient.Search().Index("module3").Query(bq).Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var billsEs []utils.ElasticDocs

	for _,hit:=range res.Hits.Hits{
		var bill utils.ElasticDocs

		jsonString,_ := hit.Source.MarshalJSON()
		if err := json.Unmarshal(jsonString,&bill);err != nil {
			log.Fatal(err)
			continue
		}
		
		billsEs = append(billsEs,bill)
	}



	DB := utils.AccessDB()
	Wg := sync.WaitGroup{}
	var mux sync.Mutex
	if len(billsEs) == 0{
		Query:=`SELECT b.Bill_Id,b.Bill_Status,b.Amount,b.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
		FROM Bill AS b
		LEFT JOIN Item AS i ON b.Item_Id=i.Item_Id
		LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id
		WHERE b.Bill_Id="` + strconv.Itoa(int(request.GetBillId())) + `" OR b.Bill_Status="` + strconv.Itoa(int(request.GetBillStatus())) + `" OR b.Amount="` + strconv.Itoa(int(request.GetAmount())) + `" OR b.Item_Id="` + strconv.Itoa(int(request.GetItemId())) +`" OR b.Customer_Id="` + strconv.Itoa(int(request.GetCustomerId())) +`" OR i.Item_Name="` + request.GetItemName()+`" OR c.Customer_Name="` + request.GetCustomerName()+`" OR c.Customer_Phone="` + request.GetCustomerPhone()+`" OR b.Bill_Desc="` + request.GetBillDesc()+`";`
		billDb,err := DB.Query(Query)
		fmt.Println(1)

		utils.LogError(err)
		var billsQueryDB []utils.ElasticDocs
	
		for billDb.Next(){
			var billQueryDB utils.ElasticDocs
			err := billDb.Scan(&billQueryDB.Bill_Id,&billQueryDB.Bill_Status,&billQueryDB.Amount,&billQueryDB.Item_Id,&billQueryDB.Customer_Id,&billQueryDB.Item_Name,&billQueryDB.Customer_Name,&billQueryDB.Customer_Phone,&billQueryDB.Bill_Desc)
			
			utils.LogError(err)
			Wg.Add(1)
			go func(bill utils.ElasticDocs){
				defer Wg.Done()
				billRespone := pb.SearchBillRespone{BillId: int64(bill.Bill_Id),BillStatus: int64(bill.Bill_Status),ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id), CustomerPhone: bill.Customer_Phone, BillDesc: bill.Bill_Desc, ItemName: bill.Item_Name, CustomerName: bill.Customer_Name}

				mux.Lock()
				if err = stream.Send(&billRespone);err!=nil{
					log.Fatal(err)
				}
				mux.Unlock()

				billRequest := pb.CreateBillRequest{ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id),BillDesc: bill.Bill_Desc}

			
				utils.InsertEs(esclient2,&billRequest,bill.Bill_Id)

				fmt.Println(billRespone)
			}(billQueryDB)
			billsQueryDB = append(billsQueryDB,billQueryDB)
		}
		Wg.Wait()
	} else {
		for _,billEs := range billsEs {
			go func(billEs utils.ElasticDocs){
				defer Wg.Done()

				billId := billEs.Bill_Id
				Query := `SELECT b.Bill_Id,b.Bill_Status,b.Amount,b.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
				FROM Bill AS b
				LEFT JOIN Item AS i ON b.Item_Id=i.Item_Id
				LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id
				WHERE Bill_Id="` + strconv.Itoa(billId) + `";`

				mux.Lock()
				billDb,err := DB.Query(Query)
				mux.Unlock()

				utils.LogError(err)


				var bill utils.ElasticDocs

				
				for billDb.Next(){
					err := billDb.Scan(&bill.Bill_Id,&bill.Bill_Status,&bill.Amount,&bill.Item_Id,&bill.Customer_Id,&bill.Item_Name,&bill.Customer_Name,&bill.Customer_Phone,&bill.Bill_Desc)

					utils.LogError(err)
				}

				billRespone := pb.SearchBillRespone{BillId: int64(bill.Bill_Id),BillStatus: int64(bill.Bill_Status),ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id), CustomerPhone: bill.Customer_Phone, BillDesc: bill.Bill_Desc, ItemName: bill.Item_Name, CustomerName: bill.Customer_Name}

				mux.Lock()
				if err = stream.Send(&billRespone);err!=nil{
					log.Fatal(err)
				}
				mux.Unlock()
			}(billEs)
		}
	}
	Wg.Wait()

	return nil
}

func (ms *ManageBillServer) CreateBill(context context.Context, request *pb.CreateBillRequest) (*pb.CreateBillRespone,error){
	DB := utils.AccessDB()
	
	sqlStatement := `INSERT INTO Bill(Bill_Status,Amount,Item_Id,Customer_Id,Bill_Desc) VALUES('`+strconv.Itoa(1)+`','`+strconv.Itoa(int(request.GetAmount()))+`','`+strconv.Itoa(int(request.GetItemId()))+`','`+strconv.Itoa(int(request.GetCustomerId()))+`','`+request.GetBillDesc()+`')`

	result,err:= DB.Exec(sqlStatement)

	if err!=nil {
		log.Fatal(err)
	}

	esclient,err := utils.GetESClient()

	LastInsertId,_ := result.LastInsertId()
	
	err = utils.InsertEs(esclient,request,int(LastInsertId))

	return &pb.CreateBillRespone{IsSaved:1},nil
}

func main(){

	grpcServer := grpc.NewServer()
	pb.RegisterManageBillServer(grpcServer,new(ManageBillServer))

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