ALTER TABLE organization_members
    ADD COLUMN images JSON DEFAULT NULL;

ALTER TABLE organization_details
DROP COLUMN structure_image_url;