CREATE TABLE IF NOT EXISTS public.users
(
    id              serial PRIMARY KEY,
    phone           varchar(20) NOT NULL,
    first_name      varchar(50) NOT NULL,
    second_name     varchar(50) NOT NULL,
    email           varchar(50) NOT NULL,
    code            smallint,
    "role"          varchar(50) NOT NULL,
    avatar          varchar(250),
    created_date    timestamp NOT NULL,
    updated_date    timestamp NOT NULL,
    deleted_date    timestamp NULL
);
