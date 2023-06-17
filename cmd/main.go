package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"robinhood/cmd/httpserver"
	"robinhood/config"
	infrastructure "robinhood/infrastructures"
	"robinhood/internal/core/services/blogsvc"
	"robinhood/internal/core/services/commentsvc"
	"robinhood/internal/core/services/usersvc"
	"robinhood/internal/handlers/bloghdl"
	"robinhood/internal/handlers/userhdl"
	"robinhood/internal/repositories"
	"syscall"
	"time"
)

// @title           Robinhood test API
// @version         1.0
// @description     This is a Robinhood test API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Tanatorn Nateesanpraser
// @contact.email  tanatorn.nateesanprasert@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func init() {
	config.New()
}

func main() {
	// infrastructures
	mc := infrastructure.NewMongoDB()

	// repositories
	br := repositories.NewBlogRepository(mc, config.Get().Mongo.Database)
	cr := repositories.NewCommentRepository(mc, config.Get().Mongo.Database)
	ur := repositories.NewUserRepository(mc, config.Get().Mongo.Database)
	// services
	bs := blogsvc.New(br, ur)
	cs := commentsvc.New(cr, ur)
	us := usersvc.New(ur)
	// handlers
	bh := bloghdl.New(bs, cs)
	uh := userhdl.New(us)

	e := httpserver.NewHTTPServer(bh, uh)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", config.Get().Endpoint.Port)); err != nil {
			e.Logger.Info("shutting down the server")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer e.Shutdown(ctx)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
