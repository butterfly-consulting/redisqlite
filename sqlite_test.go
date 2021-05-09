package redisqlite

import (
	"fmt"
	"os"
	"testing"
)

func ExampleQuery() {
	fmt.Println(1, Exec("create table t(i int)"))
	fmt.Println(2, Exec("insert into t(i) values(1),(2),(3)"))
	res, err := Query("select * from t", 1)
	fmt.Println(3, err, res)
	res, err = Query("select * from t", 2)
	fmt.Println(4, err, res)
	res, err = Query("select * from t", 0)
	fmt.Println(5, err, res)
	// Output:
	// 1 <nil>
	// 2 <nil>
	// 3 <nil> [{"i":1}]
	// 4 <nil> [{"i":1},{"i":2}]
	// 5 <nil> [{"i":1},{"i":2},{"i":3}]
}

func TestMain(m *testing.M) {
	os.Remove("./sqlite.db")
	Open()
	os.Exit(m.Run())
}
