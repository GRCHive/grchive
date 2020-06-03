ALTER TABLE database_refresh
ADD COLUMN refresh_has_diff BOOLEAN NOT NULL DEFAULT false;
