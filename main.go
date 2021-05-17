package main

import (
	"github.com/braineo/genshin-wish/server"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func main() {
	s := server.New()
	s.Run()
}
