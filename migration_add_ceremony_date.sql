-- Migration: Add ceremony_date and is_upcoming columns to awards table
-- Run this if your database was created before these columns were added

-- Add ceremony_date column (nullable, for future ceremonies)
ALTER TABLE awards ADD COLUMN IF NOT EXISTS ceremony_date DATE;

-- Add is_upcoming column (default false for existing awards)
ALTER TABLE awards ADD COLUMN IF NOT EXISTS is_upcoming BOOLEAN NOT NULL DEFAULT false;

