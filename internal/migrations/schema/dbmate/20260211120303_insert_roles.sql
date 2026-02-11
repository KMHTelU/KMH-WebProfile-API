-- migrate:up
INSERT INTO roles (id, name, description) VALUES
 	(gen_random_uuid(), 'SuperUser', 'Administrator role with full permissions'),
 	(gen_random_uuid(), 'Guest', 'Standard user role with limited permissions');

-- migrate:down
DELETE FROM roles WHERE name IN ('SuperUser', 'Guest');
