ALTER TABLE activities
    ADD COLUMN event_date DATETIME NULL;

CREATE INDEX idx_activities_event_date ON activities(event_date);