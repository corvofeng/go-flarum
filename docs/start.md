
## How to start

### Init the project

```bash
git clone https://github.com/corvofeng/go-flarum.git
cd go-flarum

git submodule init
git submodule update
```


### Build and run

```bash

cd go-flarum
docker-compose build

docker-compose up
go-flarum_redis-db_1 is up-to-date
go-flarum_mysql-db_1 is up-to-date
Starting go-flarum_go-flarum-init_1 ... done
Creating go-flarum_go-flarum_1      ... done
Attaching to go-flarum_redis-db_1, go-flarum_mysql-db_1, go-flarum_go-flarum-init_1, go-flarum_go-flarum_1
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [logger.go:35] In Debug log level
go-flarum_1       | 2024-07-01 15:00:27 ▶ I [main.go:41] Git version: Development
go-flarum-init_1  |
go-flarum-init_1  | 2024/07/01 15:00:25 /home/go-flarum/cmd/migration/main.go:29 SLOW SQL >= 200ms
go-flarum-init_1  | [202.032ms] [rows:0] CREATE TABLE `tags` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`name` varchar(191),`url_name` varchar(191),`articles` bigint unsigned,`about` longtext,`parent_id` bigint unsigned,`position` bigint unsigned,`description` longtext,`hidden` boolean,`color` longtext,`icon_img` longtext,PRIMARY KEY (`id`),INDEX `idx_tags_deleted_at` (`deleted_at`),UNIQUE INDEX `idx_name` (`name`),UNIQUE INDEX `idx_urlname` (`url_name`))
go-flarum-init_1  |
go-flarum-init_1  | 2024/07/01 15:00:26 /home/go-flarum/cmd/migration/main.go:29 SLOW SQL >= 200ms
go-flarum-init_1  | [215.310ms] [rows:0] CREATE TABLE `topic_tags` (`topic_id` bigint unsigned,`tag_id` bigint unsigned,PRIMARY KEY (`topic_id`,`tag_id`),CONSTRAINT `fk_topic_tags_topic` FOREIGN KEY (`topic_id`) REFERENCES `topics`(`id`),CONSTRAINT `fk_topic_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`))
go-flarum-init_1  |
go-flarum-init_1  | 2024/07/01 15:00:26 /home/go-flarum/cmd/migration/main.go:31 SLOW SQL >= 200ms
go-flarum-init_1  | [260.099ms] [rows:0] CREATE TABLE `replies` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`topic_id` bigint unsigned,`user_id` bigint unsigned,`number` bigint unsigned,`content` longtext,`client_ip` longtext,`add_time` bigint unsigned,PRIMARY KEY (`id`),INDEX `idx_replies_deleted_at` (`deleted_at`))
go-flarum-init_1  |
go-flarum-init_1  | 2024/07/01 15:00:26 /home/go-flarum/cmd/migration/main.go:32 SLOW SQL >= 200ms
go-flarum-init_1  | [256.020ms] [rows:0] CREATE TABLE `reply_likes` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`user_id` bigint unsigned,`reply_id` bigint unsigned,PRIMARY KEY (`id`),INDEX `idx_reply_likes_user_id` (`user_id`),INDEX `idx_reply_likes_reply_id` (`reply_id`),INDEX `idx_reply_likes_deleted_at` (`deleted_at`))
go-flarum-init_1  | 2024-07-01 15:00:26 ▶ I [main.go:34] &{0xc00007c090 <nil> 0 0xc00041c000 1} Redis<redis-db:6379 db:0>
go-flarum-init_1  | 2024-07-01 15:00:26 ▶ I [main.go:64] Migrate the db
redis-db_1        | 1:C 01 Jul 2024 15:00:02.377 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
redis-db_1        | 1:C 01 Jul 2024 15:00:02.377 # Redis version=6.0.20, bits=64, commit=00000000, modified=0, pid=1, just started
redis-db_1        | 1:C 01 Jul 2024 15:00:02.377 # Warning: no config file specified, using the default config. In order to specify a config file use redis-server /path/to/redis.conf
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * Running mode=standalone, port=6379.
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 # Server initialized
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * Loading RDB produced by version 6.0.20
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * RDB age 14 seconds
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * RDB memory usage when created 0.77 Mb
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * DB loaded from disk: 0.000 seconds
redis-db_1        | 1:M 01 Jul 2024 15:00:02.378 * Ready to accept connections
go-flarum_go-flarum-init_1 exited with code 0
mysql-db_1        | 2024-07-01 15:00:02+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.0.37-1.el9 started.
mysql-db_1        | 2024-07-01 15:00:02+00:00 [Note] [Entrypoint]: Switching to dedicated user 'mysql'
mysql-db_1        | 2024-07-01 15:00:02+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.0.37-1.el9 started.
mysql-db_1        | '/var/lib/mysql/mysql.sock' -> '/var/run/mysqld/mysqld.sock'
mysql-db_1        | 2024-07-01T15:00:02.950789Z 0 [Warning] [MY-011068] [Server] The syntax '--skip-host-cache' is deprecated and will be removed in a future release. Please use SET GLOBAL host_cache_size=0 instead.
mysql-db_1        | 2024-07-01T15:00:02.952129Z 0 [System] [MY-010116] [Server] /usr/sbin/mysqld (mysqld 8.0.37) starting as process 1
mysql-db_1        | 2024-07-01T15:00:02.958017Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.
mysql-db_1        | 2024-07-01T15:00:03.337360Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.
mysql-db_1        | 2024-07-01T15:00:03.512556Z 0 [System] [MY-010229] [Server] Starting XA crash recovery...
mysql-db_1        | 2024-07-01T15:00:03.520832Z 0 [System] [MY-010232] [Server] XA crash recovery finished.
mysql-db_1        | 2024-07-01T15:00:03.698560Z 0 [Warning] [MY-010068] [Server] CA certificate ca.pem is self signed.
mysql-db_1        | 2024-07-01T15:00:03.698663Z 0 [System] [MY-013602] [Server] Channel mysql_main configured to support TLS. Encrypted connections are now supported for this channel.
mysql-db_1        | 2024-07-01T15:00:03.718257Z 0 [Warning] [MY-011810] [Server] Insecure configuration for --pid-file: Location '/var/run/mysqld' in the path is accessible to all OS users. Consider choosing a different directory.
mysql-db_1        | 2024-07-01T15:00:03.750870Z 0 [System] [MY-011323] [Server] X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /var/run/mysqld/mysqlx.sock
mysql-db_1        | 2024-07-01T15:00:03.750916Z 0 [System] [MY-010931] [Server] /usr/sbin/mysqld: ready for connections. Version: '8.0.37'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server - GPL.
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [core.go:88] Get redis db url: redis://:@redis-db:6379
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [core.go:99] PONG <nil>
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [core.go:101] Get mongo db url: mongodb://172.17.0.1
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [core.go:108] Get mysql db url: root:password@tcp(mysql-db:3306)/flarum?charset=utf8mb4&parseTime=True
go-flarum_1       | 2024-07-01 15:00:27 ▶ I [main.go:58] This is not a cron worker
go-flarum_1       | 2024-07-01 15:00:27 ▶ N [router.go:68] Init flarum router
go-flarum_1       | 2024-07-01 15:00:27 ▶ N [router.go:37] Init flarum admin router
go-flarum_1       | 2024-07-01 15:00:27 ▶ D [main.go:78] Web server Listen port 8082

```



