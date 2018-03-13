package main

import (
	"encoding/json"
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, err := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()
	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)
	rows, _ := queryBuilder.Select("uid,username,price,flag").From("userinfo", "").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").Query()

	jsonString, _ := json.Marshal(&rows)
	fmt.Print(string(jsonString))

	queryBuilder2 := mysql.NewQueryBuilder(database)
	sql := queryBuilder2.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParameter("johnny2").SetParameter("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").GetSQL()

	fmt.Print(sql)
	fmt.Print("\n")

	queryBuilder3 := mysql.NewQueryBuilder(database)
	rows2, _ := queryBuilder3.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParameter("johnny2").SetParameter("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").Query()

	jsonString2, _ := json.Marshal(&rows2)
	fmt.Print(string(jsonString2))

}
