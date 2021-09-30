package routers

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
)

func AddVpnTeam(team *models.Team) error {
	vpnAddr := os.Getenv("OVPN_ADMIN")
	if vpnAddr == ""{
		vpnAddr = "http://localhost:9000"
	}
	urlAddr := vpnAddr + "/api/user/create"
	_, httpErr := http.PostForm(urlAddr, url.Values{
		"username": {team.Name},
		"password": {"kb4ctf"},
	})

	if httpErr != nil{
		log.Println(httpErr)
		return httpErr
	}
	return nil
}

func CreateVpnTeams(c *gin.Context){
	teams, dbErr := database.GetTeams()
	if dbErr != nil{
		c.JSON(http.StatusBadRequest, dbErr)
		return
	}
	for _, team := range teams{
		AddVpnTeam(team)
	}
}