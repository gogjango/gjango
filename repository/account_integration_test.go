package repository_test

import (
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AccountTestSuite struct {
	suite.Suite
	db       *pg.DB
	dbErr    *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
}

func (suite *AccountTestSuite) SetupTest() {
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

func (suite *AccountTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func (suite *AccountTestSuite) TestAccountCreateAndVerify() {
	cases := []struct {
		name       string
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.User
	}{}

	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			accountRepo := repository.NewAccountRepo(tt.db, log)
			_, err := accountRepo.CreateAndVerify(tt.user)
			assert.Equal(t, tt.wantError, err)
			// if u != nil {
			// 	assert.Equal(t, tt.wantResult.Email, u.Email)
			// } else {
			// 	assert.Nil(t, u)
			// }
		})
	}
}

func (suite *AccountTestSuite) TestAccountCreate() {
	cases := []struct {
		name       string
		user       *model.User
		db         *pg.DB
		wantError  error
		wantResult *model.User
	}{
		{
			name: "Success: user created",
			user: &model.User{
				Email: "user@example.org",
			},
			db:        suite.db,
			wantError: nil,
			wantResult: &model.User{
				Email: "user@example.org",
			},
		},
		{
			name: "Failure: user already exists",
			user: &model.User{
				Email: "user@example.org",
			},
			db:         suite.db,
			wantError:  apperr.New(http.StatusBadRequest, "User already exists."),
			wantResult: nil,
		},
		{
			name: "Failure: db connection failed",
			db:   suite.dbErr,
			user: &model.User{
				Email: "user2@example.org",
			},
			wantError:  apperr.DB,
			wantResult: nil,
		},
		{
			name: "Failure",
			db:   suite.db,
			user: &model.User{
				ID:    1,
				Email: "user2@example.org",
			},
			wantError:  apperr.DB,
			wantResult: nil,
		},
	}

	for _, tt := range cases {
		suite.T().Run(tt.name, func(t *testing.T) {
			log, _ := zap.NewDevelopment()
			accountRepo := repository.NewAccountRepo(tt.db, log)
			u, err := accountRepo.Create(tt.user)
			assert.Equal(t, tt.wantError, err)
			if u != nil {
				assert.Equal(t, tt.wantResult.Email, u.Email)
			} else {
				assert.Nil(t, u)
			}
		})
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func createSchema(db *pg.DB, models ...interface{}) {
	for _, model := range models {
		opt := &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}
		err := db.CreateTable(model, opt)
		if err != nil {
			log.Fatal(err)
		}
	}
}
