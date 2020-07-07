<!--
 * @Descripttion:
 * @Author: weihaoyu
-->

# gin-frame

基于 go-gin 开发的 api 框架，封装各种常用组件
有疑问随时联系本人，验证信息放我github地址即可
QQ：444216978
微信：why444216978

# run

```
go run main.go

curl localhost:777/ping
```

# app.ini example:

```
[app]
env = development
port = 777
app_id = moments-server
product = gin-frame
module = gin-frame
```

# mysql.ini example:

```
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
```

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
