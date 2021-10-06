package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"net/http"
	"os"
	"path"
)

func (r *routes) createSshKeyFile(name string, key string) error {
	fileName := fmt.Sprintf("%s.pub", name)
	dirPath := path.Join(r.cfg.TerraformProjectPath, r.cfg.SshKeys)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {return err}

	filePath := path.Join(dirPath, fileName)

	f, err := os.Create(filePath)
	if err != nil {return err}
	defer f.Close()


	if _, err := f.Write([]byte(key)); err != nil {return err}

	return nil
}
func (r *routes) GenerateSshKeysDir(c *gin.Context) {
	teams, err := r.teamRepo.GetTeams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, team := range teams {
		sshKeyErr := r.createSshKeyFile(team.Name, team.SshPubKey)
		if sshKeyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": sshKeyErr.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("dir %s generated", r.cfg.SshKeys),
	})
}

func (r *routes) GenerateVariables(c *gin.Context) {
	fileName := "teams.tf"
	filePath := path.Join(r.cfg.TerraformProjectPath, fileName)

	// osImageFilename := "focal-server-cloudimg-amd64.img"
	teams, err := r.teamRepo.GetTeams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tfFile.Close()

	// initialize the body of the new file object
	rootBody := hclFile.Body()

	// generate os_image variable
	vmBlock := rootBody.AppendNewBlock("variable", []string{"teams"})
	vmBlockBody := vmBlock.Body()

	var teamsList []cty.Value
	var ipsList []cty.Value
	for _, team := range teams {
		teamsList = append(teamsList, cty.StringVal(team.Name))
		ipsList = append(ipsList, cty.StringVal(fmt.Sprintf(team.Address)))
	}

	vmBlockBody.SetAttributeValue("default", cty.ListVal(teamsList))

	ips := rootBody.AppendNewBlock("variable", []string{"ips"})
	ipsBody := ips.Body()
	ipsBody.SetAttributeValue("default", cty.ListVal(ipsList))

	hclFile.WriteTo(tfFile)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("file %s generated", fileName),
	})
}

func (r *routes) GeneratePrometheus(c *gin.Context) {

}
