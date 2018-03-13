package main

import (
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	connection, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	queryBuilder := mysql.NewQueryBuilder(connection)

	affect, _ := queryBuilder.Delete("userinfo").Where("uid=?").SetParameter(7).PrepareAndExecute()

	fmt.Print(affect)
}
