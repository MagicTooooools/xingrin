-- Add GIN indexes for PostgreSQL array fields
-- GIN indexes enable efficient queries on array columns (e.g., @>, &&, etc.)

-- Website tech array
CREATE INDEX IF NOT EXISTS idx_website_tech_gin ON website USING GIN (tech);

-- Endpoint arrays
CREATE INDEX IF NOT EXISTS idx_endpoint_tech_gin ON endpoint USING GIN (tech);
CREATE INDEX IF NOT EXISTS idx_endpoint_matched_gf_patterns_gin ON endpoint USING GIN (matched_gf_patterns);

-- Website snapshot tech array
CREATE INDEX IF NOT EXISTS idx_website_snap_tech_gin ON website_snapshot USING GIN (tech);

-- Endpoint snapshot arrays
CREATE INDEX IF NOT EXISTS idx_endpoint_snap_tech_gin ON endpoint_snapshot USING GIN (tech);
CREATE INDEX IF NOT EXISTS idx_endpoint_snap_matched_gf_patterns_gin ON endpoint_snapshot USING GIN (matched_gf_patterns);

-- Scan arrays
CREATE INDEX IF NOT EXISTS idx_scan_engine_ids_gin ON scan USING GIN (engine_ids);
CREATE INDEX IF NOT EXISTS idx_scan_container_ids_gin ON scan USING GIN (container_ids);

-- Scheduled scan arrays
CREATE INDEX IF NOT EXISTS idx_scheduled_scan_engine_ids_gin ON scheduled_scan USING GIN (engine_ids);
