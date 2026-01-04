DROP INDEX idx_activities_event_date ON activities;

ALTER TABLE activities
DROP COLUMN event_date;