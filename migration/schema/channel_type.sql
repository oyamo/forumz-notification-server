create table if not exists channel_type (
    id varchar(2) not null primary key,
    k varchar(64) not null,
    name varchar(128) not null
)