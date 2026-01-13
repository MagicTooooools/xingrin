-- Drop words and lines columns from directory table
ALTER TABLE directory DROP COLUMN IF EXISTS words;
ALTER TABLE directory DROP COLUMN IF EXISTS lines;
