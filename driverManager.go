package godbal

import "github.com/xujiajun/godbal/driver/mysql"

type DriveManager struct {
	drivers map[string]string
}

func NewDriveManager() *DriveManager {
	return &DriveManager{
		drivers: map[string]string{"mysql": "mysql"},
	}
}

func (driverManager *DriveManager) GetAvailableDrivers() map[string]string {
	return driverManager.drivers
}

func (driverManager *DriveManager) GetMysqlDB(dataSourceName string) *mysql.Database {
	return NewMysql(dataSourceName)
}

func NewMysql(dataSourceName string) *mysql.Database {
	return mysql.New(dataSourceName)
}
