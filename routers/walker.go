package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/walker"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunWalkerHandler(c *gin.Context) {
	result := walker.PutFlags()
	var data string
	for k, v := range result {
		data += fmt.Sprintf("%s %d\n", k, v)
	}
	c.Data(http.StatusOK, "text/plain", []byte(data))
}
