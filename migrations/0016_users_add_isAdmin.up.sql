-- No-op migration: isAdmin column already exists in initial users schema.
-- This ensures idempotent migration application without version-specific ALTER TABLE syntax.
SELECT 1;