create database if not exists dbname;

create table if not exists dbname.user (firstname varchar(60) not null) engine = innodb;

insert into dbname.user (firstname)
values
    ('Marcin')
;