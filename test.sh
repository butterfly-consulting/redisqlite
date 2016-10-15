#!/usr/bin/env bash

r(){
#    echo -n "@$BASH_SOURCE:${BASH_LINENO[-2]}:"
    local expected=`sed "$((BASH_LINENO[-2] + 1))q;d" $BASH_SOURCE`

    echo '$' "$@"
    if [[ "$expected" =~ ^#\> ]];then
        expected=`echo "$expected" | sed -re 's/^#>\s*//'`
        local actual="$(redis-cli --no-raw "$@")"
        echo $actual
        if [[ "$expected" != "$actual" ]]; then
            echo -e "FAILED: Expected $expected"
        fi
    else
        redis-cli "$@"
    fi

}

r flushall

# Test hgetset
r hset a a 1
#> (integer) 1
r hgetset a a 2
#> "1"
r hget a a
#> "2"
# Return nil if field not exists
r hgetset a b 2
#> (nil)
r hgetset a b 3
#> "2"
# Will create a key
r hgetset b a 1
#> (nil)
r hgetset b a 2
#> "1"
r flushall

# Test hgetdel
r hgetdel a a
#> (nil)
r hset a a 1
r hgetdel a a
#> "1"

# TODO
