package boardapi

import (
	"log"
	"time"

	a "github.com/ferealqq/wienerlist/pkg/appenv"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(appEnv a.AppEnv) {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Access-Control-Allow-Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-CSRF-Token", "X-Powered-By"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	apiRouter := router.Group("/api/v1")
	WorkspaceRouter(apiRouter, appEnv)
	BoardRouter(apiRouter, appEnv)
	SectionRouter(apiRouter, appEnv)
	ItemRouter(apiRouter, appEnv)
	isDevelopment := appEnv.Env == "LOCAL"
	if !isDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}

	secureMiddleware := secure.New(secure.Options{
		// This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains
		// options to be ignored during development. When deploying to production,
		// be sure to set this to false.
		IsDevelopment: isDevelopment,
		// AllowedHosts is a list of fully qualified domain names that are allowed (CORS)
		AllowedHosts: []string{},
		// If ContentTypeNosniff is true, adds the X-Content-Type-Options header
		// with the value `nosniff`. Default is false.
		ContentTypeNosniff: true,
		// If BrowserXssFilter is true, adds the X-XSS-Protection header with the
		// value `1; mode=block`. Default is false.
		BrowserXssFilter: true,
	})
	// start now
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	startupMessage := "===> Starting app (v" + appEnv.Version + ")"
	startupMessage = startupMessage + " on port " + appEnv.Port
	startupMessage = startupMessage + " in " + appEnv.Env + " mode."
	log.Println(startupMessage)
	if appEnv.Env == "LOCAL" {
		n.Run("localhost:" + appEnv.Port)
	} else {
		n.Run(":" + appEnv.Port)
	}
}
