package web

import (
	"github.com/UQuark/housecontrol/strip"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var stripDashboardTemplate *template.Template

func init() {
	f, err := os.Open("./html/strip_dashboard.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	stripDashboardTemplate, err = template.New("strip_dashboard").Parse(string(buffer))
	if err != nil {
		panic(err)
	}
}

func selected(s bool) string {
	if s {
		return "selected"
	} else {
		return ""
	}
}

func (w *Web) HandleStripDashboard(ctx *gin.Context) {
	cmd := strip.NewCommandBuilder().
		Get().
		Width().
		Speed().
		Brightness().
		Mode().
		Build()

	stripConfig, err := w.Strip.Execute(cmd)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, "text/plaintext", []byte(err.Error()))
		return
	}
	log.Println(stripConfig)

	data := map[string]interface{}{
		"width": stripConfig[1],
		"speed": stripConfig[2],
		"brightness": stripConfig[3],
		"noiseSelected": selected(strip.Mode(stripConfig[4]) == strip.ModeNoise),
		"rainbowSelected": selected(strip.Mode(stripConfig[4]) == strip.ModeRainbow),
		"epilepticSelected": selected(strip.Mode(stripConfig[4]) == strip.ModeEpileptic),
		"turnoffSelected": selected(strip.Mode(stripConfig[4]) == strip.ModeTurnoff),
		"nightSelected": selected(strip.Mode(stripConfig[4]) == strip.ModeNight),
	}

	err = stripDashboardTemplate.Execute(ctx.Writer, data)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, "text/plaintext", []byte(err.Error()))
		return
	}
	ctx.Status(http.StatusOK)
}

func (w *Web) HandleStripUpdate(ctx *gin.Context) {
	data := make(map[string]byte)
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plaintext", []byte(err.Error()))
		return
	}
	cmd := strip.NewCommandBuilder().Set()
	for k, v := range data {
		switch k {
		case "width":
			cmd.Width(v)
		case "speed":
			cmd.Speed(v)
		case "brightness":
			cmd.Brightness(v)
		case "mode":
			cmd.Mode(strip.Mode(v))
		}
	}
	buffer, err := w.Strip.Execute(cmd.Build())
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/plaintext", []byte(err.Error()))
		return
	}
	log.Println(buffer)
}