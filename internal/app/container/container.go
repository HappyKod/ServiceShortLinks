package container

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/links-storage"
	"HappyKod/ServiceShortLinks/internal/storage/links-storage/files-links-torage"
	"HappyKod/ServiceShortLinks/internal/storage/links-storage/mem-links-storage"
	"HappyKod/ServiceShortLinks/internal/storage/users-storage"
	"HappyKod/ServiceShortLinks/internal/storage/users-storage/mem-users-storage"
	"github.com/sarulabs/di"
	"log"
)

func BuildContainer(cfg models.Config) error {
	var linksStorage links_storage.LinksStorages
	if cfg.FileStoragePATH != "" {
		store, err := files_links_torage.New(cfg.FileStoragePATH)
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован file-links-storage")
	} else {
		store, err := mem_links_storage.New()
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован mem-links-storage")
	}
	var usersStorage users_storage.UsersStorage
	usersStorage, err := mem_users_storage.New()
	if err != nil {
		return err
	}
	log.Println("Задействован mem_users_storage")
	builder, _ := di.NewBuilder()
	if err := builder.Add(di.Def{
		Name:  "links-storage",
		Build: func(ctn di.Container) (interface{}, error) { return linksStorage, nil }}); err != nil {
		return err
	}
	if err := builder.Add(di.Def{
		Name:  "users-storage",
		Build: func(ctn di.Container) (interface{}, error) { return usersStorage, nil }}); err != nil {
		return err
	}
	if err := builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return err
	}
	if err := builder.Add(di.Def{
		Name:  "secret-key",
		Build: func(ctn di.Container) (interface{}, error) { return []byte(cfg.SecretKey), nil }}); err != nil {
		return err
	}
	constans.GlobalContainer = builder.Build()
	return nil
}
