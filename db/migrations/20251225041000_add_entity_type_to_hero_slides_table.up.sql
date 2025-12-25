ALTER TABLE `hero_slides`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_hero_slides_entity_type` ON `hero_slides` (`entity_type`);
