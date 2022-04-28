package boardapi

import (
	"log"

	. "github.com/ferealqq/golang-trello-copy/server/pkg/container"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(appContainer AppContainer) {
	router := gin.Default()
	BoardRouter(router, appContainer)
	isDevelopment := appContainer.Env == "LOCAL"
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
	startupMessage := "===> Starting app (v" + appContainer.Version + ")"
	startupMessage = startupMessage + " on port " + appContainer.Port
	startupMessage = startupMessage + " in " + appContainer.Env + " mode."
	log.Println(startupMessage)
	if appContainer.Env == "LOCAL" {
		n.Run("localhost:" + appContainer.Port)
	} else {
		n.Run(":" + appContainer.Port)
	}
}
