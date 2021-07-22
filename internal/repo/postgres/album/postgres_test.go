package album_test

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
	"github.com/tupyy/gophoto/internal/repo/postgres/album"
	"github.com/tupyy/gophoto/utils/pgclient"
	"github.com/tupyy/gophoto/utils/pgtestcontainer"
)

const (
	// name of the root folder for the project.
	parentFolder = "gphotos"
	// sql setup file relative to parent folder.
	setupFile = "sql/setup/02_setup.sql"
	// fixtures file relative to parent folder.
	fixtureFile = "sql/fixtures/fixtures.sql"
)

type AlbumTestSuite struct {
	suite.Suite
	container testcontainers.Container
	pgClient  pgclient.Client
	repo      *album.AlbumPostgresRepo
}

func (as *AlbumTestSuite) TestGetAllAlbums() {
	asserter := assert.New(as.T())

	entities, err := as.repo.Get(context.Background())
	asserter.Nil(err)
	asserter.Len(entities, 8)
}

func (as *AlbumTestSuite) TestGetAlbumByID() {
	asserter := assert.New(as.T())

	ent, err := as.repo.GetByID(context.Background(), 1)
	asserter.Nil(err)
	asserter.Len(ent.UserPermissions, 1)
	asserter.Len(ent.GroupPermissions, 3)
	asserter.Equal(int32(1), ent.OwnerID)

	_, err = as.repo.GetByID(context.Background(), 100)
	asserter.NotNil(err)
}

func (as *AlbumTestSuite) TestGetAlbumByOwnerID() {
	asserter := assert.New(as.T())

	ent, err := as.repo.GetByOwnerID(context.Background(), 1)
	asserter.Nil(err)
	asserter.Len(ent, 4)

	_, err = as.repo.GetByOwnerID(context.Background(), 100)
	asserter.NotNil(err)
}

func (as *AlbumTestSuite) SetupSuite() {
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

	as.pgClient, err = c.GetInitialClient(ctx)
	if err != nil {
		panic(err)
	}

	as.container = c.Container

	albumRepo, err := album.NewPostgresRepo(as.pgClient)
	if err != nil {
		panic(err)
	}

	as.repo = albumRepo
}

func (as *AlbumTestSuite) TearDownSuite() {
	err := as.pgClient.Shutdown(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot close pgclient")
	}

	err = as.container.Terminate(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot terminate pgcontainer")
	}
}

func TestAlbumTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumTestSuite))
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

type AlbumTestSuite1 struct {
	suite.Suite
	container testcontainers.Container
	pgClient  pgclient.Client
	repo      *album.AlbumPostgresRepo
}

func (as *AlbumTestSuite1) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parentFolder, err := getParentFolder()
	if err != nil {
		panic(err)
	}

	initMap := make(map[string]string)
	setupKey := path.Join(parentFolder, setupFile)
	initMap[setupKey] = "/docker-entrypoint-initdb.d/setup.sql"

	fixtureKey := path.Join(parentFolder, "sql/fixtures/album_test1.sql")
	initMap[fixtureKey] = "/docker-entrypoint-initdb.d/zz_fixtures.sql"

	c, err := pgtestcontainer.NewPostgreSQLContainer(ctx, pgtestcontainer.PostgreSQLContainerRequest{
		BindMounts: initMap,
	})
	if err != nil {
		panic(err)
	}

	as.pgClient, err = c.GetInitialClient(ctx)
	if err != nil {
		panic(err)
	}

	as.container = c.Container

	albumRepo, err := album.NewPostgresRepo(as.pgClient)
	if err != nil {
		panic(err)
	}

	as.repo = albumRepo
}

func (as *AlbumTestSuite1) TearDownSuite() {
	err := as.pgClient.Shutdown(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot close pgclient")
	}

	err = as.container.Terminate(context.Background())
	if err != nil {
		logrus.WithError(err).Error("cannot terminate pgcontainer")
	}
}

func TestAlbumTestSuite1(t *testing.T) {
	suite.Run(t, new(AlbumTestSuite1))
}

func (as *AlbumTestSuite1) TestCreateAlbum() {
	asserter := assert.New(as.T())

	ptrString := func(s string) *string {
		return &s
	}

	album := entity.Album{
		Name:        "name",
		CreatedAt:   time.Now(),
		OwnerID:     1,
		Description: ptrString("test"),
		Location:    ptrString("test"),
		UserPermissions: map[int32][]entity.Permission{
			3: {entity.PermissionDeleteAlbum},
			4: {entity.PermissionReadAlbum, entity.PermissionEditAlbum},
		},
		GroupPermissions: map[int32][]entity.Permission{
			1: {entity.PermissionReadAlbum},
			2: {entity.PermissionDeleteAlbum},
		},
	}

	id, err := as.repo.Create(context.Background(), album)
	asserter.Nil(err)
	asserter.Greater(id, int32(-1))

	a1, err := as.repo.GetByID(context.Background(), id)
	asserter.Nil(err)
	asserter.Equal(a1.Name, "name")
}
