-- Remove GIN indexes

DROP INDEX IF EXISTS idx_scheduled_scan_engine_ids_gin;
DROP INDEX IF EXISTS idx_scan_container_ids_gin;
DROP INDEX IF EXISTS idx_scan_engine_ids_gin;
DROP INDEX IF EXISTS idx_endpoint_snap_matched_gf_patterns_gin;
DROP INDEX IF EXISTS idx_endpoint_snap_tech_gin;
DROP INDEX IF EXISTS idx_website_snap_tech_gin;
DROP INDEX IF EXISTS idx_endpoint_matched_gf_patterns_gin;
DROP INDEX IF EXISTS idx_endpoint_tech_gin;
DROP INDEX IF EXISTS idx_website_tech_gin;
