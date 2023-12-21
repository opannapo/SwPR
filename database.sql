CREATE TABLE users
(
    id         serial4      NOT NULL,
    full_name  varchar(50)  NOT NULL,
    "password" varchar(100) NOT NULL,
    phone      varchar(20)  NOT NULL,
    created_at timestamptz NULL DEFAULT now(),
    updated_at timestamptz NULL,
    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT user_un UNIQUE (phone)
);