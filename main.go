package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

func createCmd(request bool, key byte, value byte) []byte {
	key &= 0b111
	var requestByte byte
	if request {
		requestByte = 1 << 3
	} else {
		requestByte = 0
	}
	part1 := requestByte | key
	part2 := value
	return []byte{part1, part2}
}

func setLEDMode(responseWriter http.ResponseWriter, request *http.Request) {
	const keyMode byte = 4
	mode := request.URL.Query().Get("mode")
	if mode == "" {
		http.Error(responseWriter, "No mode specified", http.StatusBadRequest)
	}
	modeByte, ok := modes[mode]
	if !ok {
		http.Error(responseWriter, "Invalid mode", http.StatusBadRequest)
	}
	runtime.ledTTY.Write(createCmd(false, keyMode, modeByte))
}

func setLEDBrightness(responseWriter http.ResponseWriter, request *http.Request) {
	const keyBrightness byte = 3
	brightness := request.URL.Query().Get("brightness")
	if brightness == "" {
		http.Error(responseWriter, "No brightness specified", http.StatusBadRequest)
	}
	brightnessInt, err := strconv.Atoi(brightness)
	if err != nil {
		http.Error(responseWriter, "Invalid brightness", http.StatusBadRequest)
	}
	if brightnessInt > 255 {
		brightnessInt = 255
	}
	if brightnessInt < 0 {
		brightnessInt = 0
	}
	brightnessByte := byte(brightnessInt)
	runtime.ledTTY.Write(createCmd(false, keyBrightness, brightnessByte))
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
	http.HandleFunc("/set_led_brightness", setLEDBrightness)
	log.Println("Led mode handler set")

	log.Println("Listening")
	http.ListenAndServe(":8080", nil)
}
