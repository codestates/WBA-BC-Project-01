package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() (Controller, error) {
	return Controller{}, nil
}

func (cc *Controller) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{"isLogined": false})
}
