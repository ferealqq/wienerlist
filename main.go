package main

import (
	"os"
	"strings"

	api "github.com/ferealqq/golang-trello-copy/server/boardapi"
	"github.com/ferealqq/golang-trello-copy/server/migrations"
	appenv "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	database "github.com/ferealqq/golang-trello-copy/server/pkg/database"
	vparse "github.com/ferealqq/golang-trello-copy/server/pkg/version"
	log "github.com/sirupsen/logrus"
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
	if err := database.InitDB(); err != nil {
		// FIXME do something?
		panic(err)
	}
	if env == LOCAL {
		migrations.MigrateSeedAfterwards(
			database.DBConn)
	} else {
		err = migrations.Migrate(database.DBConn)
		if err != nil {
			// TODO we probably should not panic here. We should log it and try again?
			panic(err)
		}
	}
	// ===========================================================================
	// Initialise application context
	// ===========================================================================
	appEnv := appenv.AppEnv{
		Version: version,
		Env:     env,
		Port:    port,
	}
	// ===========================================================================
	// Start application
	// ===========================================================================
	api.StartServer(appEnv)

	defer database.Close()
}
