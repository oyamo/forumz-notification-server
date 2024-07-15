create table if not exists notification (
    id uuid not null primary key,
    email_address  varchar(64),
    content varchar(1024) not null ,
    priority smallint default 3,
    retry_count int default 0,
    sent boolean default false,
    datetime_created timestamp not null default current_timestamp,
    last_modified timestamp with time zone  not null default current_timestamp
)