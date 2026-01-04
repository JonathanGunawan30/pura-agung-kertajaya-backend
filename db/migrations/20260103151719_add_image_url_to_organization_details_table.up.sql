ALTER TABLE organization_details
    ADD COLUMN structure_image_url TEXT DEFAULT NULL;

ALTER TABLE organization_members
DROP COLUMN images;