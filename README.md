# Forum Z Notification Server

## Download

```bash
git clone git@github.com:oyamo/forumz-notification-server.git
cd forumz-notification-server
```

## Build
```shell
docker build -t notification-server:1.0.0 .
```

## Liquibase migration
```shell
docker run --rm --network local-sandbox \
    --volume `pwd`/migration:/liquibase/changelog liquibase/liquibase:4.13 \
    --url="jdbc:postgresql://postgres:5432/notification" \
    --changeLogFile=master-changeLog.yml \
    --username=dev \
    --password='Test@12345' \
    --database-changelog-table-name=database_changelog \
    --database-changelog-lock-table-name=database_changelog_lock \
    --driver=org.postgresql.Driver \
    --log-level=info \
    update
```

## Run
```shell
docker run -d \
  --network sandbox \
  -e NOTIFICATIONS_SERVICE_DATABASE_DSN='postgresql://postgres:5432/notification?user=dev&password=Test@12345' \
  -e NOTIFICATIONS_SERVICE_REDIS_SERVER='redis:6379' \
  -e NOTIFICATIONS_SERVICE_KAFKA_CONSUMER='' \
  -e NOTIFICATIONS_SERVICE_KAFKA_PRODUCER='' \
  -e NOTIFICATIONS_SERVICE_SENDER_NAME='Forum Z' \
  -e NOTIFICATIONS_SERVICE_SENDER_EMAIL='user@example.com' \
  -e NOTIFICATIONS_SERVICE_SENDGRID_API_KEY='' \
  -e NOTIFICATIONS_SERVICE_MAILTRAP_API_KEY='' \
  notification-server:1.0.0

```