CREATE TABLE site_identity
(
    id                    VARCHAR(100) PRIMARY KEY,
    site_name             VARCHAR(150) NOT NULL,
    logo_url              TEXT,
    tagline               VARCHAR(255),
    primary_button_text   VARCHAR(50),
    primary_button_link   VARCHAR(255),
    secondary_button_text VARCHAR(50),
    secondary_button_link VARCHAR(255),
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);