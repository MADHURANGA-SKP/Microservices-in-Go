package main

import (
	common "common"
	pb "common/api"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR",":8080")
	orderServiceAddr = "localhost:2000"
)

func main(){

	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing orders service at", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)


	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start the server")
	}
}