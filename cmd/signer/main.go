package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/sshilko/go-signer/handler/web"
	"github.com/sshilko/go-signer/internal/server"
	"github.com/sshilko/go-signer/internal/signaling"
	"github.com/sshilko/go-signer/pkg/db"
	"github.com/sshilko/go-signer/pkg/jwt"
	"github.com/sshilko/go-signer/service/homework"
	"net/http"
	"sync"
)

type Config struct {
	DatabaseDir      string `envconfig:"DB_DIR" default:"/tmp/signerdb"`
	DatabaseName     string `envconfig:"DB_NAME" default:"nutsdb1"`
	ServerListenPort uint16 `envconfig:"LISTEN_PORT" default:"8000"`
	JWTSecret        []byte `envconfig:"JWT_SECRET" default:"secret123"`
}

//var exampleUserID = uuid.New()

func main() {
	logger := log.New("example")
	logger.SetHeader("${level} ")

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx, ctxClose := context.WithCancel(context.Background())

	jwtService := jwt.NewSimpleJWT(config.JWTSecret)

	//exampleToken, _ := jwtService.ExampleToken(exampleUserID)
	//helper(exampleUserID, exampleToken)

	srv := server.NewServer(config.ServerListenPort, jwtService)

	storage, err := db.NewDB(config.DatabaseDir, config.DatabaseName)
	if err != nil {
		logger.Fatal(err.Error())
	}
	homeworkRepo := homework.NewRepository(storage)
	homeworkService := homework.NewService(homeworkRepo)

	homeworkHandler := web.NewHomeworkHandler(homeworkService)
	srv.AddHandler(homeworkHandler)

	var wg sync.WaitGroup

	wg.Add(1)
	go signaling.HandleSignals(ctx, func(c context.Context) {
		defer wg.Done() // indicate Done for WG no matter what
		logger.Info("Stop command received, finishing all jobs")
		if err := srv.StopServer(c); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Stopped web server job")

		ctxClose()
		logger.Info("Context closed")

		if err := storage.Disconnect(); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("Closed db")
	})

	logger.Info("Starting web server job")
	if err = srv.StartServer(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		logger.Fatal(err.Error())
	}

	wg.Wait()
	logger.Info("All jobs stopped")
}

//
//func helper(exampleUserID uuid.UUID, jwtToken string) {
//	fmt.Printf("\nExample USER ID is %s\n", exampleUserID.String())
//	fmt.Println("JWT authentication for example USER ID is")
//	fmt.Printf("Authorization: Bearer %s\n\n", jwtToken)
//
//	fmt.Println("You can now make requests to:")
//	fmt.Println("POST /users/:userID/homework")
//	fmt.Println("GET /users/:userID/homework/:homeworkID")
//
//	fmt.Printf("\nUSER1=%s\n", exampleUserID.String())
//	fmt.Printf("TOKEN=%s\n\n", jwtToken)
//}
