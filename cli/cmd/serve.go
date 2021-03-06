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
	"github.com/tupyy/gophoto/internal/api"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	keycloakRepo "github.com/tupyy/gophoto/internal/domain/keycloak"
	miniorepo "github.com/tupyy/gophoto/internal/domain/minio"
	"github.com/tupyy/gophoto/internal/domain/postgres/album"
	"github.com/tupyy/gophoto/internal/domain/postgres/tag"
	"github.com/tupyy/gophoto/internal/domain/postgres/user"
	"github.com/tupyy/gophoto/internal/entity"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/media"
	tagService "github.com/tupyy/gophoto/internal/services/tag"
	usersService "github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
	"github.com/tupyy/gophoto/utils/minioclient"
	"github.com/tupyy/gophoto/utils/pgclient"

	router "github.com/tupyy/gophoto/internal/routes"
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

		repos, err := createRepos(client, minioClient)
		if err != nil {
			panic(err)
		}
		logutil.GetDefaultLogger().Info("repositories created")

		// create services
		mediaService := media.New(repos[domain.MinioRepoName].(domain.Store))

		albumRepo := repos[domain.AlbumRepoName].(domain.Album)
		tagRepo := repos[domain.TagRepoName].(domain.Tag)

		albumService := albumService.New(albumRepo, mediaService)
		usersService := usersService.New(repos)
		tagService := tagService.New(tagRepo)

		logutil.GetDefaultLogger().Info("services created")

		// create keycloak
		keycloakAuthenticator := auth.NewKeyCloakAuthenticator(conf.GetKeycloakConfig(), conf.GetServerAuthCallback())

		// create new router
		r := router.NewRouter(store, keycloakAuthenticator)

		api.Logout(r.PrivateGroup, keycloakAuthenticator)

		api.RegisterIndexHandler(r.PrivateGroup, usersService)
		api.RegisterAlbumHandler(r.PrivateGroup, albumService, usersService, tagService)
		api.RegisterMediaHandler(r.PrivateGroup, albumService, mediaService)
		api.RegisterTagHandler(r.PrivateGroup, albumService, tagService)

		// run server
		r.Run()

		// TODO shutdown
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func createRepos(client pgclient.Client, mclient *minio.Client) (domain.Repositories, error) {
	repos := make(domain.Repositories)

	// create keycloak repo
	kr, err := keycloakRepo.New(context.Background(), conf.GetKeycloakConfig())
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return repos, err
	}

	ttl, interval := conf.GetRepoCacheConfig()

	repos[domain.KeycloakRepoName] = keycloakRepo.NewCacheRepo(kr, ttl, interval)

	// create album repo
	albumRepo, err := album.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create album repo")

		return repos, err
	}
	repos[domain.AlbumRepoName] = albumRepo

	// create tag repo
	tagRepo, err := tag.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create tag repo")

		return repos, err
	}
	repos[domain.TagRepoName] = tagRepo

	// create user repo
	userRepo, err := user.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("failed to create user pg repo")

		return repos, err
	}

	repos[domain.UserRepoName] = userRepo

	// create minio repo
	minioRepo := miniorepo.New(mclient)
	minioCache := miniorepo.NewCacheRepo(minioRepo, ttl, interval)
	repos[domain.MinioRepoName] = minioCache

	return repos, nil
}
