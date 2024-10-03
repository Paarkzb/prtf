package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	ssov1 "prtf-gateway/protos/gen/go/sso"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cc, err := grpc.NewClient("localhost:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	defer cc.Close()

	authClient := ssov1.NewAuthClient(cc)

	http.HandleFunc("/signUp", func(w http.ResponseWriter, r *http.Request) {
		respReg, err := authClient.SignUp(context.Background(), &ssov1.SignUpRequest{
			Username: "testserv",
			Email:    "testserv@test.ru",
			Password: "testserv",
		})
		if err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(respReg)
	})
	http.HandleFunc("/signIn", func(w http.ResponseWriter, r *http.Request) {
		respReg, err := authClient.SignIn(context.Background(), &ssov1.SignInRequest{
			Username: "testserv",
			Password: "testserv",
		})
		if err != nil {
			log.Println(err)
		}

		json.NewEncoder(w).Encode(respReg)
	})

	http.ListenAndServe(":8071", nil)
}
