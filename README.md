# gin-frame
基于go-gin开发的api框架，封装各种常用组件

# run
go run main.go 

curl localhost:777/ping

# mysql.ini example:
[default]
host = 127.0.0.1
user = why
password = why123
port = 3306
db = why
charset = utf8
is_log = true
max_open = 8
max_idle = 4

# redis.ini example:
```
[default]
host = 127.0.0.1
port = 6379
db = 0
auth =
max_active = 600
max_idle = 10
is_log = true
```

# log.ini example:
```
[run]
run_dir = ./logs/run/
dir = ./logs/run/
area = 1

[error]
error_dir = ./logs/error/
dir = ./logs/error/
area = 1

[amqp]
amqp_dir = ./logs/amqp/
dir = ./amqp/

[mysql_open]
turn = true

[redis_open]
turn = true

[rabbitmq_open]
turn = true

[es_open]
turn = true
```

# es.ini example:
```
[default]
host = http://127.0.0.1
port = 9200
```


