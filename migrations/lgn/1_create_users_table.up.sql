create table "users"
(
    id varchar(36) not null,
    name varchar(55) not null,
    password_hash varchar(255) not null,
    created timestamp default now() not null,
    updated timestamp default now() not null
);

create unique index user_id_uindex
    on users (id);

create unique index user_name_uindex
    on users (lower(name));

alter table users
    add constraint user_pk
        primary key (id);
