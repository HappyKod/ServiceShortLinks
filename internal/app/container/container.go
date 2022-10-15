package container

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/Linksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/Linksstorage/FilesLinksStorage"
	"HappyKod/ServiceShortLinks/internal/storage/Linksstorage/MemLinksStorage"
	"HappyKod/ServiceShortLinks/internal/storage/UsersStorage"
	"HappyKod/ServiceShortLinks/internal/storage/UsersStorage/MemUsersStorage"
	"github.com/sarulabs/di"
	"log"
)

func BuildContainer(cfg models.Config) error {
	var linksStorage Linksstorage.LinksStorages
	if cfg.FileStoragePATH != "" {
		store, err := FilesLinksStorage.New(cfg.FileStoragePATH)
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован file-Linksstorage")
	} else {
		store, err := MemLinksStorage.New()
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован MemLinksStorage")
	}
	var usersStorage UsersStorage.UsersStorage
	usersStorage, err := MemUsersStorage.New()
	if err != nil {
		return err
	}
	log.Println("Задействован MemUsersStorage")
	builder, _ := di.NewBuilder()
	if err := builder.Add(di.Def{
		Name:  "Linksstorage",
		Build: func(ctn di.Container) (interface{}, error) { return linksStorage, nil }}); err != nil {
		return err
	}
	if err := builder.Add(di.Def{
		Name:  "UsersStorage",
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
