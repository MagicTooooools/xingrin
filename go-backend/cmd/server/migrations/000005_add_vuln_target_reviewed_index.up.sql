-- Add composite index for target_id + is_reviewed queries
-- Optimizes: COUNT pending vulnerabilities by target, filter by target + review status
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_vuln_target_reviewed ON vulnerability(target_id, is_reviewed);
