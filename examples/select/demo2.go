package main

import (
	"encoding/json"
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	queryBuilder := mysql.NewQueryBuilder(database)
	sql := queryBuilder.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParameter("johnny2").SetParameter("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").GetSQL()

	fmt.Print(sql)
	fmt.Print("\n")

	queryBuilder2 := mysql.NewQueryBuilder(database)
	rows, _ := queryBuilder2.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParameter("johnny2").SetParameter("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").Query()

	jsonString, _ := json.Marshal(&rows)
	fmt.Print(string(jsonString))

}
