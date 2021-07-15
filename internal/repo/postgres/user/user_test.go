package user_test

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/tupyy/gophoto/internal/entity"
	userrepo "github.com/tupyy/gophoto/internal/repo/postgres/user"
	pgclient "github.com/tupyy/gophoto/utils/pgclient"
	"github.com/tupyy/gophoto/utils/pgtestcontainer"
)

const (
	// name of the root folder for the project.
	parentFolder = "gophoto"
	// sql setup file relative to parent folder.
	setupFile = "sql/setup/02_setup.sql"
	// fixtures file relative to parent folder.
	fixtureFile = "sql/fixtures/user_test.sql"
)

type UserTestSuite struct {
	suite.Suite
	container testcontainers.Container
	pgClient  pgclient.Client
}

func (u *UserTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parentFolder, err := getParentFolder()
	if err != nil {
		panic(err)
	}

	initMap := make(map[string]string)
	setupKey := path.Join(parentFolder, setupFile)
	initMap[setupKey] = "/docker-entrypoint-initdb.d/setup.sql"

	fixtureKey := path.Join(parentFolder, fixtureFile)
	initMap[fixtureKey] = "/docker-entrypoint-initdb.d/zz_fixtures.sql"

	c, err := pgtestcontainer.NewPostgreSQLContainer(ctx, pgtestcontainer.PostgreSQLContainerRequest{
		BindMounts: initMap,
	})
	if err != nil {
		panic(err)
	}

	u.pgClient, err = c.GetInitialClient(ctx)
	if err != nil {
		panic(err)
	}

	u.container = c.Container
}

func (u *UserTestSuite) TearDownSuite() {
	err := u.pgClient.Shutdown(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot close pgclient")
	}

	err = u.container.Terminate(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot terminate pgcontainer")
	}
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (u *UserTestSuite) TestUserRepo() {
	userRepo, err := userrepo.New(u.pgClient)
	if err != nil {
		u.T().Error("cannot create user repo")
	}

	u.T().Run("get user", func(t *testing.T) {
		user, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		assert.Equal(u.T(), user.Username, "batman")
		assert.Len(u.T(), user.Groups, 2)
		assert.Equal(u.T(), user.Groups[0].Name, "admins")
	})

	u.T().Run("update user", func(t *testing.T) {
		user, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		// remove one group from batman user
		user.Groups = []entity.Group{user.Groups[1]}

		_, err = userRepo.Update(context.Background(), user)
		if err != nil {
			u.T().Error(err)
		}

		user1, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		assert.Len(u.T(), user1.Groups, 1)
	})

	u.T().Run("update user2", func(t *testing.T) {
		user, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		// add another group
		var id int32 = 3
		user.Groups = append(user.Groups, entity.Group{ID: &id, Name: "editor"})

		_, err = userRepo.Update(context.Background(), user)
		if err != nil {
			u.T().Error(err)
		}

		user1, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		assert.Len(u.T(), user1.Groups, 2)
		assert.Equal(u.T(), user1.Groups[1].ID, int32(3))
	})

	u.T().Run("update user3", func(t *testing.T) {
		user, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		user.CanShare = true

		_, err = userRepo.Update(context.Background(), user)
		if err != nil {
			u.T().Error(err)
		}

		user1, err := userRepo.Get(context.Background(), "batman")
		if err != nil {
			u.T().Error(err)
		}

		assert.True(u.T(), user1.CanShare)
	})

	u.T().Run("create user", func(t *testing.T) {
		user := entity.User{
			Username: "bob",
			Role:     entity.RoleAdmin,
			UserID:   "id",
			CanShare: true,
		}

		id, err := userRepo.Create(context.Background(), user)
		if err != nil {
			u.T().Error(err)
		}

		assert.Equal(u.T(), id, 3)
	})
}

func getParentFolder() (string, error) {
	cwFolder, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// start removing folder till we reach the parentFolder
	folder := cwFolder

	for {
		folder = path.Dir(folder)
		if folder == "" || folder == "/" {
			break
		}

		// check if last folder is the parentFolder
		_, f := path.Split(folder)

		if f == parentFolder {
			return folder, nil
		}
	}

	return "", errors.New("parent folder not found") //nolint:goerr113
}
