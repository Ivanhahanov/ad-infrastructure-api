package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"time"
)

type Team struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	SshPubKey string `json:"ssh_pub_key"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetTeamInfo(c *gin.Context) {
	team := Team{
		Name: "Test",
	}
	c.JSON(http.StatusOK, team)
}

func generateIp(number int) string {
	ip, _, err := net.ParseCIDR(config.Conf.Network)
	if err != nil {
		log.Println(err.Error())
	}
	ip = ip.To4()
	ip[2] += byte(number)
	return ip.String()
}

func CreateTeam(c *gin.Context) {
	var team Team
	jsonErr := c.BindJSON(&team)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		log.Println(dbErr)
	}
	ipAddress := generateIp(len(teams))
	hash, hashErr := HashPassword(team.Password)
	if hashErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": hashErr.Error()})
		return
	}

	dbTeam := &models.Team{
		ID:        primitive.NewObjectID(),
		Name:      team.Name,
		Address:   ipAddress,
		Hash:      hash,
		SshPubKey: team.SshPubKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dbErr = database.CreateTeam(dbTeam)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}

	//sshErr := CreateSshKeyFile(team.Name, team.SshPubKey)
	//if sshErr != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": sshErr.Error()})
	//	return
	//}
	//log.Println("ssh key created for team", team.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("The team %s created", team.Name),
	})
}

func DeleteTeam(c *gin.Context) {
	teamName := c.Param("name")
	dbErr := database.DeleteTeam(teamName)
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", teamName),
	})
}
