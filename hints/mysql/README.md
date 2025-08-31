```
sudo systemctl start docker
sudo systemctl status docker

sudo docker run --name x-mysql -e MYSQL_ROOT_PASSWORD=x -d mysql:latest
# -v /my/own/datadir:/var/lib/mysql
sudo docker stop x-mysql
sudo docker ps -a
sudo docker rm x-mysql

sudo docker exec -it x-mysql bash
mysql -u root -p

craete database db;
use db;
create table data(x int);

sudo docker stop $(sudo docker ps -a -q)
sudo docker rm $(sudo docker ps -a -q)

sudo docker container prune
sudo docker image prune
sudo docker system prune

sudo systemctl stop docker.socket
sudo systemctl stop docker.service
```

How to edit configuration

```
microdnf install -y vim
microdnf install -y iproute  # install ss
vim /etc/my.cnf
```

How to connect outside. And without password

```
docker run --name x-mysql -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_ROOT_PASSWORD= -d mysql:latest
/opt/homebrew/Cellar/mysql-client/9.4.0/bin/mysql -h 127.0.0.1 -P 3306 -u root
```

Ideas

```
https://github.com/nitso/colour-mysql-console
brew install grc
https://github.com/dbcli/mycli
brew install mycli

brew install mysql-shell
```

Useful

```
show processlist;
show full processlist;

-- table/columns summary
select * from information_schema.tables t inner join information_schema.columns c on t.table_schema = c.table_schema and t.table_name = c.table_name where t.table_schema="db";

-- online status
select * from performance_schema.users;
select * from performance_schema.hosts;
select * from performance_schema.socket_instances;

select * from performance_schema.data_lock_waits;
select * from performance_schema.data_locks;
```

Tricks

```
create table jsn (j json, id int generated always as (j->"$.id"));
insert into jsn (j) values ('{"id":1}');
insert into jsn (j) values ('{}');
select * from jsn;
+-----------+------+
| j         | id   |
+-----------+------+
| {"id": 1} |    1 |
| {}        | NULL |
+-----------+------+
```

Mistakes: multiply indexes allowed

```
create table test (id int not null primary key, unique(id), index(id)) engine=InnoDB;
select i.name, t.name from information_schema.innodb_indexes i inner join information_schema.innodb_tables t on i.table_id = t.table_id where t.name like '%/test';
+---------+---------+
| name    | name    |
+---------+---------+
| PRIMARY | db/test |
| ID      | db/test |
| ID_2    | db/test |
+---------+---------+
```

Check index usage:

```
select * from test where id=1;
select * from performance_schema.table_io_waits_summary_by_index_usage where OBJECT_NAME='test';
+-------------+---------------+-------------+------------+-----------
| OBJECT_TYPE | OBJECT_SCHEMA | OBJECT_NAME | INDEX_NAME | COUNT_STAR
+-------------+---------------+-------------+------------+-----------
| TABLE       | db            | test        | PRIMARY    |          1 <---
| TABLE       | db            | test        | ID         |          0
| TABLE       | db            | test        | ID_2       |          0
+-------------+---------------+-------------+------------+-----------
```
