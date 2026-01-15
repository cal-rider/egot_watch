-- EGOT Tracker Database Setup
-- Run this in your Supabase SQL Editor

-- Create enum type for award types
CREATE TYPE award_type AS ENUM ('Emmy', 'Grammy', 'Oscar', 'Tony');

-- Create celebrities table
CREATE TABLE IF NOT EXISTS celebrities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL,
    photo_url TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_celebrities_name ON celebrities (LOWER(name));

-- Create awards table
CREATE TABLE IF NOT EXISTS awards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    celebrity_id UUID NOT NULL REFERENCES celebrities(id) ON DELETE CASCADE,
    type award_type NOT NULL,
    year INTEGER NOT NULL,
    work TEXT NOT NULL,
    category TEXT NOT NULL,
    is_winner BOOLEAN NOT NULL DEFAULT false
);

-- Create index on celebrity_id for faster award lookups
CREATE INDEX IF NOT EXISTS idx_awards_celebrity_id ON awards (celebrity_id);

-- Create index on type for filtering by award type
CREATE INDEX IF NOT EXISTS idx_awards_type ON awards (type);

-- Insert sample data for testing
INSERT INTO celebrities (name, slug, photo_url) VALUES
    ('Viola Davis', 'viola-davis', NULL),
    ('John Legend', 'john-legend', NULL),
    ('Audrey Hepburn', 'audrey-hepburn', NULL)
ON CONFLICT (name) DO NOTHING;

-- Insert sample awards for Viola Davis (EGOT winner)
INSERT INTO awards (celebrity_id, type, year, work, category, is_winner)
SELECT c.id, a.type, a.year, a.work, a.category, a.is_winner
FROM celebrities c
CROSS JOIN (VALUES
    ('Emmy'::award_type, 2015, 'How to Get Away with Murder', 'Outstanding Lead Actress in a Drama Series', true),
    ('Emmy'::award_type, 2023, 'The First Lady', 'Outstanding Lead Actress in a Limited or Anthology Series', false),
    ('Grammy'::award_type, 2023, 'Finding Me', 'Best Audio Book, Narration & Storytelling Recording', true),
    ('Oscar'::award_type, 2017, 'Fences', 'Best Supporting Actress', true),
    ('Oscar'::award_type, 2012, 'The Help', 'Best Actress', false),
    ('Tony'::award_type, 2001, 'King Hedley II', 'Best Actress in a Play', true),
    ('Tony'::award_type, 2010, 'Fences', 'Best Actress in a Play', true)
) AS a(type, year, work, category, is_winner)
WHERE c.name = 'Viola Davis'
ON CONFLICT DO NOTHING;

-- Insert sample awards for John Legend (EGOT winner)
INSERT INTO awards (celebrity_id, type, year, work, category, is_winner)
SELECT c.id, a.type, a.year, a.work, a.category, a.is_winner
FROM celebrities c
CROSS JOIN (VALUES
    ('Emmy'::award_type, 2018, 'Jesus Christ Superstar Live in Concert', 'Outstanding Variety Special (Live)', true),
    ('Grammy'::award_type, 2006, 'Get Lifted', 'Best New Artist', true),
    ('Grammy'::award_type, 2011, 'Shine', 'Best R&B Song', true),
    ('Oscar'::award_type, 2015, 'Selma', 'Best Original Song - Glory', true),
    ('Tony'::award_type, 2017, 'Jitney', 'Best Revival of a Play', true)
) AS a(type, year, work, category, is_winner)
WHERE c.name = 'John Legend'
ON CONFLICT DO NOTHING;

-- Insert sample awards for Audrey Hepburn (EGOT winner - posthumous Grammy)
INSERT INTO awards (celebrity_id, type, year, work, category, is_winner)
SELECT c.id, a.type, a.year, a.work, a.category, a.is_winner
FROM celebrities c
CROSS JOIN (VALUES
    ('Emmy'::award_type, 1993, 'Gardens of the World with Audrey Hepburn', 'Outstanding Individual Achievement - Informational Programming', true),
    ('Grammy'::award_type, 1994, 'Audrey Hepburn''s Enchanted Tales', 'Best Spoken Word Album for Children', true),
    ('Oscar'::award_type, 1954, 'Roman Holiday', 'Best Actress', true),
    ('Tony'::award_type, 1954, 'Ondine', 'Best Actress in a Play', true)
) AS a(type, year, work, category, is_winner)
WHERE c.name = 'Audrey Hepburn'
ON CONFLICT DO NOTHING;
