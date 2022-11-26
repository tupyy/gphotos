/*
Copyright © 2021 Cosmin Tupangiu <cosmin.tupangiu@gmail.com>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/gin-contrib/sessions/memstore"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/auth"
	minioclient "github.com/tupyy/gophoto/internal/clients/minio"
	pgclient "github.com/tupyy/gophoto/internal/clients/pg"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	handlersv1 "github.com/tupyy/gophoto/internal/handlers/v1"
	keycloakRepo "github.com/tupyy/gophoto/internal/repos/keycloak"
	miniorepo "github.com/tupyy/gophoto/internal/repos/minio"
	"github.com/tupyy/gophoto/internal/repos/postgres/album"
	"github.com/tupyy/gophoto/internal/repos/postgres/tag"
	"github.com/tupyy/gophoto/internal/repos/postgres/user"
	"github.com/tupyy/gophoto/internal/router"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/encryption"
	"github.com/tupyy/gophoto/internal/services/media"
	tagService "github.com/tupyy/gophoto/internal/services/tag"
	usersService "github.com/tupyy/gophoto/internal/services/users"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := setupLogger()
		defer logger.Sync()

		undo := zap.ReplaceGlobals(logger)
		defer undo()

		zap.S().Infof("Configuration %s", conf.GetConfiguration())

		// initialize cookie store
		store := memstore.NewStore([]byte(conf.GetServerSecretKey()))

		// register sessionData
		gob.Register(entity.Session{})

		// initialize postgres client
		client, err := pgclient.New(conf.GetPostgresConf())
		if err != nil {
			panic(err)
		}
		zap.S().Infow("connected to db", "conf", conf.GetPostgresConf())

		// init minio client
		minioClient, err := minioclient.New(conf.GetMinioConfig())
		if err != nil {
			panic(err)
		}
		zap.S().Infow("connected to minio", "conf", conf.GetMinioConfig())

		// create keycloak authenticator
		keycloakAuthenticator := auth.NewKeyCloakAuthenticator(conf.GetKeycloakConfig(), conf.GetServerAuthCallback())

		// create new router
		engine := gin.New()
		engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		engine.Use(ginzap.RecoveryWithZap(logger, true))
		router.InitEngine(engine, store, keycloakAuthenticator)

		//api.Logout(r.PrivateGroup, keycloakAuthenticator)

		opt := apiv1.GinServerOptions{
			Middlewares: make([]apiv1.MiddlewareFunc, 0),
		}

		server, err := createServer(client, minioClient)
		if err != nil {
			panic(err)
		}

		apiv1.RegisterHandlersWithOptions(engine, server, opt)

		// run server
		engine.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func createServer(client pgclient.Client, mclient *minio.Client) (*handlersv1.Server, error) {
	services := make(map[string]interface{})

	// create keycloak repo
	kr, err := keycloakRepo.New(context.Background(), conf.GetKeycloakConfig())
	if err != nil {
		return nil, err
	}

	// create album repo
	albumRepo, err := album.NewPostgresRepo(client)
	if err != nil {
		return nil, err
	}

	// create tag repo
	tagRepo, err := tag.NewPostgresRepo(client)
	if err != nil {
		return nil, err
	}
	// create user repo
	userRepo, err := user.NewPostgresRepo(client)
	if err != nil {
		return nil, err
	}

	// create minio repo
	minioRepo := miniorepo.New(mclient)
	mediaService := media.New(minioRepo)

	albumService := albumService.New(albumRepo, mediaService)
	usersService := usersService.New(kr, userRepo)
	tagService := tagService.New(tagRepo)

	services["album"] = albumService
	services["user"] = usersService
	services["tag"] = tagService

	encryption, err := encryption.New()
	if err != nil {
		return nil, err
	}

	server := handlersv1.NewServer(albumService, usersService, tagService, mediaService, encryption)
	return server, nil
}

func setupLogger() *zap.Logger {
	loggerCfg := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder, EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	loggerCfg.Level = zap.NewAtomicLevelAt(conf.GetLogLevel())

	logger, err := loggerCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}

	return logger
}
