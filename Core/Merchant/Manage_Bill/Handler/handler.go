package handler

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	elastic "github.com/elastic/go-elasticsearch/v8"
	elastic1 "github.com/olivere/elastic/v7"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"context"
	"github.com/golang/glog"
	pb "google.golang.org/grpc/examples/App/Proto/ManageBill"
	"github.com/go-redis/redis/v8"  


)

type ElasticDocs struct {
	Bill_Id int
	Bill_Status int
	Amount int
	Item_Id int
	Customer_Id int
	Item_Name string
	Customer_Name string
	Customer_Phone string
	Bill_Desc string
}

type Authorized struct {
	Authorized int
}
func AccessDB() (*sql.DB){
	dbDriver := "mysql"
 
    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	return db
}

func GetESClient() (*elastic.Client,error){
	client, err :=  elastic.NewClient(elastic.Config{Addresses: []string{"http://localhost:9200"}})

	return client,err
}

func JsonStruct(doc ElasticDocs) string {

    docStruct := &ElasticDocs{
        Bill_Id: doc.Bill_Id,
		Bill_Status: doc.Bill_Status,
		Amount: doc.Amount,
		Item_Id: doc.Item_Id,
		Customer_Id: doc.Customer_Id,
		Item_Name: doc.Item_Name,
		Customer_Name: doc.Customer_Name,
		Customer_Phone: doc.Customer_Phone,
		Bill_Desc: doc.Bill_Desc}

    b, err := json.Marshal(docStruct)
    if err != nil {
        fmt.Println("json.Marshal ERROR:", err)
        return string(err.Error())
    }
    return string(b)
}

func InsertEs(esclient *elastic.Client,request *pb.CreateBillRequest,Bill_Id int) (error){
	DB := AccessDB()

	Query := `SELECT b.Bill_Id,b.Bill_Status,b.Amount,b.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
	FROM Bill AS b
	LEFT JOIN Item AS i ON b.Item_Id=i.Item_Id
	LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id
	WHERE Bill_Id=` + strconv.Itoa(Bill_Id) + `;`
	bill,err := DB.Query(Query)
	column,err := bill.Columns()
	fmt.Println(column)
	if err != nil {
		log.Fatal(err)
	}
	var doc1 ElasticDocs

	for bill.Next(){
		err = bill.Scan(&doc1.Bill_Id,&doc1.Bill_Status,&doc1.Amount,&doc1.Item_Id,&doc1.Customer_Id,&doc1.Item_Name,&doc1.Customer_Name,&doc1.Customer_Phone,&doc1.Bill_Desc)
	}

	fmt.Println(doc1)
	req := esapi.IndexRequest{
		Index:      "module3",
		Body:       strings.NewReader(JsonStruct(doc1)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), esclient)
	
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
			
	var resMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			log.Printf("Error parsing the response body: %s", err)
	} else {
		log.Printf("\nIndexRequest() RESPONSE:")
		fmt.Println("Status:", res.Status())
		fmt.Println("Result:", resMap["result"])
		fmt.Println("Version:", int(resMap["_version"].(float64)))
	}

	return nil
}

func LogError(err error){
	if err !=nil {
		log.Fatal(err)
	}
}

type ManageBillServer struct {
	pb.UnimplementedManageBillServer
}


