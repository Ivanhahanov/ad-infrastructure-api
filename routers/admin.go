package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func UsersList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []string{"test"},
	})
}

func TeamsList(c *gin.Context) {
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

func DeleteUsers(c *gin.Context) {
	user := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("user %s deleted", user),
	})
}

func DeleteTeams(c *gin.Context) {
	team := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("team %s deleted", team),
	})
}

type TeamsForAnsible struct {
	IP        string `json:"ip"`
	Netmask      string `json:"netmask"`
	Mode      string `json:"mode"`
	Name      string `json:"name"`
	DHCPStart string `json:"dhcp_start"`
	DHCPEnd   string `json:"dhcp_end"`
}

func CountTeamsHandler(c *gin.Context) {
	var result []TeamsForAnsible
	teams, _ := database.GetTeams()
	for _, team := range teams{
		data := TeamsForAnsible{
			Name: team.Name,
			Netmask: "255.255.255.0",
			Mode: "nat",
		}
		// network address
		ip := net.ParseIP(team.Address)
		ip = ip.To4()
		ip[3] = 0
		data.IP = ip.String()
		// dhcp start
		ip[3] = 11
		data.DHCPStart = ip.String()
		ip[3] = 253
		data.DHCPEnd = ip.String()
		result = append(result, data)
	}

	c.JSON(http.StatusOK, gin.H{"teams": result})
}
