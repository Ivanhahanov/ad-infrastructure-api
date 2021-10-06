package admin

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager/repositories"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type routes struct {
	l         logger.Interface
	cfg      *config.Config
	teamRepo *repositories.TeamRepository
}

func isAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("id"); user.(*entity.JWTTeam).TeamName == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func NewAdminRoutes(handler *gin.RouterGroup, l logger.Interface, teamRepo *repositories.TeamRepository, jwtMiddleware *jwt.GinJWTMiddleware, cfg *config.Config) {
	r := &routes{l, cfg, teamRepo}

	admin := handler.Group("/admin")

	admin.Use(jwtMiddleware.MiddlewareFunc())
	admin.Use(isAdmin())

	teamsGroup := admin.Group("/teams")

	teamsGroup.GET("/", r.TeamsList)
	teamsGroup.DELETE("/:name", r.DeleteTeam)

	admin.POST("/vpn", r.CreateVpnTeams)

	generateGroup := admin.Group("/generate")

	generateGroup.POST("/variables", r.GenerateVariables)
	generateGroup.POST("/ssh", r.GenerateSshKeysDir)
	generateGroup.POST("/prometheus", r.GeneratePrometheus)
}
