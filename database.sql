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

CREATE TABLE login_attempt
(
    id         serial  NOT NULL,
    user_id    int8    NOT NULL,
    status     boolean NOT NULL,
    created_at timestamptz NULL DEFAULT now(),
    CONSTRAINT login_attempt_pk PRIMARY KEY (id),
    CONSTRAINT login_attempt_fk FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX login_attempt_user_id_idx ON login_attempt (user_id, status);