func (ms *ManageBillServer) SearchBill(request *pb.SearchBillRequest,stream pb.ManageBill_SearchBillServer) error{
	esclient,err := elastic1.NewClient()
	LogError(err)
	esclient2,err:= GetESClient()
	LogError(err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	userId,_ := redisClient.Get(context.Background(),request.GetUserToken()).Result()

	DB := AccessDB()

	user,err := DB.Query(`SELECT Authorized FROM User WHERE User_Id="` + userId + `";`)

	var Authorized Authorized
	for user.Next(){
		err = user.Scan(&Authorized.Authorized)

		if err !=nil {
			log.Fatal(err)
		}
	}
	bq := elastic1.NewBoolQuery()

	if request.GetBillId() != 0{
		bq.Must(elastic1.NewMatchQuery("Bill_Id",request.GetBillId()))
	}

	if request.GetBillStatus() != 0{
		bq.Must(elastic1.NewMatchQuery("Bill_Status",request.GetBillStatus()))
	}

	if request.GetItemId() != 0{
		bq.Must(elastic1.NewMatchQuery("Item_Id",request.GetItemId()))
	}

	if request.GetAmount() != 0{
		bq.Must(elastic1.NewMatchQuery("Amount",request.GetAmount()))
	}

	if Authorized.Authorized==1 && request.GetCustomerId() != 0{
		bq.Must(elastic1.NewMatchQuery("Customer_Id",request.GetCustomerId()))
	}

	if Authorized.Authorized == 2 {
		bq.Must(elastic1.NewMatchQuery("Customer_Id",userId))
	}

	if request.GetCustomerPhone() != "" && request.GetCustomerPhone() != " " {
		bq.Must(elastic1.NewMatchQuery("Customer_Phone",request.GetCustomerPhone()))
	}

	if request.GetBillDesc() != "" && request.GetBillDesc() != " "{
		bq.Must(elastic1.NewMatchQuery("Bill_Desc",request.GetBillDesc()))
	}

	if request.GetItemName() != ""&&request.GetItemName() != " "{
		bq.Must(elastic1.NewMatchQuery("Item_Name",request.GetItemName()))
	}

	if request.GetCustomerName() != "" &&request.GetCustomerName() != " "{
		bq.Must(elastic1.NewMatchQuery("Customer_Name",request.GetCustomerName()))
	}


	res, err := esclient.Search().Index("module3").Size(100).Query(bq).Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var billsEs []ElasticDocs

	for _,hit:=range res.Hits.Hits{
		var bill ElasticDocs

		jsonString,_ := hit.Source.MarshalJSON()
		if err := json.Unmarshal(jsonString,&bill);err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Println(bill)
		billsEs = append(billsEs,bill)
	}

	Wg := sync.WaitGroup{}
	var mux sync.Mutex
	if len(billsEs) == 0{
		Query:=`SELECT b.Bill_Id,b.Bill_Status,b.Amount,b.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
		FROM Bill AS b
		LEFT JOIN Item AS i ON b.Item_Id=i.Item_Id
		LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id `
		isFirst := true

		if request.GetBillId() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`
			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Bill_Id="` +strconv.Itoa(int(request.GetBillId()))+`"`
		}

		if request.GetBillStatus() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Bill_Status="` +strconv.Itoa(int(request.GetBillStatus()))+`"`
		}

		if request.GetItemId() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Item_Id="` +strconv.Itoa(int(request.GetItemId()))+`"`
		}

		if request.GetAmount() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`


			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Amount="` +strconv.Itoa(int(request.GetAmount()))+`"`
		}

		if Authorized.Authorized==1 && request.GetCustomerId() != 0 {
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Customer_Id="` + strconv.Itoa(int(request.GetCustomerId()))+`"`
		}

		if Authorized.Authorized==2 {
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Customer_Id="` + userId +`"`
		}
		if request.GetCustomerPhone() != "" && request.GetCustomerPhone() != " "{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` c.Customer_Phone="` + request.GetCustomerPhone() + `"`
		}

		if request.GetBillDesc() != "" && request.GetBillDesc() != " "{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` b.Bill_Desc LIKE '%` + request.GetBillDesc()+ `%'`
		}

		if request.GetItemName() != "" && request.GetItemName() != " "{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` I.Item_Name LIKE '%` + request.GetItemName()+ `'`
		}

		if request.GetCustomerName() != "" && request.GetCustomerName() != " "{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`

			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` c.Customer_Name LIKE '%` + request.GetCustomerName()+ `%'`
		}

		Query = Query + `;`
		
		fmt.Println(Query)
		billDb,err := DB.Query(Query)
		//fmt.Println(1)

		LogError(err)
		//var billsQueryDB []ElasticDocs
	
		for billDb.Next(){
			var billQueryDB ElasticDocs
			fmt.Println(billDb)
			err := billDb.Scan(&billQueryDB.Bill_Id,&billQueryDB.Bill_Status,&billQueryDB.Amount,&billQueryDB.Item_Id,&billQueryDB.Customer_Id,&billQueryDB.Item_Name,&billQueryDB.Customer_Name,&billQueryDB.Customer_Phone,&billQueryDB.Bill_Desc)

			LogError(err)

			Wg.Add(1)
			go func(bill ElasticDocs){
				defer Wg.Done()

				billRespone := pb.SearchBillRespone{BillId: int64(bill.Bill_Id),BillStatus: int64(bill.Bill_Status),ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id), CustomerPhone: bill.Customer_Phone, BillDesc: bill.Bill_Desc, ItemName: bill.Item_Name, CustomerName: bill.Customer_Name}

				mux.Lock()
				if err = stream.Send(&billRespone);err!=nil{
					log.Fatal(err)
				}
				mux.Unlock()
				billRequest := pb.CreateBillRequest{ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id),BillDesc: bill.Bill_Desc}

				InsertEs(esclient2,&billRequest,bill.Bill_Id)
				fmt.Println(billRespone)
			}(billQueryDB)
			mux.Lock()
			//billsQueryDB = append(billsQueryDB,billQueryDB)
			mux.Unlock()
			Wg.Wait()
		}
		//
	} else {
		for _,billEs := range billsEs {
			Wg.Add(1)
			go func(billEs ElasticDocs){
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

				LogError(err)


				var bill ElasticDocs

				
				for billDb.Next(){
					err := billDb.Scan(&bill.Bill_Id,&bill.Bill_Status,&bill.Amount,&bill.Item_Id,&bill.Customer_Id,&bill.Item_Name,&bill.Customer_Name,&bill.Customer_Phone,&bill.Bill_Desc)

					LogError(err)
				}

				billRespone := pb.SearchBillRespone{BillId: int64(bill.Bill_Id),BillStatus: int64(bill.Bill_Status),ItemId: int64(bill.Item_Id),Amount: int64(bill.Amount), CustomerId: int64(bill.Customer_Id), CustomerPhone: bill.Customer_Phone, BillDesc: bill.Bill_Desc, ItemName: bill.Item_Name, CustomerName: bill.Customer_Name}

				mux.Lock()
				if err = stream.Send(&billRespone);err!=nil{
					log.Fatal(err)
				}
				mux.Unlock()
			}(billEs)
		}
		Wg.Wait()
	}
	

	return nil
}

func (ms *ManageBillServer) CreateBill(context context.Context, request *pb.CreateBillRequest) (*pb.CreateBillRespone,error){
	DB := AccessDB()
	glog.Info("Create Bill request",request)
	sqlStatement := `INSERT INTO Bill(Bill_Status,Amount,Item_Id,Customer_Id,Bill_Desc) VALUES('`+strconv.Itoa(1)+`','`+strconv.Itoa(int(request.GetAmount()))+`','`+strconv.Itoa(int(request.GetItemId()))+`','`+strconv.Itoa(int(request.GetCustomerId()))+`','`+request.GetBillDesc()+`')`

	result,err:= DB.Exec(sqlStatement)

	if err!=nil {
		log.Fatal(err)
	}
	esclient,err := GetESClient()

	LastInsertId,_ := result.LastInsertId()
	
	err = InsertEs(esclient,request,int(LastInsertId))

	return &pb.CreateBillRespone{IsSaved:1},nil
}