ALTER TABLE `about_section`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_about_section_entity_type` ON `about_section` (`entity_type`);
