package cmd

import (
	"context"
	"fmt"
	"os"

	data "backend-service/data"
	defn "backend-service/defn"
	rest "backend-service/rest"
	user "backend-service/rest/user"
	"backend-service/service"
	util "backend-service/util"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

// var cfgFile, dc, cluster, serviceInstanceID string
// var config_url, config_apikey string
// var servicediscovery_url, servicediscovery_apikey string
// var rootcacert, svccert, svckey, jwt_public_key, jwt_private_key string
// var secureservestr string
var (
	output string
	port   string = "7000"
)

var RootCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Backend Service for portfolio website",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//instantiate with psql server
		ctx := context.Background()
		ctx = context.WithValue(ctx, "output-format", output)

		util.InitiateGlobalLogger(ctx)
		log := util.GetGlobalLogger(ctx)
		log.Println("Connecting to database...")
		database := data.ConnectDatabase(ctx)
		log.Println("Database connected successfully")
		if err := database.InitialiseDatabaseTables(ctx); err != nil {
			fmt.Printf("failed to Initialise database tables: %s\n", err.Error())
			os.Exit(1)
		}

		// stopCh := make(chan os.Signal, 1)
		// signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

		//Finally start the server
		StartServer(ctx, database)

		// <-stopCh
		// fmt.Println("Stop signal received")

		//Shutting down the server
		ShutdownServer(ctx)
	},
}

func StartServer(ctx context.Context, database data.Database) {
	router := httprouter.New()

	readTimeout := defn.ReadTimeout
	writeTimeout := defn.WriteTimeout
	readHeaderTimeout := defn.ReadHeaderTimeout

	router.PanicHandler = util.PanicHandler

	userHandler := SetupRouters(database)
	rest.AddRoutes(router, userHandler)
	handler := rest.ApplyMiddleware(router)
	handler = cors.Default().Handler(handler)
	// handler := router

	host := "localhost"

	var err error
	_, _, err = util.StartHTTPServer(fmt.Sprintf("%s:%s", host, port), handler, readTimeout, writeTimeout, readHeaderTimeout)
	if err != nil {
		fmt.Printf("failed to StartHTTPServer: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("server started")
}

func ShutdownServer(ctx context.Context) {
	util.ShutdownHTTPServer()
	fmt.Println("shutdown complete.")
}

func SetupRouters(database data.Database) *user.UserRoutesHandler {
	return user.NewUserRoutesHandler(service.NewUserService(data.NewUserRepository(database)))
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "location for log output (default: stdout)")
	RootCmd.PersistentFlags().StringVarP(&port, "port", "p", "7000", "port number for the server")
}
