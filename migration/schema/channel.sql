create table if not exists channel (
    id UUID NOT NULL primary key,
    channel_type_id varchar(2) not null,
    name varchar(128),
    datetime_created timestamp not null default current_timestamp,
    last_modified timestamp with time zone  not null default current_timestamp,
    constraint fk_channel_type_channel foreign key(channel_type_id) references channel_type(id)
)