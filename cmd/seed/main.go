package main

import (
	"github.com/ic3network/mccs-alpha/global"
	"github.com/ic3network/mccs-alpha/internal/seed"
)

func main() {
	global.Init()
	seed.LoadData()
	seed.Run()
}
