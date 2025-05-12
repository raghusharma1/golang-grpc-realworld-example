package store

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreGetByEmail(t *testing.T) {

	tests := []struct {
		name          string
		isError       bool
		setupMock     func(mock sqlmock.Sqlmock, email string)
	}{
		{
			"Valid User Account Retrieval",
			false,
			func(mock sqlmock.Sqlmock, email string) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow("1", "test@example.com")
				mock.ExpectQuery("SELECT \\* FROM").WithArgs(email).WillReturnRows(rows)
			},
		},
		{
			"Non-existent User Account Retrieval",
			true,
			func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery("SELECT \\* FROM").WithArgs(email).WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			"Invalid email",
			true,
			func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery("SELECT \\* FROM").WithArgs(email).WillReturnError(gorm.ErrInvalidSQL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v", r)
					t.Fail()
				}
			}()

			db, mock, _ := sqlmock.New()
			dbInstance, _ := gorm.Open("mysql", db) // todo - replace "mysql" with the appropriate database driver name
			defer dbInstance.Close()
			userStore := UserStore{db: dbInstance}

			tt.setupMock(mock, "test@example.com") // using constant email for all tests, TODO: customize it as per requirement

			res, err := userStore.GetByEmail("test@example.com") // using constant email for all tests, TODO: customize it as per requirement

			if tt.isError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "test@example.com", res.Email)
			}

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

