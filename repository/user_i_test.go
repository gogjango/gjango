package repository_test

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserTestSuite struct {
	suite.Suite
	db       *pg.DB
	dbErr    *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
}

func (suite *UserTestSuite) SetupTest() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, "tmp")
	testConfig := embeddedpostgres.DefaultConfig().
		Username("db_test_user").
		Password("db_test_password").
		Database("db_test_database").
		Version(embeddedpostgres.V12).
		RuntimePath(tmpDir).
		Port(9876)

	suite.postgres = embeddedpostgres.NewDatabase(testConfig)
	err := suite.postgres.Start()
	assert.Equal(suite.T(), err, nil)

	suite.db = pg.Connect(&pg.Options{
		Addr:     "localhost:9876",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})
	suite.dbErr = pg.Connect(&pg.Options{
		Addr:     "localhost:9875",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})
	createSchema(suite.db, &model.Company{}, &model.Location{}, &model.Role{}, &model.User{}, &model.Verification{})
}

func (suite *UserTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func (suite *UserTestSuite) TestUserView() {
	cases := []struct {
		name       string
		create     bool
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.Verification
	}{
		{
			name:   "Fail: user not found",
			create: false,
			user: &model.User{
				Email:       "user@example.org",
				CountryCode: "+65",
				Mobile:      "91919191",
			},
			db:        suite.db,
			wantError: apperr.NotFound,
		},
		{
			name:   "Success: view user, find user",
			create: true,
			user: &model.User{
				Email:       "user@example.org",
				CountryCode: "+65",
				Mobile:      "91919191",
			},
			db:        suite.db,
			wantError: nil,
		},
	}
	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			userRepo := repository.NewUserRepo(tt.db, log)

			if tt.create {
				accountRepo := repository.NewAccountRepo(tt.db, log)
				_, err := accountRepo.Create(tt.user)
				assert.Nil(t, err)
				u, err := userRepo.View(tt.user.ID)
				assert.Nil(t, err)
				assert.Equal(t, tt.user.Mobile, u.Mobile)
			} else {
				u, err := userRepo.View(tt.user.ID)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByMobile(tt.user.CountryCode, tt.user.Mobile)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
				u, err = userRepo.FindByEmail(tt.user.Email)
				assert.Nil(t, u)
				assert.Equal(t, tt.wantError, err)
			}
		})
	}
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
