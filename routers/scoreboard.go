package routers

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowScoreboard(c *gin.Context) {
	sla := 1
	database.GetMetricsNames("naliway", "ExampleTask")
	c.JSON(http.StatusOK, sla)
}
