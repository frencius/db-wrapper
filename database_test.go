package db

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockCtx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	mockConfig := &Configuration{
		Host:     "http://local",
		Port:     5432,
		User:     "username",
		Password: "!@$#!$",
		Driver:   "postgres",
		Name:     "dbtest",
	}

	t.Run("positive", func(t *testing.T) {
		_, err := New(mockCtx, mockConfig)
		assert.Nil(t, err, "should be nil")

	})

	t.Run("negative - sql open failed", func(t *testing.T) {
		// driver empty to make sql open failed
		mockConfig.Driver = ""

		_, err := New(mockCtx, mockConfig)
		assert.Error(t, err, "should be error")

	})
}

func TestClose(t *testing.T) {
	mockCtx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	t.Run("positive", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectClose()

		mockDB := &Database{
			Context:  mockCtx,
			Database: db,
		}

		err = mockDB.Close()
		assert.Nil(t, err, "should be nil")
	})

	t.Run("negative - close failed", func(t *testing.T) {
		// db close never expected
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mockDB := &Database{
			Context:  mockCtx,
			Database: db,
		}

		err = mockDB.Close()
		assert.Error(t, err, "should be error")
	})
}

func TestPing(t *testing.T) {
	mockCtx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	t.Run("positive", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectPing()

		mockDB := &Database{
			Context:  mockCtx,
			Database: db,
		}

		err = mockDB.Ping()
		assert.Nil(t, err, "should be nil")
	})

	t.Run("negative - ping failed", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectPing().WillReturnError(errors.New("error"))

		mockDB := &Database{
			Context:  mockCtx,
			Database: db,
		}

		err = mockDB.Ping()
		assert.Error(t, err, "should be error")
	})
}
