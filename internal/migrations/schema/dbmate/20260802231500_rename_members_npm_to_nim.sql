-- migrate:up
ALTER TABLE members RENAME COLUMN npm TO nim;

-- migrate:down
ALTER TABLE members RENAME COLUMN nim TO npm;
