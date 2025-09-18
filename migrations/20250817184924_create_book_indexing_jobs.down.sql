-- Drop trigger and function
DROP TRIGGER IF EXISTS set_updated_at ON book_indexing_jobs;
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the table
DROP TABLE IF EXISTS book_indexing_jobs;

-- Drop the ENUM type
DROP TYPE IF EXISTS job_status;
