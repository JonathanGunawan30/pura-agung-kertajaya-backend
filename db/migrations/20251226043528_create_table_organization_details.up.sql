CREATE TABLE `organization_details`
(
    `id`            VARCHAR(100) PRIMARY KEY,
    `entity_type`  ENUM('pura', 'yayasan', 'pasraman') NOT NULL,

    `vision`       LONGTEXT,
    `mission`      LONGTEXT,
    `rules`        LONGTEXT,
    `work_program` LONGTEXT,
    `image_url`    TEXT,

    `created_at`   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY `unique_entity_detail` (`entity_type`)
) ENGINE=InnoDB;