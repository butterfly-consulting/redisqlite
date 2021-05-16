function r {
    redis-cli "$@" 2>/dev/null
}
function re {
    echo ">" "$@"
    r "$@"
}

r sqlexec "drop table t" >/dev/null
r sqlexec "create table t(i int, s varchar)" >/dev/null
re sqlexec "insert into t(s,i) values('a',1),('b',2)"
re sqlexec "insert into t(s,i) values(?,?)" c 3
re sqlmap 0 "select * from t"
re sqlarr 0 "select * from t"
sel=$(r sqlprep 'select * from t where s=?')
ins=$(r sqlprep 'insert into t(i,s) values(?,?)')
#echo "id=$id"
echo "> sqlmap and sqlarr with prep"
r sqlmap 0 $sel c
r sqlexec $ins 4 d
r sqlarr 0 $sel 4
r sqlprep $sel
r sqlprep $ins
re sqlexec "begin ; insert into t(i,s) values(5,'e'); commit"
re sqlmap 0 "select * from t"
re sqlexec "begin ; insert into t(i,s) values(6,'f'); rollback"
re sqlmap 0 "select * from t"

