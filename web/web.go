package web

import (
	"github.com/UQuark/housecontrol/strip"
	"github.com/gin-gonic/gin"
)

type Web struct {
	Bind string
	Strip strip.Strip `json:"-"`

	router *gin.Engine
}

func (w *Web) Initialize() {
	w.router = gin.New()

	strip := w.router.Group("/strip")
	strip.GET("/dashboard", w.HandleStripDashboard)
	strip.PUT("/update", w.HandleStripUpdate)
	strip.PUT("/reset", w.HandleStripReset)
}

func (w *Web) Run() error {
	return w.router.Run(w.Bind)
}