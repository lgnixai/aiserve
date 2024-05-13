package main

import (
	"log"
	`log/slog`
	"os"

	`aurora/config`
	database `aurora/pkg/db`
	`aurora/pkg/platform`
	`aurora/pkg/server`

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

const (
	wsString  = "ws://localhost:9876/rpc"
	tableName = "user"
	namespace = "test"
	dbName    = "test"
	port      = 9876
	username  = "root"
	password  = "root"
)

func main() {
	config.InitLlm()
	gin.SetMode(gin.ReleaseMode)
	db, err := database.NewDatabase(wsString, username, password, namespace, dbName)
	if err != nil {
		slog.Error("failed to create shortener repository", "error", err)

	}

	slog.Info("Connected to database")
	// defer database.DB.Close()
	// Close connections to the database at program shutdown
	defer func() {
		slog.Info("Closing database")
		db.Close()
	}()
	//router := initialize.RegisterRouter(db)
	//subFS, err := fs.Sub(staticFiles, "web")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//router.StaticFS("/web", http.FS(subFS))

	_ = godotenv.Load(".env")
	host := os.Getenv("SERVER_HOST")
	addr := os.Getenv("SERVER_PORT")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	if addr == "" {
		addr = os.Getenv("PORT")
		if addr == "" {
			addr = "9877"
		}
	}

	srv := server.NewServer(db, addr)

	//if basepath != "" {
	//	srv.BasePath = "/" + strings.Trim(basepath, "/")
	//}
	if host != "" {
		srv.Host = host
	} else {
		srv.Host = "0.0.0.0"
	}

	if tlsCert != "" && tlsKey != "" {
		srv.CertFile = tlsCert
		srv.KeyFile = tlsKey
	}

	if username != "" && password != "" {
		srv.Username = username
		srv.Password = password
	}

	log.Printf("starting server at %s", srv.GetAddr())
	//if open {
	//	platform.Open(srv.GetAddr())
	//}
	platform.Start(srv)

}
