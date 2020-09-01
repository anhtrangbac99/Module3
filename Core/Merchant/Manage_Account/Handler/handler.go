package handler

import (
	pb "google.golang.org/grpc/examples/App/Proto/SignIn"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"  
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
	"time"
	"strconv"
)

type selectedUser struct {
	User_Id int
	Username string 
	Password string
	Authorized int
}
type SignInServer struct{
	pb.UnimplementedSignInServer
}

func AccessDB() (*sql.DB,*redis.Client){
	dbDriver := "mysql"
	
    db, err := sql.Open(dbDriver, "root:@(127.0.0.1:4000)/module3")
	
	if err != nil {
        log.Fatal(err)
	}
	
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})
	return db,redisClient
}

func GenerateToken() string {
	timestamp := time.Now()
    hash, err := bcrypt.GenerateFromPassword([]byte(timestamp.Format("20060102150405")), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
	}
	
    hasher := md5.New()
    hasher.Write(hash)
    return hex.EncodeToString(hasher.Sum(nil))
}

func (signInServer *SignInServer) UserAuthor(ctx context.Context,request *pb.AuthorRequest) (*pb.AuthorRespone,error){
	DB,RDB:= AccessDB()
	

	fmt.Println(request.GetUsername())

	Query := `SELECT User_Id,Username,Password,Authorized FROM User WHERE`

	Query = Query + ` Username = "` + request.GetUsername() + `";`

	selectedUsers, err := DB.Query(Query)
	
	isExisted := -1
	user_Id := "-1"
	authorized := -1

	if err != nil {
		log.Fatal(err)
	} else {
		for selectedUsers.Next(){
			var selectedUser selectedUser
			err = selectedUsers.Scan(&selectedUser.User_Id,&selectedUser.Username,&selectedUser.Password,&selectedUser.Authorized)
			//fmt.Println(selectedUser)
			if err != nil {
				log.Fatal(err)
			}

			if selectedUser.Password == request.GetPassword() {
				isExisted = 1
				user_Id = GenerateToken()
				authorized = selectedUser.Authorized
				fmt.Println(user_Id)
				
				err = RDB.Set(context.Background(),user_Id,strconv.Itoa(selectedUser.User_Id),15*time.Minute).Err()

				if err !=nil{
					log.Fatal(err)
				}
				break;
			} else {
				isExisted = 1
			}
		}
	}

	return &pb.AuthorRespone{IsExisted:int64(isExisted),User_Id:user_Id,Authorized:int64(authorized)},nil
}
