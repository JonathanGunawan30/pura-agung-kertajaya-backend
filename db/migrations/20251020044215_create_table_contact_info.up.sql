CREATE TABLE contact_info
(
    id             VARCHAR(100) PRIMARY KEY,
    address        TEXT NOT NULL,
    phone          VARCHAR(50),
    email          VARCHAR(100),
    visiting_hours VARCHAR(100),
    map_embed_url  TEXT,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);