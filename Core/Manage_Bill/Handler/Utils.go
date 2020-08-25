package Utils

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"encoding/json"
	"log"
	"fmt"
	"strconv"
	"strings"
	"context"
	pb "google.golang.org/grpc/examples/App/Proto/ManageBill"

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