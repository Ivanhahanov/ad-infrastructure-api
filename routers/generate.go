package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"net/http"
	"os"
)

func CreateFile(filename string) {

}

type Teams struct {
	Name string
	Img  string
}

func GenerateTerraformConfig(c *gin.Context) {
	filename := "main.tf"
	requiredVersion := ">= 0.13"
	providerUri := "qemu:///system"
	osImageFilename := "focal-server-cloudimg-amd64.img"
	teams := []Teams{
		{
			Name: "naliway",
			Img:  osImageFilename,
		},
		{
			Name: "nakateam",
			Img:  osImageFilename,
		},
	}
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	// generate os_image variable
	tfBlock := rootBody.AppendNewBlock("terraform", nil)
	tfBlockBody := tfBlock.Body()
	tfBlockBody.SetAttributeValue("required_version", cty.StringVal(requiredVersion))
	reqProvsBlock := tfBlockBody.AppendNewBlock("required_providers", nil)
	reqProvsBlockBody := reqProvsBlock.Body()

	reqProvsBlockBody.SetAttributeValue("libvirt", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("dmacvicar/libvirt"),
		"version": cty.StringVal("0.6.11"),
	}))
	hosts := rootBody.AppendNewBlock("provider", []string{"libvirt"})
	hostsBody := hosts.Body()
	hostsBody.SetAttributeValue("uri", cty.StringVal(providerUri))

	interfaceVar := rootBody.AppendNewBlock("resource", []string{"libvirt_pool", "os_pools"})
	interfaceBody := interfaceVar.Body()
	interfaceBody.SetAttributeValue("name", cty.StringVal("vm"))
	interfaceBody.SetAttributeValue("type", cty.StringVal("dir"))
	interfaceBody.SetAttributeValue("path", cty.StringVal("./data"))

	for _, team := range teams {
		volume := rootBody.AppendNewBlock("resource", []string{"libvirt_volume", "os-qcow2"})
		volumeBody := volume.Body()
		volumeBody.SetAttributeValue("name", cty.StringVal(team.Name+"-qcow2"))
		volumeBody.SetAttributeValue("pool", cty.StringVal("vm"))
		volumeBody.SetAttributeValue("source", cty.StringVal(team.Img))
		volumeBody.SetAttributeValue("format", cty.StringVal("qcow2"))


		cloudInit := rootBody.AppendNewBlock("resource", []string{"libvirt_cloudinit_disk", "commoninit"})
		cloudInitBody := cloudInit.Body()
		cloudInitBody.SetAttributeValue("name", cty.StringVal(team.Name+"-commoninit.iso"))
		cloudInitBody.SetAttributeValue("pool", cty.StringVal("vm"))
		cloudInitBody.SetAttributeValue("source", cty.StringVal(team.Img))
		cloudInitBody.SetAttributeValue("user_data", cty.StringVal("qcow2"))

	}
	// save file
	fmt.Printf("%s", hclFile.Bytes())
	tfFile.Write(hclFile.Bytes())
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("file %s generated", filename),
	})
}
