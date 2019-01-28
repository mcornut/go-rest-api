create table documents (
    id serial primary key,
    name varchar(60) not null,
    filePath varchar(250) not null,
    thumbPath varchar(250) not null,
    created_at timestamp default current_timestamp
    updated_at timestamp
);