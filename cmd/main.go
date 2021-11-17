package main

import (
	"github.com/AliceDiNunno/go-image-database/src/adapters/events"
	"github.com/AliceDiNunno/go-image-database/src/adapters/persistence/postgres"
	"github.com/AliceDiNunno/go-image-database/src/adapters/rest"
	"github.com/AliceDiNunno/go-image-database/src/adapters/storage/localStorage"
	"github.com/AliceDiNunno/go-image-database/src/config"
	"github.com/AliceDiNunno/go-image-database/src/core/usecases"
	glc "github.com/AliceDiNunno/go-logger-client"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()

	ginConfiguration := config.LoadGinConfiguration()
	dbConfig := config.LoadGormConfiguration()
	logConfig := config.LoadLogConfiguration()
	storageConfig := config.LoadStorageConfig()
	fileStorage := localStorage.NewLocalStorage(storageConfig)

	glc.SetupHook(logConfig)

	var albumRepo usecases.AlbumRepo
	var pictureRepo usecases.PictureRepo
	var tagRepo usecases.TagRepo

	var db *gorm.DB
	if dbConfig.Engine == "POSTGRES" {
		db = postgres.StartGormDatabase(dbConfig)
		albumRepo = postgres.NewAlbumRepo(db)
		pictureRepo = postgres.NewPictureRepo(db)
		tagRepo = postgres.NewTagRepo(db)

		db.AutoMigrate(&postgres.Tag{})
		db.AutoMigrate(&postgres.Album{})
		db.AutoMigrate(&postgres.Picture{})
	}

	usecasesHandler := usecases.NewInteractor(albumRepo, pictureRepo, tagRepo, fileStorage)

	restServer := rest.NewServer(ginConfiguration)
	routesHandler := rest.NewRouter(usecasesHandler)
	eventManager := events.NewEventManager(usecasesHandler)

	rest.SetRoutes(restServer.Router, routesHandler)

	go eventManager.Start()
	restServer.Start()
}
