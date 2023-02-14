package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"HappyKod/ServiceShortLinks/internal/app/container"
	"HappyKod/ServiceShortLinks/internal/appproto/midleware"
	pb "HappyKod/ServiceShortLinks/internal/appproto/proto"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer(grpc.UnaryInterceptor(midleware.WorkCooke))

	pb.RegisterLinksServiceServer(server, &LinksService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestLinksService_PutLink(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
		link    string
	}{
		{
			"Получаем ошибку при генерации ссылки",
			codes.InvalidArgument,
			"111",
		},
		{
			"Генерируем валидную ссылку",
			codes.OK,
			"https://www.google.com",
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	md := metadata.New(map[string]string{constans.CookeSessionName: "39636466363139662d326137632d3439463dede60c4a8829edda83bfd46dc019add7966c028655357a1e086706157be8"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.PutLinkRequest{
				Link: tt.link,
			}

			response, errPutLink := client.PutLink(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}

func TestLinksService_PutBatchLink(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
		link    *pb.BatchLink
	}{
		{
			"Получаем ошибку при генерации ссылки",
			codes.InvalidArgument,
			&pb.BatchLink{
				Link: "111",
				Id:   "1",
			},
		},
		{
			"Генерируем валидную ссылку",
			codes.OK,
			&pb.BatchLink{
				Link: "https://www.google.com",
				Id:   "1",
			},
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.PutBatchLinkRequest{
				BatchLink: []*pb.BatchLink{tt.link},
			}
			response, errPutLink := client.PutBatchLink(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}

func TestLinksService_PingDataBase(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
	}{
		{
			"Проверяем пинг",
			codes.OK,
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.PingDataBaseRequest{}
			response, errPutLink := client.PingDataBase(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}

func TestLinksService_GivLink(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
		link    string
	}{
		{
			"не находим ссылку",
			codes.NotFound,
			"1",
		},
		{
			"находим ссылку",
			codes.OK,
			"test",
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	if err = constans.GetLinksStorage().PutShortLink("test", models.Link{FullURL: "https://google.com"}); err != nil {
		t.Fatal(err)
	}
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.GivLinkRequest{Link: tt.link}
			response, errPutLink := client.GivLink(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}

func TestLinksService_GivUserLinks(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
		putLink bool
	}{
		{
			"не находим ссылку",
			codes.NotFound,
			false,
		},
		{
			"находим ссылку",
			codes.OK,
			true,
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	md := metadata.New(map[string]string{constans.CookeSessionName: "39636466363139662d326137632d3439463dede60c4a8829edda83bfd46dc019add7966c028655357a1e086706157be8"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.putLink {
				if err = constans.GetLinksStorage().PutShortLink("test",
					models.Link{FullURL: "https://google.com", UserID: "9cdf619f-2a7c-49"}); err != nil {
					t.Fatal(err)
				}
			}
			request := &pb.GivUserLinksRequest{}
			response, errPutLink := client.GivUserLinks(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}

func TestLinksService_DelUserLinks(t *testing.T) {
	tests := []struct {
		name    string
		errCode codes.Code
		putLink bool
	}{
		{
			"валидное удаление",
			codes.OK,
			true,
		},
	}

	err := container.BuildContainer(models.Config{})
	if err != nil {
		log.Fatal("ошибка инициализации контейнера", err)
	}
	md := metadata.New(map[string]string{constans.CookeSessionName: "39636466363139662d326137632d3439463dede60c4a8829edda83bfd46dc019add7966c028655357a1e086706157be8"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewLinksServiceClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.putLink {
				if err = constans.GetLinksStorage().PutShortLink("test",
					models.Link{FullURL: "https://google.com", UserID: "9cdf619f-2a7c-49"}); err != nil {
					t.Fatal(err)
				}
			}
			request := &pb.DelUserLinksRequest{Links: []string{"test"}}
			response, errPutLink := client.DelUserLinks(ctx, request)
			if response != nil {
				fmt.Println(response.String())
			}
			if errPutLink != nil {
				if er, ok := status.FromError(errPutLink); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
				}
			}
		})
	}
}
