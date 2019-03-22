create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  firstname  varchar(64) not null,
  email      varchar(40) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null   
);