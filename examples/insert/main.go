package main

import (
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	connection, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	queryBuilder := mysql.NewQueryBuilder(connection)

	id, _ := queryBuilder.Insert("userinfo").Set("username", "?").Set("departname", "?").Set("created", "?").
		SetParameter("johnny3").SetParameter("tec5").SetParameter("12312312").PrepareAndExecute()

	fmt.Print(id)
}
