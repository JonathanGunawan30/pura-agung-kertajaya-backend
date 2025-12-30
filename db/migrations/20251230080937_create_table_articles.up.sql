CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    category_id VARCHAR(100) NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    author_name VARCHAR(100) NOT NULL,
    author_role VARCHAR(100) NULL,
    excerpt TEXT,
    content LONGTEXT,
    image_url VARCHAR(255),
    status ENUM('DRAFT', 'PUBLISHED', 'ARCHIVED') DEFAULT 'DRAFT',
    is_featured BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_articles_category
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE = InnoDB;

CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_featured ON articles(is_featured);
CREATE INDEX idx_articles_published_at ON articles(published_at);