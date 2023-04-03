package main

import (
	"awesome-clean-arch/internal/auth"
	"awesome-clean-arch/internal/auth/db/mysql"
	"awesome-clean-arch/internal/config"
	"awesome-clean-arch/internal/profile"
	"awesome-clean-arch/internal/profile/db/mysql"
	"awesome-clean-arch/internal/user"
	"awesome-clean-arch/internal/user/db/mysql"
	"awesome-clean-arch/internal/user_data"
	"awesome-clean-arch/internal/user_data/db/mysql"
	"awesome-clean-arch/pkg/client/mysql"
	"awesome-clean-arch/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()

	logger.Infoln("Create router...")
	router := httprouter.New()
	logger.Infoln("...created")

	cfg := config.GetConfig()

	mysqlClient, err := mysql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	logger.Infoln("Create authRepository...")
	authRepository := mysql_auth.NewMySQLRepository(mysqlClient, logger)
	logger.Infoln("...created")

	authHandler := auth.NewHandler(logger, authRepository)
	authHandler.Register(router)
	logger.Infoln("...created")

	logger.Infoln("Create userRepository...")
	userRepository := mysql_user.NewMySQLRepository(mysqlClient, logger)
	logger.Infoln("...created")

	logger.Infoln("Create userHandler...")
	userHandler := user.NewHandler(logger, userRepository)
	userHandler.Register(router)
	logger.Infoln("...created")

	logger.Infoln("Create profileRepository...")
	profileRepository := mysql_profile.NewMySQLRepository(mysqlClient, logger)
	logger.Infoln("...created")

	logger.Infoln("Create profileHandler...")
	profileHandler := profile.NewHandler(logger, profileRepository)
	profileHandler.Register(router)
	logger.Infoln("...created")

	logger.Infoln("Create userDataRepository...")
	userDataRepository := mysql_user_data.NewMySQLRepository(mysqlClient, logger)
	logger.Infoln("...created")

	logger.Infoln("Create userDataHandler...")
	userDataHandler := user_data.NewHandler(logger, userDataRepository)
	userDataHandler.Register(router)
	logger.Infoln("...created")

	logger.Infoln("Start router...")
	start(router, cfg)
	logger.Infoln("...started")
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Infoln("Start application")

	//listenAddr := flag.String("port", ":8080", "the server address")
	//flag.Parse()

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Infoln("Detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Infoln("Create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Infoln("Listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening on unix socket :%s", socketPath)
	} else {
		logger.Info("Listen tcp")
		//listener, listenErr = net.Listen("tcp", *listenAddr)
		//logger.Infoln("The server running on port", *listenAddr)
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("Server is listening on %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
