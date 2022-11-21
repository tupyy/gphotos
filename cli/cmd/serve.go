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
	"fmt"

	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/auth"
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
	"github.com/tupyy/gophoto/internal/services/media"
	tagService "github.com/tupyy/gophoto/internal/services/tag"
	usersService "github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/internal/utils/logutil"
	"github.com/tupyy/gophoto/internal/utils/minioclient"
	"github.com/tupyy/gophoto/internal/utils/pgclient"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Conf used\n %s\n", conf.GetConfiguration())

		logrus.SetLevel(conf.GetLogLevel())
		logrus.SetReportCaller(true)

		// initialize cookie store
		store := memstore.NewStore([]byte(conf.GetServerSecretKey()))

		// register sessionData
		gob.Register(entity.Session{})

		// initialize postgres client
		client, err := pgclient.NewClient(conf.GetPostgresConf())
		if err != nil {
			panic(err)
		}
		logutil.GetDefaultLogger().Info("connected to db")

		// init minio client
		minioClient, err := minioclient.New(conf.GetMinioConfig())
		if err != nil {
			panic(err)
		}
		logutil.GetDefaultLogger().WithField("conf", conf.GetMinioConfig().String()).Info("connected at minio")

		// create keycloak authenticator
		keycloakAuthenticator := auth.NewKeyCloakAuthenticator(conf.GetKeycloakConfig(), conf.GetServerAuthCallback())

		// create new router
		engine := gin.New()
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
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return nil, err
	}

	// create album repo
	albumRepo, err := album.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create album repo")

		return nil, err
	}

	// create tag repo
	tagRepo, err := tag.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create tag repo")

		return nil, err
	}
	// create user repo
	userRepo, err := user.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create user pg repo")

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

	server := handlersv1.NewServer(albumService, usersService, tagService, mediaService)
	return server, nil
}
