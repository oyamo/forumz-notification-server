#!/usr/bin/bash

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


NOTIFICATIONS_SERVICE_DATABASE_DSN='postgresql://localhost:5432/notification?user=dev&password=Test@12345' \
NOTIFICATIONS_SERVICE_REDIS_SERVER='localhost:6379' \
NOTIFICATIONS_SERVICE_KAFKA_CONSUMER='' \
NOTIFICATIONS_SERVICE_KAFKA_PRODUCER='' \
NOTIFICATIONS_SERVICE_SENDER_NAME='Forum Z' \
NOTIFICATIONS_SERVICE_SENDER_EMAIL='oyamo.xyz@gmail.com' \
NOTIFICATIONS_SERVICE_SENDGRID_API_KEY='' \
NOTIFICATIONS_SERVICE_MAILTRAP_API_KEY='' \
go run ./cmd/server
