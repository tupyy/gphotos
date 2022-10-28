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
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	internalApiv1 "github.com/tupyy/gophoto/internal/api/v1"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
	keycloakRepo "github.com/tupyy/gophoto/internal/domain/keycloak"
	miniorepo "github.com/tupyy/gophoto/internal/domain/minio"
	"github.com/tupyy/gophoto/internal/domain/postgres/album"
	"github.com/tupyy/gophoto/internal/domain/postgres/tag"
	"github.com/tupyy/gophoto/internal/domain/postgres/user"
	"github.com/tupyy/gophoto/internal/entity"
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
		keycloakConf := conf.GetKeycloakConfig()

		fmt.Printf("Conf used\n %s\n", keycloakConf.String())

		logrus.SetLevel(conf.GetLogLevel())
		logrus.SetReportCaller(true)
		logrus.SetFormatter(conf.GetLogFormatter())

		// initialize cookie store
		store := memstore.NewStore([]byte(conf.GetServerSecretKey()))

		// register sessionData
		gob.Register(entity.Session{})
		gob.Register(entity.Alert{})

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

		services, err := initServices(client, minioClient)
		if err != nil {
			panic(err)
		}
		logutil.GetDefaultLogger().Info("services created")

		// create keycloak
		keycloakAuthenticator := auth.NewKeyCloakAuthenticator(conf.GetKeycloakConfig(), conf.GetServerAuthCallback())

		// create new router
		r := router.NewRouter(store, keycloakAuthenticator)

		//api.Logout(r.PrivateGroup, keycloakAuthenticator)

		server := internalApiv1.NewServer(services)
		apiv1.RegisterHandlers(r.PrivateGroup, server)

		// run server
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func initServices(client pgclient.Client, mclient *minio.Client) (map[string]interface{}, error) {
	services := make(map[string]interface{})

	// create keycloak repo
	kr, err := keycloakRepo.New(context.Background(), conf.GetKeycloakConfig())
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return services, err
	}

	// create album repo
	albumRepo, err := album.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create album repo")

		return services, err
	}

	// create tag repo
	tagRepo, err := tag.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create tag repo")

		return services, err
	}
	// create user repo
	userRepo, err := user.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create user pg repo")

		return services, err
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

	return services, nil
}
