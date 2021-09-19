package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
)

func CreateSshKeyFile(name string, key string) error {
	fileName := fmt.Sprintf("%s.pub", name)
	filePath := path.Join(config.Conf.TerraformProjectPath, config.Conf.SshKeys, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, writeErr := f.Write([]byte(key))
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func GetTeamInfo(c *gin.Context) {
	team := models.Team{
		Name: "Test",
		Players: []string{
			"test1",
			"test2",
			"test3",
		},
	}
	c.JSON(http.StatusOK, team)
}
func CreateTeam(c *gin.Context) {
	var team models.Team
	jsonErr := c.BindJSON(&team)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	sshErr := CreateSshKeyFile(team.Name, team.SshPubKey)
	if sshErr != nil {
		log.Println(sshErr)
	}
	log.Println("ssh key created for team", team.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("The team %s created", team.Name),
	})
}

func DeleteTeam(c *gin.Context) {
	teamName := "Test"
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", teamName),
	})
}
