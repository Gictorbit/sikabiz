package api

import (
	"github.com/gictorbit/sikabiz/db/userdb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net"
)

type UserService struct {
	db      userdb.UserDBPgConn
	logger  *zap.Logger
	server  *echo.Echo
	address string
}

func NewUserService(logger *zap.Logger, db userdb.UserDBPgConn, host, port string) *UserService {
	return &UserService{
		db:      db,
		logger:  logger,
		address: net.JoinHostPort(host, port),
	}
}

func (us *UserService) InitHandlers() {

	// Register the route with the GET handler
	us.server.GET("/users/:id", us.GetUser)
}

func (us *UserService) Run() error {
	// Start the server
	us.InitHandlers()
	return us.server.Start(us.address)
}
