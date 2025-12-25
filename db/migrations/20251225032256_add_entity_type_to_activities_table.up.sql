ALTER TABLE `activities`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_activities_entity_type` ON `activities` (`entity_type`);