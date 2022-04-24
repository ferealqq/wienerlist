package boardapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestListBoardsHandler(t *testing.T) {
	appEnv := CreateContextForTestSetup()
	r, _ := http.NewRequest(http.MethodGet, "/boards", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	// recreate route from routes.go
	router.
		Methods(http.MethodGet).
		Path("/boards").
		Name("ListBoardsHandler").
		Handler(MakeHandler(appEnv, ListBoardsHandler))
	n := negroni.New()
	n.UseHandler(router)
	n.ServeHTTP(w, r)
	// test response headers and codes
	assert.Equal(t, http.StatusOK, w.Code, "they should be equal")
	assert.Equal(t, "GNU Terry Pratchett", w.HeaderMap["X-Clacks-Overhead"][0], "they should be equal")
	// parse json body
	var f interface{}
	json.Unmarshal(w.Body.Bytes(), &f)
}