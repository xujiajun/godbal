package mysql_test

import (
	"testing"
	//"fmt"
	//"encoding/json"

	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	INNER = "INNER"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
)

func TestQueryBuilder_Select(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()
	database := godbal.NewMysql("")
	database.SetDB(db)

	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)
	sql := queryBuilder.Select("id, title").From("posts", "").GetSQL()

	expectedSql := "SELECT id, title FROM posts "

	if sql != expectedSql {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectQuery(sql).WillReturnRows(rows)

	_, err = queryBuilder.Query()

	if err != nil {
		t.Errorf("error '%s' was not expected, while SELECT a row", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)

	}
}

func testJoinCommon(t *testing.T, joinFlag string) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()
	database := godbal.NewMysql("")
	database.SetDB(db)

	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)

	sql, expectedSql := "", ""
	switch joinFlag {
	case INNER:
		sql = queryBuilder.Select("p.id, p.title").From("posts", "p").SetFirstResult(0).
			SetMaxResults(3).InnerJoin("user", "u", "u.uid = p.uid").GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p INNER JOIN user u ON u.uid = p.uid LIMIT 0,3"
	case LEFT:
		sql = queryBuilder.Select("p.id, p.title").From("posts", "p").SetFirstResult(0).
			SetMaxResults(3).LeftJoin("user", "u", "u.uid = p.uid").GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p LEFT JOIN user u ON u.uid = p.uid LIMIT 0,3"
	case RIGHT:
		sql = queryBuilder.Select("p.id, p.title").From("posts", "p").SetFirstResult(0).
			SetMaxResults(3).RightJoin("user", "u", "u.uid = p.uid").GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p RIGHT JOIN user u ON u.uid = p.uid LIMIT 0,3"
	}

	if sql != expectedSql {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectQuery(sql).WillReturnRows(rows)

	//rows2, err := queryBuilder.Query()
	_, err = queryBuilder.Query()

	//jsonString, _ := json.Marshal(&rows2)
	//fmt.Print(string(jsonString))

	if err != nil {
		t.Errorf("error '%s' was not expected, while SELECT a row", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQueryBuilder_InnerJoin(t *testing.T) {
	testJoinCommon(t, INNER)
}

func TestQueryBuilder_LeftJoin(t *testing.T) {
	testJoinCommon(t, LEFT)
}

func TestQueryBuilder_RightJoin(t *testing.T) {
	testJoinCommon(t, RIGHT)
}

func TestQueryBuilder_GetSQL(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()
	database := godbal.NewMysql("")
	database.SetDB(db)

	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)

	sql := queryBuilder.Select("uid,username,created,textVal,price,name").From("userinfo", "").Where("username = ? AND departname = ?").
		SetParam("johnny2").SetParam("tec").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").GetSQL()

	expectedSql := "SELECT uid,username,created,textVal,price,name FROM userinfo  WHERE username = ? AND departname = ? ORDER BY uid DESC LIMIT 0,3"

	if sql != expectedSql {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql)
	}

	//GET
	queryBuilder.GetParameter()
	queryBuilder.GetMaxResults()
	queryBuilder.GetFirstResult()

	queryBuilder2 := mysql.NewQueryBuilder(database)
	sql = queryBuilder2.Select("u.uid,u.username,p.address,count(*) as num").From("userinfo", "u").SetFirstResult(0).
		SetMaxResults(3).RightJoin("profile", "p", "u.uid = p.uid").Having("num > 1").GetSQL()

	expectedSql2 := "SELECT u.uid,u.username,p.address,count(*) as num FROM userinfo u RIGHT JOIN profile p ON u.uid = p.uid HAVING num > 1 LIMIT 0,3"
	if sql != expectedSql2 {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql2)
	}
}

func TestQueryBuilder_Transaction(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database := godbal.NewMysql("")
	database.SetDB(db)

	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)

	sql := queryBuilder.Update("userinfo", "u").Set("u.username", "johnny").Where("u.uid=?").
		SetParam(1).GetSQL()

	expectedSql := "UPDATE userinfo u SET u.username = ?  WHERE u.uid=?"

	if sql != expectedSql {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql)
	}

	queryBuilder2 := mysql.NewQueryBuilder(database)

	sql2 := queryBuilder2.Insert("userinfo").Value("username", "johnny3").Value("departname", "tec5").GetSQL()

	expectedSql2 := "INSERT INTO userinfo (username,departname) VALUES(?,?)"

	if sql2 != expectedSql2 {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql2, expectedSql2)
	}

	// expect transaction begin
	mock.ExpectBegin()
	// expect user balance update
	mock.ExpectPrepare(sql).ExpectExec().
		WithArgs("xujiajun", 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row

	mock.ExpectPrepare("INSERT INTO userinfo ").ExpectExec().
		WithArgs("tec", "xxx").
		WillReturnResult(sqlmock.NewResult(1, 1)) // no insert id, 1 affected row

	// expect a transaction commit
	mock.ExpectCommit()

	if err != nil {
		t.Errorf("An error '%s' was not expected while queryBuilder PrepareAndExecute", err)
	}

	defer database.Close()

	transaction, err := database.Begin()

	if err != nil {
		t.Errorf("An error '%s' was not expected while database begin", err)
	}

	_, err = transaction.Tx.Exec(sql, "xujiajun", 1)
	_, err = transaction.Tx.Exec(sql2, "tec", "xxx")

	if err != nil {
		return
	}

	if err := transaction.Commit(); err != nil {
		return
	}
}

func TestQueryBuilder_Delete(t *testing.T) {
	db, _, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mysql database connection", err)
	}

	defer db.Close()
	database := godbal.NewMysql("")
	database.SetDB(db)

	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)

	sql := queryBuilder.Delete("userinfo").Where("uid=?").SetParam(7).GetSQL()

	expectedSql := "DELETE  FROM userinfo WHERE uid=?"

	if sql != expectedSql {
		t.Errorf("returned unexpected sql: got %v want %v",
			sql, expectedSql)
	}
}
