// Package server запуск сервера.
package server

import (
	"fmt"
	"log"
	"net"

	"HappyKod/ServiceShortLinks/internal/appproto/handlers"
	"HappyKod/ServiceShortLinks/internal/appproto/midleware"
	pb "HappyKod/ServiceShortLinks/internal/appproto/proto"
	"HappyKod/ServiceShortLinks/internal/models"

	"google.golang.org/grpc"
)

// NewServer создания сервера с настройками.
func NewServer(cfg models.Config) {
	listen, err := net.Listen("tcp", cfg.AddressProto)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(midleware.WorkCooke))
	// регистрируем сервис
	pb.RegisterLinksServiceServer(s, &handlers.LinksService{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
