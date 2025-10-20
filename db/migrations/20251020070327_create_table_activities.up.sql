CREATE TABLE activities
(
    id          VARCHAR(100) PRIMARY KEY,
    title       VARCHAR(150) NOT NULL,
    description TEXT         NOT NULL,
    time_info   VARCHAR(100),
    location    VARCHAR(100),
    order_index INT       DEFAULT 1,
    is_active   BOOLEAN   DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);