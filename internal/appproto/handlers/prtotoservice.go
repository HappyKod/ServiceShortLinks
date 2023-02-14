// Package handlers Описана работа с proto.
package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	pb "HappyKod/ServiceShortLinks/internal/appproto/proto"
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage/pglinkssotorage"
	"HappyKod/ServiceShortLinks/utils"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LinksService поддерживает все необходимые методы сервера.
type LinksService struct {
	pb.UnimplementedLinksServiceServer
}

const errorReadUserID = "ошибка чтения userId"

// PutLink реализует интерфейс добавление ссылки и создания короткой.
func (s *LinksService) PutLink(ctx context.Context, in *pb.PutLinkRequest) (*pb.PutLinkResponse, error) {
	var userID string
	var cooke string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userID = md.Get(constans.CookeUserIDName)[0]
		cooke = md.Get(constans.CookeSessionName)[0]
	} else {
		log.Println(errorReadUserID)
		return nil, status.Error(codes.Unauthenticated, errorReadUserID)
	}
	linksStorage := constans.GetLinksStorage()
	var response pb.PutLinkResponse
	response.Cooke = cooke
	if !utils.ValidatorURL(in.Link) {
		return nil, status.Errorf(codes.InvalidArgument, constans.ErrorInvalidURL)
	}
	link := models.Link{
		ShortKey: utils.GeneratorStringUUID(),
		FullURL:  in.Link,
		UserID:   userID,
		Created:  time.Now(),
	}
	var key string
	if err := linksStorage.PutShortLink(link.ShortKey, link); err != nil {
		if errPG, ok := err.(*pq.Error); ok {
			if errPG.Code != pgerrcode.UniqueViolation {
				log.Println(constans.ErrorWriteStorage, in.Link, err)
				return nil, status.Errorf(codes.Aborted, constans.ErrorWriteStorage)
			}
		} else {
			if !errors.Is(constans.ErrorNoUNIQUEFullURL, err) {
				log.Println(constans.ErrorWriteStorage, in.Link, err)
				return nil, status.Errorf(codes.Aborted, constans.ErrorWriteStorage)
			}
		}
		key, err = linksStorage.GetKey(in.Link)
		if err != nil {
			log.Println(constans.ErrorGetKeyStorage, in.Link, err)
			return nil, status.Errorf(codes.Aborted, constans.ErrorGetKeyStorage)
		}
	} else {
		key = link.ShortKey
	}
	uri, err := utils.GenerateURL(key)
	if err != nil {
		log.Println(constans.ErrorGenerateURL, key, err)
		return nil, status.Errorf(codes.Aborted, constans.ErrorGenerateURL)
	}
	response.ShortLink = uri
	return &response, nil
}

// PutBatchLink реализует интерфейс добавление множества ссылок и создания коротких.
func (s *LinksService) PutBatchLink(ctx context.Context, in *pb.PutBatchLinkRequest) (*pb.PutBatchLinkResponse, error) {
	var userID string
	var cooke string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userID = md.Get(constans.CookeUserIDName)[0]
		cooke = md.Get(constans.CookeSessionName)[0]
	} else {
		log.Println(errorReadUserID)
		return nil, status.Error(codes.Unauthenticated, errorReadUserID)
	}
	var response pb.PutBatchLinkResponse
	response.Cooke = cooke
	var links []models.Link
	linksStorage := constans.GetLinksStorage()
	for _, v := range in.BatchLink {
		if !utils.ValidatorURL(v.Link) {
			return nil, status.Errorf(codes.InvalidArgument, constans.ErrorInvalidURL)
		}
		links = append(links, models.Link{
			FullURL:  v.Link,
			ShortKey: utils.GeneratorStringUUID(),
			UserID:   userID,
			Created:  time.Now(),
		})
	}
	if err := linksStorage.ManyPutShortLink(links); err != nil {
		log.Println(constans.ErrorWriteStorage, err.Error())
		return nil, status.Error(codes.Aborted, constans.ErrorWriteStorage)
	}
	for _, link := range links {
		for _, v := range in.BatchLink {
			if v.Link == link.FullURL {
				shortURL, err := utils.GenerateURL(link.ShortKey)
				if err != nil {
					log.Println(constans.ErrorGenerateURL, link.ShortKey, err)
					return nil, status.Error(codes.Aborted, constans.ErrorGenerateURL)
				}
				response.BatchLink = append(response.BatchLink, &pb.BatchLink{
					Id:   v.Id,
					Link: shortURL,
				})
				break
			}
		}
	}
	return &response, nil
}

