package server

import (
	"context"
	"fmt"
	"github.com/ermos/annotation/parser"
	"github.com/julienschmidt/httprouter"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/manager"
	"github.com/mineway/excavator/internal/pkg/response"
	"github.com/mineway/excavator/utils"
	"github.com/rs/cors"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func Serve (ch chan string, port string, routes []parser.API, controllerHandler interface{}) {
	router := httprouter.New()

	for _, route := range routes {
		for _, r := range route.Routes {
			routePath := fmt.Sprintf("/api%s", r.Route)

			switch strings.ToLower(r.Method) {
			case "get":
				router.GET(routePath, call(route, controllerHandler))
			case "post":
				router.POST(routePath, call(route, controllerHandler))
			case "put":
				router.PUT(routePath, call(route, controllerHandler))
			case "patch":
				router.PATCH(routePath, call(route, controllerHandler))
			case "delete":
				router.DELETE(routePath, call(route, controllerHandler,))
			}
		}
	}

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET","PUT","POST","DELETE","PATCH","OPTIONS"},
		AllowCredentials: false,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	server := initServer(port, c.Handler(router))

	logger.Success("[API] Currently running on port \u001B[1m%s\u001B[0m..", port)

	for {
		select {
		case v := <- ch:
			switch v {
			case "stop":
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					logger.Fatal("[API] %s", err.Error())
				}
				return
			}
		}
	}
}

func initServer(port string, handler http.Handler) *http.Server {
	if !utils.IsAvailablePort(port) {
		newPort, err := utils.GetAvailablePort()
		if err != nil {
			logger.Fatal("[API] cannot fetch available port (%s)", err.Error())
		}

		logger.Warning("[API] :%s is already in use, new port is :%s, don't forget to change it into your web interface", port, newPort)

		port = newPort
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: handler}

	go func(server *http.Server) {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal("[API] %s", err.Error())
		}
	}(server)

	return server
}

func call(route parser.API, handler interface{}) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	build := reflect.ValueOf(handler).MethodByName(route.Controller)
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		m, status, err := manager.New(r, route, ps)
		if err != nil {
			response.Error(w, status, err)
			return
		}

		// Controller
		rebuild := make([]reflect.Value, 3)
		rebuild[0] = reflect.ValueOf(r.Context())
		rebuild[1] = reflect.ValueOf(m)
		rebuild[2] = reflect.ValueOf(w)
		_ = build.Call(rebuild)
	}
}