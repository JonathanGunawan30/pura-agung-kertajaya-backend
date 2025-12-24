ALTER TABLE `galleries`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_galleries_entity_type` ON `galleries` (`entity_type`);