package e2e_test

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/e2e"
	"github.com/calvinchengx/gin-go-pg/testhelper"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	db       *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
	helper   testhelper.TestHelper
}

func (suite *E2ETestSuite) SetupSuite() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, "tmp2")
	testConfig := embeddedpostgres.DefaultConfig().
		Username("db_test_user").
		Password("db_test_password").
		Database("db_test_database").
		Version(embeddedpostgres.V12).
		RuntimePath(tmpDir).
		Port(9877)

	suite.postgres = embeddedpostgres.NewDatabase(testConfig)
	_ = suite.postgres.Start()

	suite.db = pg.Connect(&pg.Options{
		Addr:     "localhost:9877",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})

	models := e2e.GetModels()
	suite.helper.CreateSchema(suite.db, models...)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.postgres.Stop()
}

func (suite *E2ETestSuite) TestGetModels() {
	models := e2e.GetModels()

	sql := `select count(*) from information_schema.tables where table_schema = 'public';`
	var count int
	res, err := suite.db.Query(pg.Scan(&count), sql, nil)

	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(models), count)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
