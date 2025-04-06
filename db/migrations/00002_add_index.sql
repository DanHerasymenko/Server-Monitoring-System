-- +goose Up
CREATE INDEX idx_metrics_server_id ON metrics(server_id);
CREATE INDEX idx_metrics_timestamp ON metrics(timestamp);
CREATE INDEX idx_servers_last_active ON servers(last_active);

-- +goose Down
DROP INDEX IF EXISTS idx_metrics_server_ip;
DROP INDEX IF EXISTS idx_metrics_timestamp;
DROP INDEX IF EXISTS idx_servers_last_active;
