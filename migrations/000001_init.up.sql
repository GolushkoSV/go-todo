CREATE SCHEMA users;

CREATE TABLE users.user
(
    id           SERIAL PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15)
);


CREATE SCHEMA tasks;

CREATE TABLE tasks.task
(
    id             SERIAL PRIMARY KEY,
    version        BIGINT       NOT NULL DEFAULT 1,
    title          VARCHAR(100) NOT NULL,
    description    VARCHAR(1000),
    completed      BOOLEAN      NOT NULL,
    author_user_id INTEGER      NOT NULL REFERENCES users.user (id),
    created_at     TIMESTAMPTZ  NOT NULL,
    completed_at   TIMESTAMPTZ,

    CHECK (
        (completed = FALSE AND completed_at IS NULL)
            OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
    )
)