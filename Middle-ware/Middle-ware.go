package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	//"io/ioutil"
	//"reflect"
	"context"
	//"strconv"
	_ "github.com/go-sql-driver/mysql"
	pbSignIn "google.golang.org/grpc/examples/App/Proto/SignIn"
	grpc "google.golang.org/grpc"

)

type User struct{
	Username string 
	Password string
}

type selectedUser struct {
	User_Id int
	Username string 
	Password string
}

type Respone struct {
	respone bool
}

func AccessDB() (*sql.DB){
	dbDriver := "mysql"

    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	
	return db
}

func main()  {
	router := mux.NewRouter()

	router.HandleFunc("/SignIn",SignIn).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":1234", handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(router)))

}

func SignIn(writer http.ResponseWriter,request *http.Request){
	var user User
	//var respone http.Response
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	connect,err := grpc.Dial("localhost:5000",grpc.WithInsecure(),grpc.WithBlock())
	defer connect.Close()

	if err != nil {
		log.Fatal(err)
	}

	client := pbSignIn.NewSignInClient(connect)

	result,err := client.UserAuthor(context.Background(),&pbSignIn.AuthorRequest{Username: user.Username , Password: user.Password})
	fmt.Println(result.GetUser_Id())
	fmt.Println(result)
	respone,err := json.Marshal(&result)

	writer.Write(respone)

	

	return
	
}

