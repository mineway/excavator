package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ermos/annotation/parser"
	"github.com/julienschmidt/httprouter"
	"github.com/mineway/excavator/internal/pkg/logger"
	"github.com/mineway/excavator/internal/pkg/manager"
	"github.com/mineway/excavator/internal/pkg/response"
	"github.com/mineway/excavator/utils"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var beforeDefaultMiddleware []string

var afterDefaultMiddleware []string

func SetBeforeDefaultMiddleware (name ...string) {
	beforeDefaultMiddleware = append(beforeDefaultMiddleware, name...)
}

func SetAfterDefaultMiddleware (name ...string) {
	afterDefaultMiddleware = append(afterDefaultMiddleware, name...)
}

func Serve (ch chan string, port, routeLocation string, controllerHandler, middlewareHandler interface{}) {
	router := httprouter.New()
	routes := getRoutes(routeLocation)

	for _, route := range routes {
		for _, r := range route.Routes {
			routePath := fmt.Sprintf("/api%s", r.Route)

			switch strings.ToLower(r.Method) {
			case "get":
				router.GET(routePath, call(route, controllerHandler, middlewareHandler))
			case "post":
				router.POST(routePath, call(route, controllerHandler, middlewareHandler))
			case "put":
				router.PUT(routePath, call(route, controllerHandler, middlewareHandler))
			case "patch":
				router.PATCH(routePath, call(route, controllerHandler, middlewareHandler))
			case "delete":
				router.DELETE(routePath, call(route, controllerHandler, middlewareHandler))
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

func call(route parser.API, handler interface{}, middlewareHandler interface{}) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	build := reflect.ValueOf(handler).MethodByName(route.Controller)
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		m := manager.New(r.Context(), r, route, ps)

		// Body Validator
		status, err := m.CheckRequest(r)
		if err != nil {
			response.Error(w, status, err)
			return
		}

		// Before Middleware
		status, err = callMiddleware(r.Context(), middlewareHandler, append(beforeDefaultMiddleware, route.Middleware.Before...), m, w, r)
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

		// After middleware
		status, err = callMiddleware(r.Context(), middlewareHandler, append(afterDefaultMiddleware, route.Middleware.After...), m, w, r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func getRoutes(location string) (annotations []parser.API){
	file, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(file, &annotations); err != nil {
		log.Fatal(err)
	}
	return
}

func callMiddleware(
	ctx context.Context,
	handler interface{},
	middlewares []string,
	m *manager.Manager,
	w http.ResponseWriter,
	r *http.Request,
) (s int, err error){
	for _, middleware := range middlewares {
		mw := reflect.ValueOf(handler).MethodByName(middleware)
		build := make([]reflect.Value, 4)

		build[0] = reflect.ValueOf(ctx)
		build[1] = reflect.ValueOf(m)
		build[2] = reflect.ValueOf(w)
		build[3] = reflect.ValueOf(r)

		res := mw.Call(build)

		if len(res) != 2 {
			return http.StatusInternalServerError, fmt.Errorf("%s's middleware don't return 2 arguments", middleware)
		}

		if res[1].Interface() != nil {
			return res[0].Interface().(int), res[1].Interface().(error)
		}
	}
	return 0, nil
}