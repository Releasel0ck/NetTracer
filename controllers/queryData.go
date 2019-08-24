package controllers

import (
	"NetTracer/models"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func QueryData(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	result := QueryDataF(string(data))
	c.String(http.StatusOK, result)

}

func QueryDataF(queryString string) string {
	if strings.Contains(queryString, "Port") {
		port := strings.Split(queryString, ":")[1]
		new_output := models.QueryByPort(port)
		return new_output
	} else if strings.Contains(queryString, "ServiceIP") {
		ip := strings.Split(queryString, ":")[1]
		new_output := models.QueryBySIP(ip)
		return new_output
	} else if strings.Contains(queryString, "ClientIP") {
		ip := strings.Split(queryString, ":")[1]
		new_output := models.QueryByCIP(ip)
		return new_output
	} else if strings.Contains(queryString, "IPAddress") {
		ip := strings.Split(queryString, ":")[1]
		new_output := models.QueryByIP(ip)
		return new_output
	} else {
		return ""
	}
}
