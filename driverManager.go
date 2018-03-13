package godbal

import "github.com/xujiajun/godbal/driver/mysql"

type DriveManager struct {
	drivers map[string]string
}

// NewDriveManager returns a newly initialized NewDriveManager that implements DriveManager
func NewDriveManager() *DriveManager {
	return &DriveManager{
		drivers: map[string]string{"mysql": "mysql"},
	}
}

// GetAvailableDrivers returns available drivers
func (driverManager *DriveManager) GetAvailableDrivers() map[string]string {
	return driverManager.drivers
}

// GetMysqlDB returns mysql Database
func (driverManager *DriveManager) GetMysqlDB(dataSourceName string) *mysql.Database {
	return NewMysql(dataSourceName)
}

// NewMysql is short for GetMysqlDB
func NewMysql(dataSourceName string) *mysql.Database {
	return mysql.New(dataSourceName)
}
