package redisqlite

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func ExampleArgs() {
	count, last, err := Exec("create table ttt(k string, i int); insert into ttt values('a',1),('b',4),('c',3)", nil)
	fmt.Println(1, err)
	count, last, err = Exec("update ttt set i=? where k=?", []interface{}{2, "b"})
	fmt.Println(2, count, last, err)
	res, err := Query("select * from ttt", nil, true, 0)
	fmt.Println(3, err, len(res), res)
	res, err = Query("select * from ttt", nil, false, 0)
	fmt.Println(3.1, err, len(res), res)
	res, err = Query("select * from ttt where i=? and k=?", []interface{}{2, "b"}, true, 0)
	fmt.Println(4, err, len(res), res)
	res, err = Query("select * from ttt where i=? and k=?", []interface{}{2, "b"}, false, 0)
	fmt.Println(4.1, err, len(res), res)
	// Output:
	// 1 <nil>
	// 2 3 1 <nil>
	// 3 <nil> 3 [{"i":1,"k":"a"} {"i":2,"k":"b"} {"i":3,"k":"c"}]
	// 3.1 <nil> 3 [["a",1] ["b",2] ["c",3]]
	// 4 <nil> 1 [{"i":2,"k":"b"}]
	// 4.1 <nil> 1 [["b",2]]
}

func ExampleExecErrUpdate() {
	count, last, err := Exec("blabla", nil)
	fmt.Println(1, count, last, err)
	count, last, err = Exec("create table tt(k string, i int); insert into tt values('a',1),('b',3)", nil)
	fmt.Println(2, count, last, err)
	count, last, err = Exec("update tt set i=2 where k='b'", nil)
	fmt.Println(3, count, last, err)
	res, err := Query("select * from tt", nil, true, 0)
	fmt.Println(4, err, len(res), res)
	// Output:
	// 1 -1 -1 near "blabla": syntax error
	// 2 2 2 <nil>
	// 3 2 1 <nil>
	// 4 <nil> 2 [{"i":1,"k":"a"} {"i":2,"k":"b"}]
}

func ExampleQueryExec() {
	count, last, err := Exec("create table t(i int)", nil)
	fmt.Println(1, err)
	count, last, err = Exec("insert into t(i) values(1),(2)", nil)
	fmt.Println(2, count, last, err)
	count, last, err = Exec("insert into t(i) values(3)", nil)
	fmt.Println(3, count, last, err)
	res, err := Query("select * from t", nil, true, 1)
	fmt.Println(4, err, len(res), res)
	res, err = Query("select * from t", nil, false, 1)
	fmt.Println(4.1, err, len(res), res)
	res, err = Query("select * from t", nil, true, 2)
	fmt.Println(5, err, len(res), res)
	res, err = Query("select * from t", nil, false, 2)
	fmt.Println(5.1, err, len(res), res)
	res, err = Query("select * from t", nil, true, 0)
	fmt.Println(6, err, len(res), res)
	res, err = Query("select * from t", nil, false, 0)
	fmt.Println(6.1, err, len(res), res)
	// Output:
	// 1 <nil>
	// 2 2 2 <nil>
	// 3 3 1 <nil>
	// 4 <nil> 1 [{"i":1}]
	// 4.1 <nil> 1 [[1]]
	// 5 <nil> 2 [{"i":1} {"i":2}]
	// 5.1 <nil> 2 [[1] [2]]
	// 6 <nil> 3 [{"i":1} {"i":2} {"i":3}]
	// 6.1 <nil> 3 [[1] [2] [3]]
}

func ExamplePrep() {
	// prepare
	bad, err := Prep("blabla")
	fmt.Println(1, bad, err)
	bad, err = Prep("9999")
	fmt.Println(1.1, bad, err)
	crt, err := Prep("create table tttt(k string, i int)")
	fmt.Println(2, crt, err)
	_, _, err = Exec(strconv.Itoa(crt), nil)
	fmt.Println(3, err)

	sel, err := Prep("select k from tttt where i=?")
	fmt.Println(4, sel >= 0, err)
	ins, err := Prep("insert into tttt values(?,?)")
	fmt.Println(5, ins >= 0, err)

	// insert
	_, _, err = Exec(strconv.Itoa(ins), nil)
	fmt.Println(6, err)
	count, lastId, err := Exec(strconv.Itoa(ins), []interface{}{"a", 1})
	fmt.Println(7, count, lastId, err)
	count, lastId, err = Exec(strconv.Itoa(ins), []interface{}{"b", 2})
	fmt.Println(8, count, lastId, err)

	// select
	_, err = Query(strconv.Itoa(sel), nil, true, 0)
	fmt.Println(9, err)
	_, err = Query(strconv.Itoa(sel), []interface{}{"b", 2}, true, 0)
	fmt.Println(10, err)
	res, err := Query(strconv.Itoa(sel), []interface{}{2}, true, 0)
	fmt.Println(11, err, res)

	// unprep
	ins1, err := Prep(strconv.Itoa(sel))
	fmt.Println(12, ins1, res)
	sel1, err := Prep(strconv.Itoa(ins))
	fmt.Println(13, sel1, res)

	// check no prepared statement
	_, _, err = Exec(strconv.Itoa(sel), nil)
	fmt.Println(14, err)
	_, _, err = Exec("999", nil)
	fmt.Println(15, err)
	_, err = Query(strconv.Itoa(ins), nil, true, 0)
	fmt.Println(16, err)
	_, err = Query("999", nil, true, 0)
	fmt.Println(17, err)
	// Output:
	// 1 -1 near "blabla": syntax error
	// 1.1 -1 invalid prepared statement index
	// 2 2 <nil>
	// 3 <nil>
	// 4 true <nil>
	// 5 true <nil>
	// 6 sql: expected 2 arguments, got 0
	// 7 1 1 <nil>
	// 8 2 1 <nil>
	// 9 sql: expected 1 arguments, got 0
	// 10 sql: expected 1 arguments, got 2
	// 11 <nil> [{"k":"b"}]
	// 12 -1 [{"k":"b"}]
	// 13 -1 [{"k":"b"}]
	// 14 no such prepared statement index
	// 15 no such prepared statement index
	// 16 no such prepared statement index
	// 17 no such prepared statement index
}

func TestMain(m *testing.M) {
	os.Remove("./sqlite.db")
	Open()
	os.Exit(m.Run())
}
