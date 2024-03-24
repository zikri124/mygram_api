package repository

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/zikri124/mygram-api/internal/infrastructure/mock"
	"github.com/zikri124/mygram-api/pkg/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockGorm() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestGetUser(t *testing.T) {
	t.Run("error get user by id", func(t *testing.T) {
		db, mock := newMockGorm()

		postgresMock := mocks.NewGormPostgres(t)
		postgresMock.On("GetConnection").Return(db)

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "users" WHERE id=1 AND deleted_at IS NULL
		`)).WillReturnError(errors.New("some error"))

		userRepo := userRepositoryImpl{db: postgresMock}
		res, err := userRepo.GetUserById(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Equal(t, uint32(0), res.ID)
	})

	t.Run("success get user by id", func(t *testing.T) {
		db, mock := newMockGorm()

		postgresMock := mocks.NewGormPostgres(t)
		postgresMock.On("GetConnection").Return(db)

		dobTime, err := helper.ParseStrToTime("2000-12-09")
		if err != nil {
			log.Fatalf("an error was not expected when parse str time to time : '%s'", err)
		}

		row := sqlmock.
			NewRows([]string{"id", "username", "email", "password", "dob"}).
			AddRow(1, "username", "test@test.com", "pppp", dobTime)

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL
		`)).WillReturnRows(row)

		userRepo := userRepositoryImpl{db: postgresMock}
		res, err := userRepo.GetUserById(context.Background(), 1)
		assert.Nil(t, err)
		assert.Equal(t, uint32(1), res.ID)
	})
}
