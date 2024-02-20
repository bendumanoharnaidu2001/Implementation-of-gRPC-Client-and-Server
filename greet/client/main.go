package main

import (
	pb "GoAuth/greet/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type UserDetails struct {
	Id         int64  `json: "id"`
	First_name string `json: "first_name"`
	Last_name  string `json: "second_name"`
	Age        int64  `json: "age"`
}

type CreateResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type GetResponse struct {
	Message string `json:"message"`
}

type UpdateResponse struct {
	Message string `json:"message"`
}

var addr string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewGreetClient(conn)

	router := mux.NewRouter()

	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		Create(c, w, r)
	}).Methods("POST")

	//router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	GetUser(c, w, r)
	//}).Methods("GET")
	//
	//router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
	//	UpdateUser(c, w, r)
	//}).Methods("PUT")

	http.ListenAndServe(":8088", router)
}
func Create(client pb.GreetClient, w http.ResponseWriter, r *http.Request) {

	var usr UserDetails

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return
	}

	fmt.Print(usr)

	user := &pb.User{
		Id:        usr.Id,
		FirstName: usr.First_name,
		LastName:  usr.Last_name,
		Age:       usr.Age,
	}
	res, err := client.CreatUser(context.Background(), &pb.CreateUserRequest{
		User: user,
	})

	if err != nil {
		response := CreateResponse{
			Token:   "",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CreateResponse{
		Token:   res.Token,
		Message: res.Message,
	}
	json.NewEncoder(w).Encode(response)

}