// PingDataBase реализует интерфейс проверки соединения с базой.
func (s *LinksService) PingDataBase(_ context.Context, _ *pb.PingDataBaseRequest) (*pb.PingDataBaseResponse, error) {
	var response pb.PingDataBaseResponse
	cfg := constans.GlobalContainer.Get("server-config").(models.Config)
	if cfg.DataBaseURL == "" {
		log.Println("NoDataBaseURL")
		return &response, nil
	}
	linkStorage, err := pglinkssotorage.New(cfg.DataBaseURL)
	if err != nil {
		log.Println(constans.ErrorConnectStorage, err)
		return nil, status.Error(codes.Aborted, constans.ErrorConnectStorage)
	}
	if err = linkStorage.Ping(); err != nil {
		log.Println(constans.ErrorConnectStorage, err)
		return nil, status.Error(codes.Aborted, constans.ErrorConnectStorage)
	}
	return &response, nil
}

// GivLink реализует интерфейс получение полной ссылки по короткой.
func (s *LinksService) GivLink(_ context.Context, in *pb.GivLinkRequest) (*pb.GivLinkResponse, error) {
	var response pb.GivLinkResponse
	link, err := constans.GetLinksStorage().GetShortLink(in.Link)
	if err != nil {
		log.Println(constans.ErrorReadStorage, in.Link, err.Error())
		return nil, status.Error(codes.Aborted, constans.ErrorReadStorage)
	}
	if link.FullURL == "" {
		return nil, status.Error(codes.NotFound, "Ошибка по ключу ничего не нашлось")
	}
	if link.Del {
		return nil, status.Error(codes.NotFound, "Ошибка данная ссылка больше не доступна")
	}
	response.Link = link.FullURL
	return &response, nil
}

// GivUserLinks реализует интерфейс получение всех ссылок пользователя.
func (s *LinksService) GivUserLinks(ctx context.Context, _ *pb.GivUserLinksRequest) (*pb.GivUserLinksResponse, error) {
	var userID string
	var cooke string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userID = md.Get(constans.CookeUserIDName)[0]
		cooke = md.Get(constans.CookeSessionName)[0]
	} else {
		log.Println(errorReadUserID)
		return nil, status.Error(codes.Unauthenticated, errorReadUserID)
	}
	var response pb.GivUserLinksResponse
	response.Cooke = cooke
	linksStorage := constans.GetLinksStorage()
	links, err := linksStorage.GetShortLinkUser(userID)
	if err != nil {
		return nil, status.Error(codes.Aborted, constans.ErrorReadStorage)
	}
	for _, link := range links {
		shortLink, errGenerateURL := utils.GenerateURL(link.ShortKey)
		if errGenerateURL != nil {
			log.Println(constans.ErrorGenerateURL, link.ShortKey, errGenerateURL)
			return nil, status.Error(codes.Aborted, constans.ErrorGenerateURL)
		}
		response.Links = append(response.Links, &pb.Links{
			ShortUri: shortLink,
			FullUri:  link.FullURL,
		})
	}
	if len(response.Links) == 0 {
		return nil, status.Error(codes.NotFound, "Not found links")
	}
	return &response, nil
}

// DelUserLinks реализует интерфейс удаления ссылок пользователя.
func (s *LinksService) DelUserLinks(ctx context.Context, in *pb.DelUserLinksRequest) (*pb.DelUserLinksResponse, error) {
	var userID string
	var cooke string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userID = md.Get(constans.CookeUserIDName)[0]
		userID = md.Get(constans.CookeSessionName)[0]
	} else {
		log.Println(errorReadUserID)
		return nil, status.Error(codes.Unauthenticated, errorReadUserID)
	}
	var response pb.DelUserLinksResponse
	response.Cooke = cooke
	linksStorage := constans.GetLinksStorage()
	go func() {
		err := linksStorage.DeleteShortLinkUser(userID, in.Links)
		if err != nil {
			log.Println(constans.ErrorUpdateStorage, in.Links, err)
		}
	}()
	return &response, nil
}
