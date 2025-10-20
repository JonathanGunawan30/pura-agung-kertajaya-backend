CREATE TABLE about_section
(
    id          VARCHAR(100) PRIMARY KEY,
    title       VARCHAR(150) NOT NULL,
    description TEXT         NOT NULL,
    image_url   TEXT,
    is_active   BOOLEAN   DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE about_values
(
    id          VARCHAR(100) PRIMARY KEY,
    about_id    VARCHAR(100),
    title       VARCHAR(100) NOT NULL,
    value       VARCHAR(100) NOT NULL,
    order_index INT       DEFAULT 1,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_about_values_section
        FOREIGN KEY (about_id) REFERENCES about_section (id)
            ON DELETE CASCADE
);