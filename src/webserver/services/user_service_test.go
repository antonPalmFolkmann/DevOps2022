package services_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"

	//"github.com/go-test/deep"

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

/* func (s *Suite) Test_CreateUser() {
	// Arrange
	var (
		id               = uint(1)
		username         = "Ronald Weasley"
		email            = "ginger6@hp.com"
		passwordUnhashed = "secrets"
		passwordHashed   = "7de38f3c3d3baa7ca58a366f09577586"

		sqlInsert = `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","pw_hash") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "users"."id"`
	)


	// Act
	s.mock.ExpectBegin() // begin transaction
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), username, email, passwordHashed).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
									AddRow(id))
	s.mock.ExpectCommit() // commit transaction

	err := s.UserService.CreateUser(username, email, passwordUnhashed)

	// Assert
	require.NoError(s.T(), err)
}

func (s *Suite) Test_ReadAllUsers()  {
	// Arrange
	rows := sqlmock.NewRows([]string{"id", "username", "email", "pw_hash"}).
		AddRow(1, "user1", "email1", "password1").
		AddRow(2, "user2", "email2", "password2").
		AddRow(3, "user3", "email3", "password3").
		AddRow(4, "user4", "email4", "password4").
		AddRow(5, "user5", "email5", "password5")

	// Act

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, pw_hash FROM "users"`)).
		WillReturnRows(rows)

	res, err := s.UserService.ReadAllUsers()

	// Assert
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(5, len(res)))
}

func (s *Suite) Test_ReadUserByID() {
	// Arrange
	var (
		id       = uint(3)
		username = "Harry Potter"
		email    = "tbwl@hp.com"
		password = "secrets"
	)

	// Act
	rows := sqlmock.NewRows([]string{"id", "username", "email", "pw_hash"}).
		AddRow(id, username, email, password)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, pw_hash FROM "users" WHERE (id = $1)`)).
		WithArgs(id).
		WillReturnRows(rows)

	res, err := s.UserService.ReadUserById(id)

	// Assert
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(storage.UserDTO{ID: id, Username: res.Username, Email: res.Email, PwHash: res.PwHash}, res))
}

func (s *Suite) Test_ReadUserByUsername() {
	// Arrange
	var (
		id       = uint(4)
		username = "Albus Dumbledore"
		email    = "oldman@hp.com"
		password = "secrets"
	)

	// Act
	rows := sqlmock.NewRows([]string{"id", "username", "email", "pw_hash"}).
		AddRow(id, username, email, password)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, pw_hash FROM "users" WHERE (username = $1)`)).
		WithArgs(username).
		WillReturnRows(rows)

	res, err := s.UserService.ReadUserByUsername(username)

	// Assert
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(storage.UserDTO{ID: id, Username: res.Username, Email: res.Email, PwHash: res.PwHash}, res))
}

func (s *Suite) Test_ReadUserIdByUsername() {
	// Arrange
	var (
		id       = uint(5)
		username = "Tom Riddle"
		email    = "nonose@hp.com"
		password = "secrets"
	)

	// Act
	rows := sqlmock.NewRows([]string{"id", "username", "email", "pw_hash"}).
		AddRow(id, username, email, password)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM "users" WHERE (username = $1)`)).
		WithArgs(username).
		WillReturnRows(rows)

	res, err := s.UserService.ReadUserIdByUsername(username)

	// Assert
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(id, res))
}  */

/* func (s *Suite) Test_UpdateUser() {

	// Arrange
	var (
		id               = uint(1)
		username         = "Ronald Weasley"
		email            = "ginger6@hp.com"
		passwordUnhashed = "secrets"
		passwordHashed   = "7de38f3c3d3baa7ca58a366f09577586"

		sqlInsert = `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","pw_hash") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "users"."id"`

		changed_id             = uint(2)
		changed_username       = "Ronild Waslib"
		changed_email          = "gingertop6@hp.com"
		changed_passwordHashed = "doublesecrets"

		sqlUpdate = `SELECT * FROM "users" WHERE ("id" = $1)`
	)

	// Insert
	s.mock.ExpectBegin() // begin transaction
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), username, email, passwordHashed).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit() // commit transaction

	// Update
	s.mock.ExpectBegin() // begin transaction
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectCommit() // commit transaction

	// Act
	err := s.UserService.CreateUser(username, email, passwordUnhashed)
	require.NoError(s.T(), err)
	
	err = s.UserService.UpdateUser(changed_id, changed_username, changed_email, changed_passwordHashed)

	// Assert
	require.NoError(s.T(), err)
} */

func (s *Suite) Test_UpdateUser() {

	// Arrange
	var (
		id               = uint(1)
		username         = "Ronald Weasley"
		email            = "ginger6@hp.com"
		passwordHashed   = "7de38f3c3d3baa7ca58a366f09577586"

		changed_id             = uint(2)
		changed_username       = "Ronild Waslib"
		changed_email          = "gingertop6@hp.com"
		//changed_passwordUnhashed = "doublesecrets"
		changed_passwordHashed = "doublesecrets"

		sqlUpdate = `UPDATE "users" SET ("id" = $1,"username" = $2,"email" = $3,"pw_hash" = $4) WHERE ("id" = $1) ` 
	)

	// Insert
	s.mock.NewRows([]string{"id", "username", "email", "pw_hash"}).AddRow(id, username, email, passwordHashed)

	// Update
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(changed_id, changed_username, changed_email, changed_passwordHashed).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(changed_id))
	s.mock.ExpectCommit()
	// Act
	
	err := s.UserService.UpdateUser(changed_id, changed_username, changed_email, changed_passwordHashed)

	// Assert
	require.NoError(s.T(), err)
}