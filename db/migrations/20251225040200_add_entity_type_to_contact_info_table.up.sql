ALTER TABLE contact_info ADD COLUMN entity_type ENUM('pura', 'yayasan', 'pasraman') DEFAULT 'pura' NOT NULL AFTER id;
CREATE INDEX idx_contact_info_entity_type ON contact_info(entity_type);