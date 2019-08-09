create table "app_user"
(
    id varchar(36) not null,
    name varchar(55) not null,
    password_hash varchar(255) not null,
    created timestamp default now() not null,
    updated timestamp default now() not null
);

create unique index user_id_uindex
    on app_user (id);

create unique index user_name_uindex
    on app_user (lower(name));

alter table app_user
    add constraint user_pk
        primary key (id);
