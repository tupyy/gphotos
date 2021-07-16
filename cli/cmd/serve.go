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
	"encoding/gob"
	"fmt"

	"github.com/gin-contrib/sessions/memstore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/controllers"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/internal/repo/postgres/album"
	groupRepo "github.com/tupyy/gophoto/internal/repo/postgres/group"
	userRepo "github.com/tupyy/gophoto/internal/repo/postgres/user"
	"github.com/tupyy/gophoto/utils/logutil"
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

		// initialize cookie store
		store := memstore.NewStore([]byte(conf.GetServerSecretKey()))

		// register sessionData
		gob.Register(entity.Session{})

		// initialize postgres client
		client, err := pgclient.NewClient(conf.GetPostgresConf())
		if err != nil {
			panic(err)
		}

		repos, err := createPostgresRepos(client)
		if err != nil {
			panic(err)
		}

		// initialize oidc provier
		oidcProvider := auth.NewOidcProvider(keycloakConf, conf.GetServerAuthCallback())

		keyCloakAuthenticator := auth.NewKeyCloakAuthenticator(oidcProvider, repos[repo.UserRepoName].(repo.UserRepo), repos[repo.GroupRepoName].(repo.GroupRepo))

		// create new router
		r := router.NewRouter(store, keyCloakAuthenticator)

		controllers.Logout(r.PrivateGroup, keyCloakAuthenticator)

		controllers.Register(r.PrivateGroup, r.PublicGroup, repos)

		// run server
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func createPostgresRepos(client pgclient.Client) (repo.Repositories, error) {
	repos := make(repo.Repositories)

	ur, err := userRepo.New(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return repos, err
	}

	repos[repo.UserRepoName] = ur

	gr, err := groupRepo.New(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return repos, err
	}

	repos[repo.GroupRepoName] = gr

	albumRepo, err := album.NewPostgresRepo(client)
	if err != nil {
		logutil.GetDefaultLogger().WithError(err).Warn("cannot create user repo")

		return repos, err
	}

	repos[repo.AlbumRepoName] = albumRepo

	return repos, nil
}
