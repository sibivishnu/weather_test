package main

import (
	"log"
	"net/http"
	"os"
	firebase "src/grpc/go/pkg/mod/firebase.google.com/go@v3.13.0+incompatible"
	"src/grpc/go/pkg/mod/github.com/urfave/cli@v1.22.12"
	"src/grpc/go/pkg/mod/google.golang.org/api@v0.122.0/option"
	"src/grpc/go/weather-service/common"

	"github.com/gorilla/mux"
)

// func main() {
// 	http.HandleFunc("/api/hello", helloHandler)

// 	fmt.Println("Server listening on port 8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Println("ListenAndServe: ", err)
// 	}
// }

func setupHTTP(port string) {
	log.Println("[WebApp] Starting the http server on port : " + port)
	router := mux.NewRouter()

	// Root
	router.HandleFunc("/", actionDisplayCheckPage).Methods("GET")

	http.ListenAndServe(":"+port, router)
}

func main() {
	app := cli.NewApp()
	app.Name = "Weather Service"
	app.Usage = "Weather Cache Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  FLAG_HTTP_PORT,
			Value: "5000",
			Usage: "Http Port",
		},
	}

	app.Action = runIt
	app.Run(os.Args)
}

// ----------------------------------------------
// runIt - Application Action
// ----------------------------------------------
func runIt(runtimeContext *cli.Context) {
	log.Println("[WebApp] begin")

	//-------------------------------------------------
	// Init Globals
	//-------------------------------------------------
	options = make(map[string]interface{})
	options["redis.host"] = os.Getenv(ENV_REDIS_HOST)
	options["accuweather.key"] = os.Getenv(ENV_ACCU_API_KEY)
	options["datastore.project"] = "lax-gateway" // os.Getenv(ENV_PROJECT_ID)
	options["config.categories"] = "/conf/categories.json"
	init.LoadCommonEnvironment(options)

	// Locals
	firebaseServiceFile := os.Getenv(ENV_FIREBASE_SERVICE_FILE)
	opt := option.WithCredentialsFile(firebaseServiceFile)

	// Globals
	httpHost = os.Getenv(FLAG_HTTP_HOST)
	httpScheme = os.Getenv(FLAG_HTTP_SCHEME)

	// Configure FireBase App
	app, err := firebase.NewApp(common.CTX, nil, opt)
	if err != nil {
		log.Println("[WebApp] Firebase NewApp Error")
		panic(err)
	}

	// Configure Firebase Client
	firebaseClient, err = app.Auth(common.CTX)
	if err != nil {
		log.Println("[WebApp] Firebase Auth Error")
		panic(err)
	}

	// Prepare Http Request Handlers
	setupHTTP(runtimeContext.String(FLAG_HTTP_PORT))
}
