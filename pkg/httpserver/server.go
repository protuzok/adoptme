package httpserver

import (
	"adoptme/pkg/logger"
	"context"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App    *echo.Echo
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine for s.App.Start()

	s := &Server{
		ctx:             ctx,
		eg:              group,
		App:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	app := echo.New()
	app.Server.ReadTimeout = s.readTimeout
	app.Server.WriteTimeout = s.writeTimeout

	s.App = app

	return s
}

func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.App.Start(s.address)
		if err != nil {
			s.notify <- err

			close(s.notify)

			return err
		}

		return nil
	})

	s.logger.Info("restapi server - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var shutdownErrors []error

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.App.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "restapi server - Server - Shutdown - s.App.Shutdown")

		shutdownErrors = append(shutdownErrors, err)
	}

	// Wait for s.App.Start() goroutine to finish and get any error
	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		// todo future свапнути ті два аргументи, в шаблоні було погано описано
		s.logger.Error(err, "restapi server - Server - Shutdown - s.eg.Wait")

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("restapi server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
