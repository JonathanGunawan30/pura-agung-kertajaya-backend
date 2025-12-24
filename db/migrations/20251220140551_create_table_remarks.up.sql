CREATE TABLE remarks
(
    id          VARCHAR(100) PRIMARY KEY,
    entity_type ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura',
    name        VARCHAR(100) NOT NULL,
    position    VARCHAR(100) NOT NULL,
    image_url   TEXT,
    content     TEXT         NOT NULL,
    is_active   BOOLEAN   DEFAULT TRUE,
    order_index INT       DEFAULT 1,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);