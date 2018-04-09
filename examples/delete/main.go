package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	queryBuilder := mysql.NewQueryBuilder(database)

	affect, _ := queryBuilder.Delete("userinfo").Where("uid=?").SetParam(7).PrepareAndExecute()

	fmt.Println(affect)
}
