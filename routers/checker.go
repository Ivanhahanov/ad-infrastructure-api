package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/walker"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CheckerHandler(c *gin.Context) {
	checkResult, err := walker.CheckFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	var data string
	for k, v := range checkResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}
	database.RemoveAllFlags()
	database.WriteTime()
	putResult, err := walker.PutFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	for k, v := range putResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}

	c.Data(http.StatusOK, "text/plain", []byte(data))
}
