package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

var logger zerolog.Logger

func main() {
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	viper.AutomaticEnv()
	setEnv()
	address := viper.GetString("address")
	//useTLS := viper.GetBool("rbac_use_tls")
	//apiuseTLS := viper.GetBool("control_api_use_tls")

	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	api.GET("/rbacState", stateHandler)
	//api.POST("/jokes/like/:jokeID", LikeJoke)

	// Start and run the server
	router.Run(address)

	//server := http.Server{
	//	Addr:    address,
	//	Handler: New(&logger),
	//}
	//
	//server.ListenAndServe()
}

func stateHandler(c *gin.Context) {
	svcState := getState()
	c.JSON(http.StatusOK, svcState)
}

func setEnv() {
	viper.SetDefault("address", "0.0.0.0:8000")
	viper.SetDefault("gm_control_api_address", "0.0.0.0:5555")
}

// New generates a new handler
func New(logger *zerolog.Logger) http.Handler {

	mux := mux.NewRouter()

	// handlers attached directly to the mux are considered outside the API
	// they do not use the RPC functions or the JSON messages

	mux.HandleFunc("/", helloWorld)
	mux.HandleFunc("/getPolicies/{proxy}", getPolicy)
	//mux.HandleFunc("/getState", getState)
	//mux.HandleFunc("/getPolicy", getPolicy)
	mux.HandleFunc("/addPolicy", addPolicy)
	return mux

	//mux.HandleFunc("/logging", st.handleLogLevelGET).Methods("GET")
	//mux.HandleFunc("/logging", st.handleLogLevelPUT).Methods("PUT").Queries("level", "{level}")
	//mux.HandleFunc("/logging", nonAPIMethodNotAllowedFactory(st.logger, "GET", "PUT"))

	//mux.HandleFunc("/add", addPolicy)
	//mux.HandleFunc("/delete", deletePolicy)
	//mux.HandleFunc("/modify", modifyPolicy)
	//mux.HandleFunc("/get", getPolicies)

	// Cluster endpoints
	//subrouter.HandleFunc("/cluster", st.handleV1ClustersGET).Methods("GET")
	//subrouter.HandleFunc("/cluster", st.handleV1ClustersPOST).Methods("POST")
	//subrouter.HandleFunc("/cluster", methodNotAllowedFactory(st.logger, "GET", "POST"))
}
