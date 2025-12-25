ALTER TABLE `organization_members`
    ADD COLUMN `entity_type` ENUM('pura', 'yayasan', 'pasraman') NOT NULL DEFAULT 'pura' AFTER `id`;

CREATE INDEX `idx_org_members_entity_type` ON `organization_members` (`entity_type`);