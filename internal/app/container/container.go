// Package container сборка DI-контейнера.
package container

import (
	"log"

	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage/fileslinksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage/memlinksstorage"
	"HappyKod/ServiceShortLinks/internal/storage/linksstorage/pglinkssotorage"

	"github.com/sarulabs/di"
)

// BuildContainer собирает в DI контейнер.
func BuildContainer(cfg models.Config) error {
	var linksStorage linksstorage.LinksStorages
	if cfg.DataBaseURL != "" {
		store, err := pglinkssotorage.New(cfg.DataBaseURL)
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован pg-linkssotorage")
	} else if cfg.FileStoragePATH != "" {
		store, err := fileslinksstorage.New(cfg.FileStoragePATH)
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован file-linksstorage")
	} else {
		store, err := memlinksstorage.New()
		if err != nil {
			return err
		}
		linksStorage = store
		log.Println("Задействован mem-linksstorage")
	}
	err := linksStorage.Ping()
	if err != nil {
		return err
	}
	builder, _ := di.NewBuilder()
	if err := builder.Add(di.Def{
		Name:  "linksstorage",
		Build: func(ctn di.Container) (interface{}, error) { return linksStorage, nil }}); err != nil {
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
