package controllers

import (
	"NetTracer/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDatabases(c *gin.Context) {
	c.String(http.StatusOK, util.GetArbitraryDB())
}
