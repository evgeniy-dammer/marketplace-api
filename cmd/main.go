package main

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/pkg/handlers/items"
	itemservice "github.com/evgeniy-dammer/emenu-api/pkg/proto"
	"log"

	"google.golang.org/grpc"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":1111")

	if err != nil {
		fmt.Println(err)
	}

	defer listen.Close()

	itemServ := handlers.ItemServiceServer{}
	grpcServer := grpc.NewServer()

	itemservice.RegisterItemServiceServer(grpcServer, &itemServ)

	log.Println("Starting server")
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
