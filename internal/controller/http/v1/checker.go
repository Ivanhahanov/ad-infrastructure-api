package v1

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"

	_manager "github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	_walker "github.com/Ivanhahanov/ad-infrastructure-api/pkg/walker"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type checkerRoutes struct {
	l         logger.Interface
	manager  *_manager.CtfManager
	walker   *_walker.Walker
}
func NewCheckerRoutes(handler *gin.RouterGroup, l logger.Interface, manager *_manager.CtfManager, walker *_walker.Walker, pass string) {
	r := &checkerRoutes{l, manager, walker}

	handler.GET("/checker",
		gin.BasicAuth(gin.Accounts{
			"checker": pass,
		}),
		r.CheckerHandler)
}

func (r *checkerRoutes) CheckerHandler(c *gin.Context) {

	log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA================================")
	checkResult, err := r.walker.CheckFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	var data string
	for k, v := range checkResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}
	r.manager.RemoveAllFlags()
	r.manager.WriteTime()
	putResult, err := r.walker.PutFlags()
	if err != nil {
		log.Println(err)
		c.Data(http.StatusOK, "text/plain", []byte(""))
	}
	for k, v := range putResult {
		data += fmt.Sprintf("%s %d\n", k, v)
	}

	c.Data(http.StatusOK, "text/plain", []byte(data))
}
