CREATE TABLE users(
    id bigserial not null primary key,
    email varchar not null unique,
    encypted_password varchar not null
);