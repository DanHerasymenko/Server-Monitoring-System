-- +goose Up
CREATE TABLE IF NOT EXISTS servers (
                                       id SERIAL PRIMARY KEY,
                                       ip_address VARCHAR(255) UNIQUE NOT NULL,
                                       last_active BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS metrics (
                                       id SERIAL PRIMARY KEY,
                                       server_id INTEGER NOT NULL,
                                       cpu_usage DOUBLE PRECISION,
                                       ram_usage DOUBLE PRECISION,
                                       disk_usage DOUBLE PRECISION,
                                       timestamp BIGINT NOT NULL,
                                       FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS servers;
