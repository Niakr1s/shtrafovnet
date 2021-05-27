package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/niakr1s/shtrafovnet/services/companyInfoGetter"
	"google.golang.org/grpc"
)

var httpAddr = flag.Int("httpAddr", 8081, "address to listen http")
var grpcAddr = flag.Int("grpcAddr", 9000, "address to listen grpc")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	startGrpcServer(ctx)
	startHttpServer(ctx)

	termChan := make(chan os.Signal, 1)

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	cancel()

	time.Sleep(time.Second)
}

func startGrpcServer(ctx context.Context) {
	grpcAddrStr := fmt.Sprintf(":%d", *grpcAddr)
	lis, err := net.Listen("tcp", grpcAddrStr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	companyInfoGetterServer := companyInfoGetter.NewServer()
	companyInfoGetter.RegisterCompanyInfoGetterServer(grpcServer, companyInfoGetterServer)

	go func() {
		log.Printf("grpc server is listening on %s", grpcAddrStr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc server failed to serve: %s", err)
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
}

func startHttpServer(ctx context.Context) {
	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := companyInfoGetter.RegisterCompanyInfoGetterHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", *grpcAddr), opts)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", swaggerHandler()))
	mux.Handle("/", gwmux)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", *httpAddr),
		Handler: mux,
	}

	go func() {
		log.Printf("http server is listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed to serve: %s", err)
		}
	}()

	go func() {
		<-ctx.Done()
		httpServer.Shutdown(context.Background())
	}()
}

//go:embed swagger-ui
var swaggerFS embed.FS

func swaggerHandler() http.Handler {
	swaggerFS, err := fs.Sub(swaggerFS, "swagger-ui")
	if err != nil {
		log.Fatalf("can't get subdirectory of swagger dist: %s", err)
	}

	return http.FileServer(http.FS(swaggerFS))
}
