process.env["__NIM_REDIS_IP"]="127.0.0.1"
process.env["__NIM_REDIS_PASSWORD"]="password"
let nim = require("../../nimbella-sdk-nodejs")
r = nim.sqlite()
r.exec("drop table t", console.log)
r.exec("create table t(n int, s varchar)", console.log)
r.exec("insert into t(n,s) values(1,'a'),(2,'b'),(3,'c')", console.log)
r.query("select * from t", 0, console.log)
r.query("select * from t", 1, console.log)

x = r.queryAsync("select * from t", 0) //.then(x => {  console.log(x)  })
