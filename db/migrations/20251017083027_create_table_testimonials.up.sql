CREATE TABLE testimonials
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    avatar_url  TEXT,
    rating      INT          NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT         NOT NULL,
    is_active   BOOLEAN   DEFAULT TRUE,
    order_index INT       DEFAULT 1,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);