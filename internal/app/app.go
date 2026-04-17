package app

import (
	"adoptme/config"
	"adoptme/internal/controller/restapi"
	"adoptme/internal/repo/persistent"
	"adoptme/internal/usecase/adoption"
	"adoptme/internal/usecase/catalog"
	"adoptme/internal/usecase/shelter"
	"adoptme/internal/usecase/volunteer"
	"adoptme/pkg/httpserver"
	"adoptme/pkg/jwt"
	"adoptme/pkg/logger"
	"adoptme/pkg/postgres"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type useCases struct {
	adoption  *adoption.UseCase
	catalog   *catalog.UseCase
	shelter   *shelter.UseCase
	volunteer *volunteer.UseCase
}

type servers struct {
	http *httpserver.Server
}

func initUseCases(pg *postgres.Postgres, jwtManager *jwt.Manager) useCases {
	animalRepo := persistent.NewAnimalRepo(pg)
	shelterRepo := persistent.NewShelterRepo(pg)
	volunteerRepo := persistent.NewVolunteerRepo(pg)

	return useCases{
		adoption:  adoption.New(animalRepo, shelterRepo, volunteerRepo),
		catalog:   catalog.New(animalRepo, shelterRepo, volunteerRepo),
		shelter:   shelter.New(shelterRepo, jwtManager),
		volunteer: volunteer.New(volunteerRepo, jwtManager),
	}
}

func initServers(cfg *config.Config, uc useCases, jwt *jwt.Manager, l logger.Interface) servers {
	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, uc.adoption, uc.catalog, uc.shelter, uc.volunteer, jwt, l)

	return servers{
		http: httpServer,
	}
}

func (s *servers) startServers() {
	s.http.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// JWT
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	// Use cases
	uc := initUseCases(pg, jwtManager)

	// Servers
	s := initServers(cfg, uc, jwtManager, l)
	s.startServers()
	s.waitForShutdown(l)
}
