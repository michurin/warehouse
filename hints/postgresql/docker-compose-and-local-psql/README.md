# run docker-compose

```
su -Pc 'sudo whoami' - a
su -Pc 'sudo docker-compose -f /tmp/do/run.yaml up' - a
```

```
http://localhost:8080/?pgsql=db&username=user&db=user&ns=public
```

# install psql from source and run it

```
curl https://ftp.postgresql.org/pub/source/v16.0/postgresql-16.0.tar.gz -O postgresql-16.0.tar.gz
tar xzf postgresql-16.0.tar.gz
cd postgresql-16.0
./configure
make
```

```
LD_PRELOAD=./src/interfaces/libpq/libpq.so.5 PGPASSWORD=pass ./src/bin/psql/psql -h localhost -U user
```
