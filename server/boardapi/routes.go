package boardapi

// Route is the model for the router setup
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

// Routes are the main setup for our Router
type Routes []Route

var routes = Routes{
	Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},
	// === BOARDS ===
	Route{"ListBoards", "GET", "/boards", ListBoardsHandler},
	Route{"CreateBoard", "POST", "/boards", CreateBoardHandler},
	Route{"GetBoard", "GET", "/boards/{bid:[0-9]+}", GetBoardHandler},
	Route{"UpdateBoard", "PUT", "/boards/{bid:[0-9]+}", UpdateBoardHandler},
	Route{"DeleteBoard", "DELETE", "/boards/{bid:[0-9]+}", DeleteBoardHandler},
}
