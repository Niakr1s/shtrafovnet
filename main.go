package main

import (
	"flag"
	"log"
	"net"

	"github.com/niakr1s/shtrafovnet/services/companyInfoGetter"
	"google.golang.org/grpc"
)

func main() {
	var addr = flag.String("addr", ":9000", "address to listen")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	companyInfoGetterServer := companyInfoGetter.NewServer()
	companyInfoGetter.RegisterCompanyInfoGetterServer(grpcServer, companyInfoGetterServer)

	log.Printf("Server is listening on %s", *addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
