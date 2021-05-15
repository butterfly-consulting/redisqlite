function r {
    redis-cli "$@" 2>/dev/null
}
function re {
    echo ">" "$@"
    r "$@"
}

r sqlexec "drop table t" >/dev/null
r sqlexec "create table t(i int, s varchar)" >/dev/null
re sqlexec "insert into t(i,s) values('a',1),('b',2),('c',3)"
re sqlmap 0 "select * from t"
re sqlarr 0 "select * from t"
id=$(r sqlprep 'select * from t where s=?')
#echo "id=$id"
echo "> sqlmap and sqlarr with prep"
r sqlmap 0 $id 3
r sqlarr 0 $id 3
r sqlprep $id
re sqlexec "begin ; insert into t(i,s) values('d',4); commit"
re sqlmap 0 "select * from t"
re sqlexec "begin ; insert into t(i,s) values('e',5); rollback"
re sqlmap 0 "select * from t"
