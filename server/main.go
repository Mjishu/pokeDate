package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/lxzan/gws"
	"github.com/mjishu/pokeDate/controller"
	"github.com/mjishu/pokeDate/database"
)

type apiConfig struct {
	jwt_secret       string
	assetPath        string
	database_url     string
	s3Bucket         string
	s3Region         string
	s3CfDistribution string
	s3Client         *s3.Client
}

// set up s3

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("could not load .env")
	}

	// * load params
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("jwt secret is emtpy!")

	}
	assetPath := os.Getenv("ASSET_PATH")
	if assetPath == "" {
		log.Fatal("asset path is empty")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Could not find database path")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("s3Bucket variable not set")
	}
	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("s3Region variable not set")
	}
	s3CfDist := os.Getenv("S3_CF_DISTRO")
	if s3CfDist == "" {
		log.Fatal("s3 cf variable dist not set")
	}

	awscfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s3Region))
	if err != nil {
		log.Fatal("unable to detect aws profile")
	}
	s3Client := s3.NewFromConfig(awscfg)

	config := apiConfig{
		jwt_secret:       jwt_secret,
		assetPath:        assetPath,
		database_url:     databaseURL,
		s3Bucket:         s3Bucket,
		s3Region:         s3Region,
		s3CfDistribution: s3CfDist,
		s3Client:         s3Client,
	}

	_, pool := database.CreateConnection()

	handler := &Handler{pool: pool}

	// web socket
	upgrader := gws.NewUpgrader(handler, &gws.ServerOption{
		ParallelEnabled:   true,
		Recovery:          gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})

	// user info
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		controller.UserController(w, r, pool, config.jwt_secret, config.s3Bucket, config.s3Region)
	})
	mux.HandleFunc("POST /users/profile_pictures/{userID}", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleUserImageUpload(w, r, pool, config.jwt_secret, config.s3Bucket, config.s3Region, config.s3Client)
	})
	mux.HandleFunc("POST /users/progress/reset", func(w http.ResponseWriter, r *http.Request) {
		controller.ResetSeenProgress(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("/organizations/", func(w http.ResponseWriter, r *http.Request) {
		controller.OrganizationController(w, r, pool, config.jwt_secret, s3Bucket, s3Region)
	})

	mux.HandleFunc("POST /animals/{animalID}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetAnimal(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("POST /animals/images/{animalID}", func(w http.ResponseWriter, r *http.Request) {
		controller.UploadAnimalImage(w, r, pool, config.jwt_secret, config.s3Bucket, config.s3Region, config.s3Client)
	})

	mux.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		controller.CardsController(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("/animals", func(w http.ResponseWriter, r *http.Request) {
		controller.AnimalController(w, r, pool, jwt_secret, config.s3Bucket, config.s3Client)
	})
	mux.HandleFunc("/shots/animals", func(w http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc("/shots", func(w http.ResponseWriter, r *http.Request) {
		controller.ShotController(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		controller.RefreshToken(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("/revoke", func(w http.ResponseWriter, r *http.Request) {
		controller.RevokeToken(w, r, pool)
	})

	//* Messages
	mux.HandleFunc("POST /messages", func(w http.ResponseWriter, r *http.Request) {
		controller.CurrentUserMessages(w, r, pool, config.jwt_secret)
	})
	mux.HandleFunc("POST /messages/create", func(w http.ResponseWriter, r *http.Request) {
		controller.CreateConversation(w, r, pool, config.jwt_secret)
	})

	mux.HandleFunc("POST /messages/{messageID}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetMessage(w, r, pool, config.jwt_secret)
	})

	mux.HandleFunc("POST /messages/{messageID}/send", func(w http.ResponseWriter, r *http.Request) {
		controller.CreateMessage(w, r, pool, config.jwt_secret)
	})

	//* Notifications
	mux.HandleFunc("POST /notifications/", func(w http.ResponseWriter, r *http.Request) {
		controller.CreateNotification(w, r, pool, config.jwt_secret, "New message request", 1)
	})

	mux.HandleFunc("GET /notifications", func(w http.ResponseWriter, r *http.Request) {
		controller.GetNotifications(w, r, pool, config.jwt_secret)
	})

	//* WebSocket
	mux.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r)
		if err != nil {
			return
		}
		go func() {
			socket.ReadLoop()
		}()
		fmt.Println("connection was made")
	})

	database.Database(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("listening on port " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
