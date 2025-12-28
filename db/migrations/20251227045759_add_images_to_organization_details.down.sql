ALTER TABLE organization_details
    ADD COLUMN image_url TEXT NULL AFTER work_program;

UPDATE organization_details
SET image_url = work_program_image_url;

ALTER TABLE organization_details
DROP COLUMN vision_mission_image_url,
DROP COLUMN work_program_image_url,
DROP COLUMN rules_image_url;