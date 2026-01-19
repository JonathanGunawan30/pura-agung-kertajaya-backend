ALTER TABLE users
    ADD COLUMN role ENUM('pura','yayasan','pasraman','super')
NOT NULL
DEFAULT 'pura';
