package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/spql"
)

// HTTPServer holds HTTP server configuration and routes
type HTTPServer struct {
	Config        *db.SparrowConfig
	mux           *http.ServeMux
	dbManager     *db.DBManager
	routers       map[string]*controllerInfo
	queryExecutor *spql.QueryExecutor
}

type controllerInfo struct {
	route      string
	httpMethod string
	method     func(request *RequestData)
}

func (httpServer *HTTPServer) add(c *controllerInfo) {
	parts := strings.Split(c.route[1:], "/")
	httpServer.routers[parts[0]] = c
}

func (httpServer *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")

	if parts[0] == "favicon.ico" {
		return
	}

	monitor.IncHTTPRequests()

	controller, ok := httpServer.routers[parts[0]]

	if ok {
		parts := strings.Split(r.URL.Path[1:], "/")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		controller.method(&RequestData{responseWriter: w, request: r, params: parts[1:]})
	}

}

// Start starts HTTP server listener
func (httpServer *HTTPServer) Start() {
	log.Printf("Starting HTTP Server %s:%s", httpServer.Config.HTTPHost, httpServer.Config.HTTPPort)

	handler := NewServeHandler(httpServer.dbManager, httpServer.queryExecutor)

	httpServer.add(&controllerInfo{route: "/g", httpMethod: "GET", method: handler.get})
	httpServer.add(&controllerInfo{route: "/upload", httpMethod: "POST", method: handler.upload})
	httpServer.add(&controllerInfo{route: "/query", httpMethod: "POST", method: handler.serveQuery})

	httpServer.mux.Handle("/", httpServer)

	http.ListenAndServe(":"+httpServer.Config.HTTPPort, httpServer.mux)
}

// Stop stops HTTP server listener
func (httpServer *HTTPServer) Stop() {
	log.Printf("Stopping HTTP Server")
}

// NewHTTPServer returns new HTTPServer
func NewHTTPServer(config *db.SparrowConfig, dbm *db.DBManager) HTTPServer {
	return HTTPServer{
		Config:        config,
		dbManager:     dbm,
		queryExecutor: spql.NewQueryExecutor(dbm),
		mux:           http.NewServeMux(),
		routers:       make(map[string]*controllerInfo),
	}
}
