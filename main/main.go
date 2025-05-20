package main

import (
	"goBackend/internal/handler"
	"log"
	"net"
	"net/http"

	_ "goBackend/docs" // подключаем Swagger-доки

	"google.golang.org/grpc"

	"goBackend/api/proto"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Параллельно запускаем gRPC и HTTP серверы
	go startGRPCServer()
	startHTTPServer()
}

// HTTP-сервер (Swagger + /calculate)
func startHTTPServer() {
	h := handler.NewHandler()

	http.HandleFunc("/calculate", h.Calculate)

	// Swagger UI
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("HTTP сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// gRPC-сервер (порт :50051)
func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterCalculatorServer(grpcServer, handler.NewGRPCServer())

	log.Println("gRPC сервер запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
