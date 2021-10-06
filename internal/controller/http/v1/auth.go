package v1

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"

	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"

	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager/repositories"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"time"
)

type authRoutes struct {
	l         logger.Interface
	cfg      *config.Config
	teamRepo *repositories.TeamRepository
}
func NewAuthRoutes(handler *gin.RouterGroup, l logger.Interface, teamRepo *repositories.TeamRepository, jwtMiddleware *jwt.GinJWTMiddleware, cfg *config.Config) {
	r := &authRoutes{l, cfg, teamRepo}

	handler.POST("/login", jwtMiddleware.LoginHandler)
	handler.POST("/register", r.CreateTeam)
	handler.GET("/refresh-token", jwtMiddleware.RefreshHandler)
}

type Team struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	SshPubKey string `json:"ssh_pub_key"`
}

func (r *authRoutes) CreateTeam(c *gin.Context) {
	var team Team

	err := c.BindJSON(&team)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teams, err := r.teamRepo.GetTeams()
	if err != nil {log.Println(err)}

	hash, err := HashPassword(team.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: check team name

	err = r.teamRepo.CreateTeam(&entity.Team{
		ID:        primitive.NewObjectID(),
		Name:      team.Name,
		Address:   r.generateIp(len(teams)),
		Hash:      hash,
		SshPubKey: team.SshPubKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("The team %s created", team.Name),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *authRoutes) generateIp(number int) string {
	ip, _, err := net.ParseCIDR(r.cfg.Network)
	if err != nil {log.Println(err.Error())}

	ip = ip.To4()
	ip[3] += byte(number + 11)
	return ip.String()
}
