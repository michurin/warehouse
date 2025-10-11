#!/bin/bash

# dump 127.0.0.1 28004 root password payment_provider
host="$1"
port="$2"
user="$3"
pass="$4"
dbname="$5"

export MYSQL_PWD=$pass # -p$pass
tables="$(mysql -h $host -P $port -u $user $dbname -N -B -A -e 'show tables')"
#echo "TABLE=$tables"

for tbl in $tables
do
    echo "--------- TABLE $tbl ($(mysql -h $host -P $port -u $user $dbname -N -B -A -e 'select count(*) from '"$tbl")) {{{"
    mysql -h $host -P $port -u $user $dbname -Nsr -A -e 'show create table '"$tbl"
    mysql -h $host -P $port -u $user $dbname --table --binary-as-hex -A -e 'select * from '"$tbl"' limit 10'
    echo '--------- }}}'
    echo
done

echo '-- v''i:ft=sql:nowrap:foldmethod=marker:foldlevel=0:'
