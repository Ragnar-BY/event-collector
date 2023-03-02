package rest

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EventCollector is interface for collecting events
type EventUsecase interface {
	SaveEvents(ctx context.Context, events []domain.Event, clientIP string, serverTime time.Time) error

	Ping(ctx context.Context) error
}

// Server is http server
type Server struct {
	srv *http.Server
	log *zap.Logger

	events EventUsecase
	Now    func() time.Time // for testing we need possibility to override function
}

// NewServer creates new server instance
func NewServer(addr string, log *zap.Logger, events EventUsecase) *Server {
	e := gin.Default()

	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}
	s := Server{
		srv:    srv,
		log:    log,
		events: events,
		Now:    time.Now,
	}
	s.routes(e)
	return &s
}

// Routes adds routes to server
func (s *Server) routes(e *gin.Engine) {
	e.GET("/healtz", s.Healtz)
	e.POST("/events", s.SendEvents)
}

// Run starts server
func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop stops server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Healtz(c *gin.Context) {
	err := s.events.Ping(c.Request.Context())
	if err != nil {
		s.log.Error("can not ping event collector",
			zap.Error(err),
			zap.Time("server_time", time.Now()))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// SendEvents is handler for event sending
func (s *Server) SendEvents(c *gin.Context) {
	clientIP := c.ClientIP()
	serverTime := s.Now()

	scanner := bufio.NewScanner(c.Request.Body)
	defer c.Request.Body.Close()
	events := make([]domain.Event, 0)
	for scanner.Scan() {
		ev := Event{}
		err := json.Unmarshal(scanner.Bytes(), &ev)
		if err != nil {
			s.log.Error("can not parse events",
				zap.Error(err),
				zap.Time("server_time", serverTime),
				zap.String("client_ip", clientIP))
			c.JSON(http.StatusBadRequest, err)
			return
		}

		events = append(events, ev.ToDomainEvent())
	}

	err := s.events.SaveEvents(c.Request.Context(), events, clientIP, serverTime)
	if err != nil {
		s.log.Error("can not save events",
			zap.Error(err),
			zap.Time("server_time", serverTime),
			zap.String("client_ip", clientIP))
	}

	c.JSON(http.StatusOK, "ok")
}
