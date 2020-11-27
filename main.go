package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var modes = map[string]byte{
	"RAINBOW":        0,
	"FLASHLIGHT":     1,
	"EPILEPTIC":      2,
	"WALKINGPIXEL":   3,
	"STARS":          4,
	"SORT":           5,
	"STACK":          6,
	"SMOOTHCHANGING": 7,
	"NIGHT":          8,
	"TURNOFF":        9,
}

var config struct {
	LEDTTY string `json:"led_tty"`
}

var runtime struct {
	ledTTY *os.File
}

func sendByteToLED(b byte) {
	fmt.Fprintf(runtime.ledTTY, "%d", b)
}

func setLEDMode(responseWriter http.ResponseWriter, request *http.Request) {
	mode := request.URL.Query().Get("mode")
	if mode == "" {
		http.Error(responseWriter, "No mode specified", http.StatusBadRequest)
	}
	modeByte, ok := modes[mode]
	if !ok {
		http.Error(responseWriter, "Invalid mode", http.StatusBadRequest)
	}
	sendByteToLED(modeByte)
}

func readConfig() error {
	file, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()
	err = json.Unmarshal(text, &config)
	if err != nil {
		return err
	}
	return nil
}

func prepareLED() error {
	var err error
	runtime.ledTTY, err = os.OpenFile(config.LEDTTY, os.O_WRONLY, os.ModeDevice)
	return err
}

func main() {
	err := readConfig()
	if err != nil {
		panic(err)
	}
	log.Println("Config read")
	err = prepareLED()
	if err != nil {
		panic(err)
	}
	log.Println("LED TTY opened")

	http.Handle("/", http.FileServer(http.Dir("./html")))
	log.Println("File server initialized")
	http.HandleFunc("/set_led_mode", setLEDMode)
	log.Println("Led mode handler set")

	log.Println("Listening")
	http.ListenAndServe(":8080", nil)
}
