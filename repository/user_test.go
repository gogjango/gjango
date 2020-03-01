package repository_test

import (
	"testing"

	"github.com/calvinchengx/gin-go-pg/mockgopg"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserUnitTestSuite struct {
	suite.Suite
	mock     *mockgopg.SQLMock
	u        *model.User
	userRepo *repository.UserRepo
}

func (suite *UserUnitTestSuite) SetupTest() {
	var err error
	var db orm.DB
	db, suite.mock, err = mockgopg.NewGoPGDBTest()
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.u = &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "91919191",
	}

	log, _ := zap.NewDevelopment()
	suite.userRepo = repository.NewUserRepo(db, log)
}

func (suite *UserUnitTestSuite) TearDownTest() {
	suite.mock.FlushAll()
}

func TestUserUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UserUnitTestSuite))
}
