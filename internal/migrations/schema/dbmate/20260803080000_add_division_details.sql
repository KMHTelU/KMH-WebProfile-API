-- migrate:up
ALTER TABLE divisions
  ADD COLUMN subtitle TEXT,
  ADD COLUMN responsibilities JSONB NOT NULL DEFAULT '[]',
  ADD COLUMN programs JSONB NOT NULL DEFAULT '[]';

-- migrate:down
ALTER TABLE divisions
  DROP COLUMN subtitle,
  DROP COLUMN responsibilities,
  DROP COLUMN programs;
