package godbal_test

import (
	"github.com/xujiajun/godbal"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

const (
	defaultDriver = "mysql"
)

func TestNewMysql(t *testing.T) {
	database := godbal.NewMysql("")
	db, _, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()

	database.SetDB(db)
}

func TestNewDriveManager(t *testing.T) {
	newDriveManager := godbal.NewDriveManager()
	drivers := newDriveManager.GetAvailableDrivers()

	if driver, ok := drivers[defaultDriver]; ok {
		if driver != defaultDriver {
			t.Errorf("returned unexpected driver: got %v want %v", driver, defaultDriver)
		}
	} else {
		t.Errorf("an error when TestNewDriveManager driver %s not found ", defaultDriver)
	}
}

func TestDriveManager_GetMysqlConnection(t *testing.T) {
	newDriveManager := godbal.NewDriveManager()
	mysqlDB := newDriveManager.GetMysqlDB("")

	db, _, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()

	mysqlDB.SetDB(db)
}
