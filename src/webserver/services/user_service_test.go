package services_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	UserService services.User
	user        *storage.User
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.UserService = *services.NewUserService(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// ------------------- TESTS -------------------------

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_ReadUserByID() {
	// Arrange
	var (
		ID = 1
		Username = ""
		Email = ""
		PwHash = ""
		Messages = make([]storage.Message, 0)
		Follows = make([]storage.User, 0)
	)

	// Act
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "user" WHERE (id = $1)`)).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "email", "pw_hash", "messages", "follows"}).
			AddRow(ID, Username, Email, PwHash, Messages, Follows))

	res, err := s.UserService.ReadUserById(uint(ID))

	// Assert
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(&storage.User{Username: res.Username, Email: res.Email}, res))
}
