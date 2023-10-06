CREATE TABLE services (
    service_id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    details VARCHAR(255) NOT NULL
);

CREATE TABLE metrics (
    metric_id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    metric_type VARCHAR(255) NOT NULL,
    details VARCHAR(255) NOT NULL
);

CREATE TABLE events (
    event_id BIGSERIAL PRIMARY KEY,
    time_stamp TIMESTAMP WITH TIME ZONE NOT NULL UNIQUE,
    service_id BIGINT REFERENCES services ON DELETE CASCADE NOT NULL
);

CREATE TABLE events_with_metrics (
    event_id BIGINT REFERENCES events ON DELETE CASCADE,
    metric_id BIGINT REFERENCES metrics ON DELETE CASCADE,
    metric_value VARCHAR(255) NOT NULL,
    PRIMARY KEY (event_id, metric_id)
);