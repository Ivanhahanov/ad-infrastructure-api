package app

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"

	v1 "github.com/Ivanhahanov/ad-infrastructure-api/internal/controller/http/v1"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"
	_db "github.com/Ivanhahanov/ad-infrastructure-api/pkg/db"

	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/httpserver"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager/repositories"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/utils"
	_walker "github.com/Ivanhahanov/ad-infrastructure-api/pkg/walker"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func SeedAdmin(repo *repositories.TeamRepository) {
	teams, _ := repo.GetUsers()
	for _, team := range teams {
		if team.Name == "admin" {return}
	}

	password := utils.GetEnv("ADMIN_PASS", "admin")

	hash, _ := v1.HashPassword(password)

	repo.CreateTeam(&entity.Team{
		ID:        primitive.NewObjectID(),
		Name:      "admin",
		Hash:      hash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func Run(cfg *config.Config) {

	rand.Seed(time.Now().UnixNano())

	l := logger.New(cfg.Log.Level)

	db := _db.New()
	repoManager := repositories.NewRepoManager(db)
	ctfManager  := manager.NewManager(repoManager, cfg)

	walker, _  := _walker.New(ctfManager)

	SeedAdmin(ctfManager.TeamRepo)

	handler := gin.New()
	v1.NewRouter(handler, l, ctfManager, walker, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		err = httpServer.Shutdown()
		if err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}