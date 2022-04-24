package main

import (
	"os"
	"strings"

	api "github.com/ferealqq/golang-trello-copy/server/boardapi"
	"github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	"github.com/ferealqq/golang-trello-copy/server/migrations"
	database "github.com/ferealqq/golang-trello-copy/server/pkg/database"
	vparse "github.com/ferealqq/golang-trello-copy/server/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

var (
	LOCAL = "LOCAL"
)

func init() {
	if LOCAL == strings.ToUpper(os.Getenv("ENV")) {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	// ===========================================================================
	// Load environment variables
	// ===========================================================================
	var (
		env     = strings.ToUpper(os.Getenv("ENV")) // LOCAL, DEV, STG, PRD
		port    = os.Getenv("PORT")                 // server traffic on this port
		version = os.Getenv("VERSION")              // path to VERSION file
	)
	// ===========================================================================
	// Read version information
	// ===========================================================================
	version, err := vparse.ParseVersionFile(version)
	if err != nil {
		log.WithFields(log.Fields{
			"env":  env,
			"err":  err,
			"path": os.Getenv("VERSION"),
		}).Fatal("Can't find a VERSION file")
		return
	}
	log.WithFields(log.Fields{
		"env":     env,
		"path":    os.Getenv("VERSION"),
		"version": version,
	}).Info("Loaded VERSION file")
	// ===========================================================================
	// Initialise data storage
	// ===========================================================================
	database.InitDB()
	if env == LOCAL{
		migrations.MigrateSeedAfterwards(database.DBConn)
	}else{
		migrations.Migrate(database.DBConn)
	}
	// ===========================================================================
	// Initialise application context
	// ===========================================================================
	appEnv := api.AppEnv{
		Render:    	render.New(),
		Version:   	version,
		Env:       	env,
		Port:      	port,
		DBConn: 	database.DBConn,
	}
	// ===========================================================================
	// Start application
	// ===========================================================================
	api.StartServer(appEnv)

	defer database.Close()
}

func seed(){
	board := models.Board{
		Title: "Test Board",
		Description: "Test Board Description",
	}
	result := database.DBConn.Create(&board)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"err": result.Error,
		}).Fatal("Can't create board")
		return
	}
}