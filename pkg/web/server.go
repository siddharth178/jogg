package web

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/toptal/sidd/jogg/pkg/config"
	"gitlab.com/toptal/sidd/jogg/pkg/datastore"
	"gitlab.com/toptal/sidd/jogg/pkg/weather"
	"go.uber.org/zap"
)

type Server struct {
	config config.HTTPConfig
	logger *zap.SugaredLogger

	ds *datastore.DS
	e  *echo.Echo
	ws weather.WeatherService
}

func NewServer(ctx context.Context,
	logger *zap.SugaredLogger,
	cfg config.HTTPConfig,
	ds *datastore.DS,
	ws weather.WeatherService) *Server {
	e := echo.New()
	e.HideBanner = true
	// e.Use(echozap.ZapLogger(logger.Desugar()))
	e.Use(middleware.Logger())

	return &Server{
		config: cfg,
		logger: logger,
		ds:     ds,
		e:      e,
		ws:     ws,
	}
}

func (s *Server) Start() error {
	s.setupRoutes()
	return s.e.Start(fmt.Sprintf(":%d", s.config.Port))
}

func (s *Server) setupRoutes() {
	s.logger.Info("setting up routes")

	// nop := NopController{s.logger}

	hc := HealthController{s.ds}

	lc := LoginController{s.ds, s.config.Secret}

	getUserByIDCtrl := GetUserByIDController{s.ds}
	addUserCtrl := AddUserController{s.ds}
	getAllUsers := GetAllUsersController{s.logger, s.ds}
	updateUserCtrl := UpdateUserController{s.logger, s.ds}
	deleteUserCtrl := DeleteUserController{s.logger, s.ds}

	addActivitiesCtrl := AddActivitiesController{s.logger, s.ds, s.ws}
	getActivityCtrl := GetActivityController{s.logger, s.ds}
	getAllActivitiesCtrl := GetAllActivitiesController{s.logger, s.ds}
	deleteActivityCtrl := DeleteActivityController{s.logger, s.ds}
	updateActivityCtrl := UpdateActivityController{s.logger, s.ds}

	weeklyReportCtrl := WeeklyReportController{s.logger, s.ds}

	// health ok
	s.e.GET("/health", hc.Health)

	// login
	s.e.POST("/login", lc.Login)
	s.e.POST("/register", addUserCtrl.AddUser) // self register

	//
	// APIs accessible only after login. i.e. valid JWT is expected
	//
	g := s.e.Group("/")
	g.Use(middleware.JWT([]byte(s.config.Secret)))

	// users
	g.POST("users", addUserCtrl.AddUser, AdminOrUserManager(s.logger))
	g.GET("users", getAllUsers.GetAllUsers, AdminOrUserManager(s.logger))
	g.GET("users/:user_id", getUserByIDCtrl.GetUserByID, ValidUser(s.logger))
	g.PUT("users/:user_id", updateUserCtrl.UpdateUser, AdminOrValidUser(s.logger))
	g.DELETE("users/:user_id", deleteUserCtrl.DeleteUser, AdminOrUserManager(s.logger))

	// activities
	g.POST("users/:user_id/activities", addActivitiesCtrl.AddActivities, AdminOrValidUser(s.logger))
	g.GET("users/:user_id/activities/:activity_id", getActivityCtrl.GetActivity, AdminOrValidUser(s.logger))
	g.GET("users/:user_id/activities", getAllActivitiesCtrl.GetAllActivities, AdminOrValidUser(s.logger))
	g.DELETE("users/:user_id/activities/:activity_id", deleteActivityCtrl.DeleteActivity, AdminOrValidUser(s.logger))
	g.PUT("users/:user_id/activities", updateActivityCtrl.UpdateActivity, AdminOrValidUser(s.logger))

	// report
	g.GET("reports/:user_id/weekly", weeklyReportCtrl.WeeklyReport, AdminOrValidUser(s.logger))
}
