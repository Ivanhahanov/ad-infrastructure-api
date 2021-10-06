package v1

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/controller/http/v1/admin"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"

	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	_manager "github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/middleware"
	_walker "github.com/Ivanhahanov/ad-infrastructure-api/pkg/walker"

	jwt "github.com/appleboy/gin-jwt/v2"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func NewRouter(handler *gin.Engine, l logger.Interface, manager *_manager.CtfManager, walker *_walker.Walker, cfg *config.Config) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(middleware.CORS())

	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entity.JWTTeam); ok {
				return jwt.MapClaims{
					"id": v.TeamName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &entity.JWTTeam{
				TeamName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			userID := loginVals.Username
			password := loginVals.Password

			team, dbErr := manager.TeamRepo.FilterTeams(map[string]interface{}{"name": userID})
			if dbErr != nil {
				log.Println(dbErr)
				return nil, jwt.ErrFailedAuthentication
			}

			log.Println(team)

			if CheckPasswordHash(password, team[0].Hash) {
				return &entity.JWTTeam{TeamName: userID}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(*entity.JWTTeam)

			return ok
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	})
	if err != nil {log.Fatal("JWT Error:" + err.Error())}

	errInit := jwtMiddleware.MiddlewareInit()
	if errInit != nil {log.Fatal("jwtMiddleware.MiddlewareInit() Error:" + errInit.Error())}

	v1Group := handler.Group("/api/v1")

	// Check health
	v1Group.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	v1Group.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	admin.NewAdminRoutes(v1Group, l, manager.TeamRepo, jwtMiddleware, cfg)
	NewAuthRoutes(v1Group, l, manager.TeamRepo, jwtMiddleware, cfg)
	NewCheckerRoutes(v1Group, l, manager, walker, cfg.CheckerPassword)
	NewFlagRoutes(v1Group, l, manager, jwtMiddleware)
	NewScoreboardRoutes(v1Group, l, manager, jwtMiddleware)
}
