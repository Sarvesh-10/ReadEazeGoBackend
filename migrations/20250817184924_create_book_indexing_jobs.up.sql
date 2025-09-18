-- 1. Create ENUM type for status
CREATE TYPE job_status AS ENUM ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'FAILED');

-- 2. Create the table using the ENUM
CREATE TABLE book_indexing_jobs (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    user_id INT NOT NULL,
    status job_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 3. Create a function to auto-update "updated_at"
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Attach trigger to the table
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON book_indexing_jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
