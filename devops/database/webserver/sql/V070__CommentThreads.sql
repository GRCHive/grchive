CREATE TABLE IF NOT EXISTS comment_threads (
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE
);

ALTER TABLE comments
ADD COLUMN IF NOT EXISTS thread_id BIGINT;

ALTER TABLE comments
DROP CONSTRAINT IF EXISTS comments_thread_id_fkey;

ALTER TABLE comments
ADD CONSTRAINT comments_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES comment_threads(id) ON DELETE CASCADE;
