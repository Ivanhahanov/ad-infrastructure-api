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
	result, err := walker.CheckFlags()
	if err != nil{
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	var data string
	for k, v := range result {
		data += fmt.Sprintf("%s %d\n", k, v)
	}
	database.RemoveAllFlags()
	c.Data(http.StatusOK, "text/plain", []byte(data))
}
