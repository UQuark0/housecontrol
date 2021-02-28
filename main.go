package main

import (
	"encoding/json"
	"github.com/UQuark/housecontrol/strip"
	"github.com/UQuark/housecontrol/web"
	"os"
)

var runtime = struct {
	Strip strip.Strip
	Web web.Web
}{}

func main() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&runtime)
	if err != nil {
		panic(err)
	}

	err = runtime.Strip.Initialize()
	if err != nil {
		panic(err)
	}

	runtime.Web.Strip = runtime.Strip

	runtime.Web.Initialize()
	err = runtime.Web.Run()
	if err != nil {
		panic(err)
	}
}
