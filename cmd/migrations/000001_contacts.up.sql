CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    name       TEXT,
    email      TEXT,
    phone      TEXT,
    external_id INTEGER NOT NULL
);

CREATE INDEX idx_id_contacts ON contacts (id);
CREATE INDEX idx_timestamp_contacts ON contacts (timestamp);
