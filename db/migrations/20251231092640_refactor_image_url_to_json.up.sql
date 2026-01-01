ALTER TABLE articles ADD COLUMN images JSON;
UPDATE articles
SET images = JSON_OBJECT(
        'blur', image_url, 'avatar', image_url, 'xs', image_url, 'sm', image_url,
        'md', image_url, 'lg', image_url, 'xl', image_url, '2xl', image_url, 'fhd', image_url
             ) WHERE image_url IS NOT NULL AND image_url != '';
ALTER TABLE articles DROP COLUMN image_url;

ALTER TABLE hero_slides ADD COLUMN images JSON;
UPDATE hero_slides
SET images = JSON_OBJECT(
        'blur', image_url, 'avatar', image_url, 'xs', image_url, 'sm', image_url,
        'md', image_url, 'lg', image_url, 'xl', image_url, '2xl', image_url, 'fhd', image_url
             ) WHERE image_url IS NOT NULL AND image_url != '';
ALTER TABLE hero_slides DROP COLUMN image_url;

ALTER TABLE about_section ADD COLUMN images JSON;
UPDATE about_section
SET images = JSON_OBJECT(
        'blur', image_url, 'avatar', image_url, 'xs', image_url, 'sm', image_url,
        'md', image_url, 'lg', image_url, 'xl', image_url, '2xl', image_url, 'fhd', image_url
             ) WHERE image_url IS NOT NULL AND image_url != '';
ALTER TABLE about_section DROP COLUMN image_url;

ALTER TABLE facilities ADD COLUMN images JSON;
UPDATE facilities
SET images = JSON_OBJECT(
        'blur', image_url, 'avatar', image_url, 'xs', image_url, 'sm', image_url,
        'md', image_url, 'lg', image_url, 'xl', image_url, '2xl', image_url, 'fhd', image_url
             ) WHERE image_url IS NOT NULL AND image_url != '';
ALTER TABLE facilities DROP COLUMN image_url;

ALTER TABLE galleries ADD COLUMN images JSON;
UPDATE galleries
SET images = JSON_OBJECT(
        'blur', image_url, 'avatar', image_url, 'xs', image_url, 'sm', image_url,
        'md', image_url, 'lg', image_url, 'xl', image_url, '2xl', image_url, 'fhd', image_url
             ) WHERE image_url IS NOT NULL AND image_url != '';
ALTER TABLE galleries DROP COLUMN image_url;