-- migrate:up
ALTER TABLE events
  ADD COLUMN division_id UUID REFERENCES divisions(id);

-- migrate:down
ALTER TABLE events
  DROP COLUMN division_id;
