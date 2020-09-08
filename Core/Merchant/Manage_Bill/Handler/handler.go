package handler

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	elastic1 "github.com/olivere/elastic/v7"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"sync"
	"context"
	"github.com/golang/glog"
	pb "git.zapa.cloud/fresher/kietcdx/Module3/App/Proto/Merchant/ManageBill"
	"github.com/go-redis/redis/v8"  


)

type ElasticDocs struct {
	BillDetailId int `json:"BillDetail_Id"`
	BillId int `json:"Bill_Id"`
	BillStatus int `json:"Bill_Status"`
	Amount int `json:"Amount"`
	ItemId int `json:"Item_Id"`
	CustomerId int `json:"Customer_Id"`
	ItemName string `json:"Item_Name"`
	CustomerName string `json:"Customer_Name"`
	CustomerPhone string `json:"Customer_Phone"`
	BillDesc string `json:"Bill_Desc"`
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

func GetESClient() (*elastic1.Client,error){
	client, err :=  elastic1.NewClient()

	return client,err
}

func JsonStruct(doc ElasticDocs) string {

    docStruct := &ElasticDocs{
        BillId: doc.BillId,
		BillStatus: doc.BillStatus,
		Amount: doc.Amount,
		ItemId: doc.ItemId,
		CustomerId: doc.CustomerId,
		ItemName: doc.ItemName,
		CustomerName: doc.CustomerName,
		CustomerPhone: doc.CustomerPhone,
		BillDesc: doc.BillDesc}

    b, err := json.Marshal(docStruct)
    if err != nil {
        fmt.Println("json.Marshal ERROR:", err)
        return string(err.Error())
    }
    return string(b)
}

func InsertEs(esclient *elastic1.Client,request *pb.CreateBillRequest,BillDetail_Id int) (error){
	DB := AccessDB()

	Query := `SELECT bd.BillDetail_Id,bd.Bill_Id,b.Bill_Status,bd.Amount,bd.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
	FROM BillDetail as bd
	LEFT JOIN Bill AS b ON bd.Bill_Id=b.Bill_Id
	LEFT JOIN Item AS i ON bd.Item_Id=i.Item_Id
	LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id
	WHERE bd.BillDetail_Id=` + strconv.Itoa(BillDetail_Id) + `;`
	bill,err := DB.Query(Query)
	if err != nil {
		log.Fatal(err)
	}
	var doc1 ElasticDocs

	for bill.Next(){
		err = bill.Scan(&doc1.BillDetailId,&doc1.BillId,&doc1.BillStatus,&doc1.Amount,&doc1.ItemId,&doc1.CustomerId,&doc1.ItemName,&doc1.CustomerName,&doc1.CustomerPhone,&doc1.BillDesc)

		fmt.Println(doc1)
		if err != nil {
			log.Fatal(err)
		}

		res,err := esclient.Index().Index(`module3`).Type("_doc").BodyJson(doc1).Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
				
		glog.Info(`Index document`,res.Id)
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


	res, err := esclient2.Search().Index("module3").Size(100).Query(bq).Do(context.Background())

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
		fmt.Println("Query in DB")

		Query:=`SELECT bd.BillDetail_Id,bd.Bill_Id,b.Bill_Status,bd.Amount,bd.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
		FROM BillDetail as bd
		LEFT JOIN Bill AS b ON bd.Bill_Id=b.Bill_Id
		LEFT JOIN Item AS i ON bd.Item_Id=i.Item_Id
		LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id `
		isFirst := true

		if request.GetBillId() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`
			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` bd.Bill_Id="` +strconv.Itoa(int(request.GetBillId()))+`"`
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
			Query = Query + ` bd.Item_Id="` +strconv.Itoa(int(request.GetItemId()))+`"`
		}

		if request.GetAmount() != 0{
			if isFirst{
				isFirst = false
				Query = Query + `WHERE`


			} else {
				Query = Query + ` AND`
			}
			Query = Query + ` bd.Amount="` +strconv.Itoa(int(request.GetAmount()))+`"`
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
			err := billDb.Scan(&billQueryDB.BillDetailId,&billQueryDB.BillId,&billQueryDB.BillStatus,&billQueryDB.Amount,&billQueryDB.ItemId,&billQueryDB.CustomerId,&billQueryDB.ItemName,&billQueryDB.CustomerName,&billQueryDB.CustomerPhone,&billQueryDB.BillDesc)

			LogError(err)

			Wg.Add(1)
			go func(bill ElasticDocs){
				defer Wg.Done()

				billRespone := pb.SearchBillRespone{BillId: int64(bill.BillId),BillStatus: int64(bill.BillStatus),ItemId: int64(bill.ItemId),Amount: int64(bill.Amount), CustomerId: int64(bill.CustomerId), CustomerPhone: bill.CustomerPhone, BillDesc: bill.BillDesc, ItemName: bill.ItemName, CustomerName: bill.CustomerName}

				mux.Lock()
				if err = stream.Send(&billRespone);err!=nil{
					log.Fatal(err)
				}
				mux.Unlock()
				temp := []*pb.CreateBillItem{}
				billRequest := pb.CreateBillRequest{Item:append(temp,&pb.CreateBillItem{ItemId:int64(bill.ItemId),Amount: int64(bill.Amount)}), CustomerId: int64(bill.CustomerId),BillDesc: bill.BillDesc}

				InsertEs(esclient2,&billRequest,bill.BillDetailId)
				fmt.Println(billRespone)
			}(billQueryDB)
			mux.Lock()
			//billsQueryDB = append(billsQueryDB,billQueryDB)
			mux.Unlock()
			Wg.Wait()
		}
		//
	} else {
		fmt.Println("Query in ES")

		for _,billEs := range billsEs {
			Wg.Add(1)
			go func(billEs ElasticDocs){
				defer Wg.Done()

				Query := `SELECT bd.BillDetail_Id,bd.Bill_Id,b.Bill_Status,bd.Amount,bd.Item_Id,b.Customer_Id,i.Item_Name,c.Customer_Name,c.Customer_Phone,b.Bill_Desc
				FROM BillDetail as bd
				LEFT JOIN Bill AS b ON bd.Bill_Id=b.Bill_Id
				LEFT JOIN Item AS i ON bd.Item_Id=i.Item_Id
				LEFT JOIN Customer AS c ON b.Customer_Id=c.Customer_Id
				WHERE bd.Bill_Id="` + strconv.Itoa(billEs.BillId) + `" AND bd.Item_Id="` + strconv.Itoa(billEs.ItemId) + `";`

				mux.Lock()
				billDb,err := DB.Query(Query)
				mux.Unlock()

				LogError(err)


				var bill ElasticDocs

				
				for billDb.Next(){
					err := billDb.Scan(&bill.BillDetailId,&bill.BillId,&bill.BillStatus,&bill.Amount,&bill.ItemId,&bill.CustomerId,&bill.ItemName,&bill.CustomerName,&bill.CustomerPhone,&bill.BillDesc)

					LogError(err)
				}

				billRespone := pb.SearchBillRespone{BillId: int64(bill.BillId),BillStatus: int64(bill.BillStatus),ItemId: int64(bill.ItemId),Amount: int64(bill.Amount), CustomerId: int64(bill.CustomerId), CustomerPhone: bill.CustomerPhone, BillDesc: bill.BillDesc, ItemName: bill.ItemName, CustomerName: bill.CustomerName}

				mux.Lock()
				glog.Info(`Bill Respone `,billRespone)
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

	sqlStatement := `INSERT INTO Bill(Bill_Status,Customer_Id,Bill_Desc,Item_Id) VALUES('`+strconv.Itoa(1)+`','`+strconv.Itoa(int(request.GetCustomerId()))+`','`+request.GetBillDesc()+`','`+strconv.Itoa(0)+`')`
	
	result1,err:= DB.Exec(sqlStatement)

	if err!=nil {
		log.Fatal(err)
	}
	insertId,_ := result1.LastInsertId()
	for _,value := range request.GetItem(){
		if value.GetAmount()==0 {
			continue
		}
		fmt.Print(value.GetItemId())

		sqlStatement := `INSERT INTO BillDetail(Bill_Id,Amount,Item_Id) VALUES('`+strconv.Itoa(int(insertId))+`','`+strconv.Itoa(int(value.GetAmount()))+`','`+strconv.Itoa(int(value.GetItemId()))+`')`
		result,err:= DB.Exec(sqlStatement)

		if err!=nil {
			log.Fatal(err)
		}

		esclient,err := GetESClient()

		LastInsertId,_ := result.LastInsertId()

		err = InsertEs(esclient,request,int(LastInsertId))
	}

	return &pb.CreateBillRespone{IsSaved:1},nil
}