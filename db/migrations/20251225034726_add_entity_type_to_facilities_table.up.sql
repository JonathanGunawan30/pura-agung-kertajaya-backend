ALTER TABLE `facilities`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_facilities_entity_type` ON `facilities` (`entity_type`);