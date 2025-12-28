ALTER TABLE organization_details
    ADD COLUMN vision_mission_image_url TEXT NULL AFTER work_program,
    ADD COLUMN work_program_image_url TEXT NULL AFTER vision_mission_image_url,
    ADD COLUMN rules_image_url TEXT NULL AFTER work_program_image_url;

UPDATE organization_details
SET work_program_image_url = image_url;

ALTER TABLE organization_details
DROP COLUMN image_url;