## 本地使用docker构建

```bash
yarn build
CGO_ENABLED=0 GOOS=linux go build -o go-flarum
docker build . -t go-flarum:dev

docker stack deploy --compose-file docker-compose.yaml go-flarum
```

## 更新插件

```bash
git submodule foreach 'git pull'

git submodule foreach '[[ "$path" =~ (view/extensions/(analytics|auth-github|diff|flarum-pipetables))|(view/locale/zh) ]]  || git checkout main '

git submodule foreach '[[ "$path" =~ view/extensions/(analytics|auth-github|diff|flarum-pipetables) ]]  || echo '
git submodule foreach 'if [[ "$path" =~ view/extensions/(analytics|auth-github|diff|flarum-pipetables) ]]; then git checkout main; fi;'



git submodule foreach '
    if [[ "$path" =~ view/extensions/analytics ]]; then
        git fetch origin && git checkout 1.0.0
    fi

    if [[ "$path" =~ view/extensions/auth-github ]]; then
        git fetch origin && git checkout v0.1.0-beta.13
    fi
'

git submodule foreach 'if [[ "$path" =~ view/extensions/analytics ]]; then git fetch origin && git checkout v1.0.0; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/auth-github ]]; then git fetch origin && git checkout v0.1.0-beta.13; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/diff ]]; then git fetch origin && git checkout 1.1.1; fi;'
git submodule foreach 'if [[ "$path" =~ view/extensions/(flarum-pipetables) ]]; then git fetch origin && git checkout v2.0; fi;'

```

