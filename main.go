package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/common"
	"github.com/speshiy/LinuxMonitorGo/cron"
	"github.com/speshiy/LinuxMonitorGo/routes"
	"github.com/speshiy/LinuxMonitorGo/service/controllers"
	serviceRoutes "github.com/speshiy/LinuxMonitorGo/service/routes"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

func main() {
	//Initialize global variables
	common.InitGlobalVars()

	if settings.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	//Initilize default routes
	router := gin.Default()
	routerService := gin.Default()

	//Setting CORS params for request Headers

	//Add cors to middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"X-CSRF-Token"},
		AllowCredentials: false,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	}))

	//Initializing app routes
	router = routes.InitRoutes(router)
	routerService = serviceRoutes.InitRoutes(routerService)

	//Auto create database and tables
	controllers.MigrateDo(nil)

	//Start CRON BACKUP MYSQL
	cron.InitCron()

	//Starting API and Service Servers
	StartServer(router, routerService)
}

//StartServer with graceful stop
func StartServer(router *gin.Engine, routerService *gin.Engine) {
	srv := &http.Server{
		Addr:         ":" + settings.Port,
		Handler:      router,
		ReadTimeout:  10800 * time.Second,
		WriteTimeout: 10800 * time.Second,
	}

	srvService := &http.Server{
		Addr:         ":" + settings.PortService,
		Handler:      routerService,
		ReadTimeout:  10800 * time.Second,
		WriteTimeout: 10800 * time.Second,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		log.Println("HTTP server Shutdown: ", srv.Addr, "successfull")

		// We received an interrupt signal, shut down.
		if err := srvService.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		log.Println("HTTP server Shutdown: ", srvService.Addr, "successfull")

		close(idleConnsClosed)
	}()

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP/S server ListenAndServe: %v", err)
		}
	}()

	go func() {
		err := srvService.ListenAndServe()
		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP/S server ListenAndServe: %v", err)
		}
	}()

	<-idleConnsClosed
}
