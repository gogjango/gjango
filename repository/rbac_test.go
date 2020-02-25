package repository_test

import (
	"fmt"
	"net/http/httptest"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type RBACTestSuite struct {
	suite.Suite
	db       *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
}

func (suite *RBACTestSuite) SetupTest() {
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
	createSchema(suite.db, &model.Company{}, &model.Location{}, &model.Role{}, &model.User{}, &model.Verification{})
}

func (suite *RBACTestSuite) TearDownTest() {
	suite.postgres.Stop()
}

func TestRBACTestSuite(t *testing.T) {
	suite.Run(t, new(RBACTestSuite))
}

func (suite *RBACTestSuite) TestRBAC() {
	// create a context for tests
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Set("role", model.SuperAdminRole)

	// create a user in our test database, which is superadmin
	log, _ := zap.NewDevelopment()
	userRepo := repository.NewUserRepo(suite.db, log)
	accountRepo := repository.NewAccountRepo(suite.db, log)
	u := createUserAndMakeActive(accountRepo, userRepo)
	fmt.Println("#######")
	fmt.Println(u.Role)
	fmt.Println("#######")

	rbac := repository.NewRBACService(userRepo)
	assert.NotNil(suite.T(), rbac)
}

func createUserAndMakeActive(accountRepo *repository.AccountRepo, userRepo *repository.UserRepo) *model.User {
	u := &model.User{
		CountryCode: "+65",
		Mobile:      "91919191",
	}
	u, _ = accountRepo.Create(u)
	u.Active = true
	u, _ = userRepo.Update(u)
	return u
}
