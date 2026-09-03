-- migrate:up
ALTER TABLE divisions
  ADD COLUMN division_type VARCHAR(20) NOT NULL DEFAULT 'internal';

-- migrate:down
ALTER TABLE divisions
  DROP COLUMN division_type;
