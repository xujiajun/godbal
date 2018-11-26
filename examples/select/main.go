package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, err := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)
	sql := queryBuilder.Select("uid,username,price,flag").From("userinfo", "").SetFirstResult(0).
		SetMaxResults(5).OrderBy("uid", "DESC").OrderBy("username", "DESC").GetSQL()

	fmt.Println(sql)

	rows, _ := queryBuilder.QueryAndGetMap()

	jsonString, _ := json.Marshal(&rows)
	fmt.Println(string(jsonString))
	queryBuilder2 := mysql.NewQueryBuilder(database)
	sql = queryBuilder2.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParam("johnny2").SetParam("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").GetSQL()

	fmt.Println(sql)
	rows2, _ := queryBuilder2.QueryAndGetMap()
	jsonString2, _ := json.Marshal(&rows2)
	fmt.Println(string(jsonString2))
	queryBuilder3 := mysql.NewQueryBuilder(database)
	rows3, _ := queryBuilder3.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ?").
		SetParam("johnny2").SetFirstResult(0).
		SetMaxResults(1).OrderBy("uid", "DESC").QueryAndGetMap()

	fmt.Println()

	jsonString3, _ := json.Marshal(&rows3)
	fmt.Println(string(jsonString3))

}
