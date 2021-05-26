package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/niakr1s/shtrafovnet/services/companyInfoGetter"
	"google.golang.org/grpc"
)

func main() {
	var inn *string = flag.String("inn", "", "inn of a company to search info for")
	var addr = flag.String("addr", ":9000", "address to listen")
	flag.Parse()

	if *inn == "" {
		flag.Usage()
		os.Exit(1)
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := companyInfoGetter.NewCompanyInfoGetterClient(conn)

	response, err := c.GetCompanyInfo(context.Background(), &companyInfoGetter.GetCompanyInfoRequest{Inn: *inn})
	if err != nil {
		log.Fatalf("Error when calling GetCompanyInfo: %s", err)
	}
	log.Printf("Response from server: %s", response)
}
