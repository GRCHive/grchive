CREATE UNIQUE INDEX ON database_refresh (db_id)
WHERE refresh_finish_time IS NULL;
