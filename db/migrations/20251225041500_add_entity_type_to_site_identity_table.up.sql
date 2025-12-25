ALTER TABLE `site_identity`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_site_identity_entity_type` ON `site_identity` (`entity_type`);
