create table documents (
    id serial primary key,
    name varchar(60) not null,
    file_path varchar(250) not null,
    thumb_path varchar(250) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp
);