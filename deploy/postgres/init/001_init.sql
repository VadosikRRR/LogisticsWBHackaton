CREATE TABLE IF NOT EXISTS raw_records (
    id BIGSERIAL PRIMARY KEY,
    route_id BIGINT NOT NULL,
    office_from_id BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    status_1 BIGINT NOT NULL,
    status_2 BIGINT NOT NULL,
    status_3 BIGINT NOT NULL,
    status_4 BIGINT NOT NULL,
    status_5 BIGINT NOT NULL,
    status_6 BIGINT NOT NULL,
    status_7 BIGINT NOT NULL,
    status_8 BIGINT NOT NULL,
    target_2h DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_raw_records_timestamp ON raw_records (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_raw_records_route_ts ON raw_records (route_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_raw_records_office_ts ON raw_records (office_from_id, timestamp DESC);

CREATE TABLE IF NOT EXISTS transport_requests (
    id TEXT PRIMARY KEY,
    request_key TEXT NOT NULL,
    route_id BIGINT NOT NULL,
    office_from_id BIGINT NOT NULL,
    slot_timestamp TIMESTAMPTZ NOT NULL,
    required_vehicles INTEGER NOT NULL,
    predicted_volume DOUBLE PRECISION NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_transport_requests_status ON transport_requests (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transport_requests_slot ON transport_requests (slot_timestamp DESC);
CREATE UNIQUE INDEX IF NOT EXISTS ux_transport_requests_active_key
ON transport_requests (request_key)
WHERE status IN ('CREATED','SENT','CONFIRMED','IN_PROGRESS');
