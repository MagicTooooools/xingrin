-- Restore words and lines columns to directory table
ALTER TABLE directory ADD COLUMN IF NOT EXISTS words INTEGER;
ALTER TABLE directory ADD COLUMN IF NOT EXISTS lines INTEGER;
