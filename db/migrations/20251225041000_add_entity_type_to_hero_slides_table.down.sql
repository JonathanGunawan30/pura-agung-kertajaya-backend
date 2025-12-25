DROP INDEX `idx_hero_slides_entity_type` ON `hero_slides`;

ALTER TABLE `hero_slides`
    DROP COLUMN `entity_type`;
