package main

import (
	"n_communication/server"

	"n_communication/handler"

	"go.uber.org/zap"
)

func initLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()

	s := server.New()

	hh := handler.NewHealthHandler()
	s.Mount("/", hh.NewHealthRouter())

	ch := handler.NewCommHandler()
	s.Mount("/comms", ch.NewCommRouter())

	s.StartServer(":8086")
}
