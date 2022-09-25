package container

import (
	"HappyKod/ServiceShortLinks/internal/constans"
	"HappyKod/ServiceShortLinks/internal/models"
	"HappyKod/ServiceShortLinks/internal/storage"
	"HappyKod/ServiceShortLinks/internal/storage/filestorage"
	"HappyKod/ServiceShortLinks/internal/storage/memstorage"
	"errors"
	"fmt"
	"github.com/sarulabs/di"
	"log"
)

func BuildContainer(cfg models.Config) error {
	var meStorage storage.Storages
	if cfg.FileStoragePATH != "" {
		store, err := filestorage.New(cfg.FileStoragePATH)
		if err != nil {
			return err
		}
		meStorage = store
		log.Println("Задействован filestorage")
	} else {
		store, err := memstorage.New()
		if err != nil {
			return err
		}
		meStorage = store
		log.Println("Задействован memstorage")
	}
	builder, _ := di.NewBuilder()
	if err := builder.Add(di.Def{
		Name:  "links-storage",
		Build: func(ctn di.Container) (interface{}, error) { return meStorage, nil }}); err != nil {
		return errors.New(fmt.Sprint("Ошибка инициализации контейнера", err))
	}
	if err := builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return errors.New(fmt.Sprint("Ошибка инициализации контейнера", err))
	}
	constans.GlobalContainer = builder.Build()
	return nil
}
