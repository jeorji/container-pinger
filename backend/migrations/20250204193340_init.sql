-- +goose Up
-- +goose StatementBegin
CREATE TABLE containers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    image TEXT NOT NULL,
    state TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE pings (
    container_id TEXT NOT NULL,
    ip TEXT NOT NULL,
    status TEXT NOT NULL,
    latency REAL,
    ping_time TIMESTAMP NOT NULL,

    CONSTRAINT pings_pk PRIMARY KEY (ping_time, container_id),
    CONSTRAINT pings_container_fk FOREIGN KEY (container_id)
        REFERENCES containers (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pings;
DROP TABLE containers;
-- +goose StatementEnd
