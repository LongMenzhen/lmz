package main

import (
	_ "github.com/cyrnicolase/lmz/config"
	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/route"
)

func main() {
	hub := engine.AttachHub()
	hub.Run()

	route.Route()
}
