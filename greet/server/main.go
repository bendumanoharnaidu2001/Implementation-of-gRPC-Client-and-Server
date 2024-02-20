package main

//setting up a gRPC server that connects to a PostgreSQL database using the GORM library.
import (
	pb "GoAuth/greet/proto"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
)

func init() {
	DatabaseConnection()
}

type Server struct {
	DB *gorm.DB
	pb.GreetServer
}

var err error

// User Defines a database model
type User struct {
	gorm.Model
	Id         int64
	First_name string
	Last_name  string
	Age        int64
	Token      string
}

func DatabaseConnection() *gorm.DB {
	host := "localhost"
	port := "5433"
	dbName := "testdb"
	dbUser := "postgres"
	password := "1234"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbUser, dbName, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")

	// Auto-migrate models
	db.AutoMigrate(&User{})

	return db
}

var addr string = "0.0.0.0:50051"

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to Listen :%v\n", err)
	}

	log.Printf("listening %s\n", addr)

	s := grpc.NewServer()
	db := DatabaseConnection()
	pb.RegisterGreetServer(s, &Server{DB: db})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v\n", err)
	}

}

func (s *Server) CreatUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	usr := req.User
	token := uuid.New().String()

	users := User{
		Id:         usr.Id,
		First_name: usr.FirstName,
		Last_name:  usr.LastName,
		Age:        usr.Age,
	}

	users.Token = token

	res := s.DB.Create(&users)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie creation unsuccessful")
	}

	response := &pb.CreateUserResponse{
		Token:   users.Token,
		Message: "User successfully created",
	}

	return response, nil

}
