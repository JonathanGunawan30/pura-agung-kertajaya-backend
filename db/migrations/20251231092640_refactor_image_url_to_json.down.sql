ALTER TABLE articles ADD COLUMN image_url VARCHAR(255);
UPDATE articles SET image_url = JSON_UNQUOTE(JSON_EXTRACT(images, '$.lg')) WHERE images IS NOT NULL;
ALTER TABLE articles DROP COLUMN images;

ALTER TABLE hero_slides ADD COLUMN image_url VARCHAR(255);
UPDATE hero_slides SET image_url = JSON_UNQUOTE(JSON_EXTRACT(images, '$.lg')) WHERE images IS NOT NULL;
ALTER TABLE hero_slides DROP COLUMN images;

ALTER TABLE about_section ADD COLUMN image_url VARCHAR(255);
UPDATE about_section SET image_url = JSON_UNQUOTE(JSON_EXTRACT(images, '$.lg')) WHERE images IS NOT NULL;
ALTER TABLE about_section DROP COLUMN images;

ALTER TABLE facilities ADD COLUMN image_url VARCHAR(255);
UPDATE facilities SET image_url = JSON_UNQUOTE(JSON_EXTRACT(images, '$.lg')) WHERE images IS NOT NULL;
ALTER TABLE facilities DROP COLUMN images;

ALTER TABLE galleries ADD COLUMN image_url VARCHAR(255);
UPDATE galleries SET image_url = JSON_UNQUOTE(JSON_EXTRACT(images, '$.lg')) WHERE images IS NOT NULL;
ALTER TABLE galleries DROP COLUMN images;