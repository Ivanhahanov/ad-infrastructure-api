package admin

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func (r *routes) AddVpnTeam(team *entity.Team) error {
	vpnAddr := utils.GetEnv("OVPN_ADMIN", "http://localhost:9000")

	_, err := http.PostForm(vpnAddr + "/api/user/create", url.Values{
		"username": {team.Name},
		"password": {"kb4ctf"},
	})

	if err == nil{return nil}

	return err
}

func (r *routes) CreateVpnTeams(c *gin.Context) {
	teams, err := r.teamRepo.GetTeams()
	if err == nil {
		for _, team := range teams {
			r.AddVpnTeam(team)
		}
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}
