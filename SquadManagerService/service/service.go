package service

import (
	"net/http"
	"time"

	"log"

	"github.com/julienschmidt/httprouter"
)

type Configuration struct {
	DatabaseName string
	Host         string
	DbTimeout    time.Duration
}

type MainHandler struct {
	context *Context
	router  *httprouter.Router
}

func (mainHandler *MainHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	LogRequest(mainHandler.router.ServeHTTP)(writer, request)
}

func (mainHandler *MainHandler) Close() {
	mainHandler.context.Close()
}

func MakeMainHandler(config Configuration) *MainHandler {
	context, err := newContext(config)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	router.GET("/squad", context.with(Handler(listSquads)))
	router.POST("/squad", context.with(NoInputHandler(createSquad)))
	router.GET("/squad/:id", context.with(SquadHandler(getSquad)))
	router.POST("/squad/:id", context.with(SquadHandler(postSquadMember)))

	return &MainHandler{context, router}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (writer *loggingResponseWriter) WriteHeader(code int) {
	writer.statusCode = code
	writer.ResponseWriter.WriteHeader(code)
}

func LogRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		loggingWriter := &loggingResponseWriter{ResponseWriter: writer, statusCode: http.StatusOK}
		next(loggingWriter, request)
		duration := time.Now().Sub(start)
		log.Println(loggingWriter.statusCode, duration, request.URL.String())
	}
}
