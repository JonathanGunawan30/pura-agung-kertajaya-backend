CREATE TABLE hero_slides
(
    id          VARCHAR(100) PRIMARY KEY,
    image_url   TEXT NOT NULL,
    order_index INT  NOT NULL DEFAULT 1,
    is_active   BOOLEAN       DEFAULT TRUE,
    created_at  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP     DEFAULT CURRENT_TIMESTAMP
);