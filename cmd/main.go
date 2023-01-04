package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/config"
	handlers "github.com/evgeniy-dammer/emenu-api/pkg/handlers/items"
	"github.com/evgeniy-dammer/emenu-api/pkg/protos"
)

func main() {
	configuration, err := config.LoadConfiguration()

	if err != nil {
		log.Fatalf("Configuration faild: %s", err)
	}

	/*if err := db.Connect(&configuration); err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}*/

	itemServer := handlers.ItemServiceServer{}
	grpcServer := grpc.NewServer()

	protos.RegisterItemServiceServer(grpcServer, &itemServer)

	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", ":"+configuration.SvPort)

	if err != nil {
		log.Fatalf("Unable to listen on %s: %s", configuration.SvPort, err)
	}

	defer listen.Close()

	log.Printf("Server started at port: %s", configuration.SvPort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Unable to serve: %s", err)
	}
}
