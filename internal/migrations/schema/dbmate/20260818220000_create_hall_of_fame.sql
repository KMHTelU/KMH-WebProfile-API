-- migrate:up
CREATE TABLE IF NOT EXISTS hof_generations (
  id          UUID PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  year_start  INT NOT NULL,
  year_end    INT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  milestones  JSONB NOT NULL DEFAULT '[]',
  accent      VARCHAR(20) NOT NULL DEFAULT '',
  sort_order  INT NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hof_people (
  id             UUID PRIMARY KEY,
  generation_id  UUID NOT NULL REFERENCES hof_generations(id) ON DELETE CASCADE,
  name           VARCHAR(150) NOT NULL,
  role           VARCHAR(150) NOT NULL DEFAULT '',
  study_program  VARCHAR(150) NOT NULL DEFAULT '',
  biography      TEXT NOT NULL DEFAULT '',
  contributions  TEXT NOT NULL DEFAULT '',
  legacy         TEXT NOT NULL DEFAULT '',
  quote          TEXT NOT NULL DEFAULT '',
  fields         JSONB NOT NULL DEFAULT '[]',
  photo_media_id UUID REFERENCES media(id),
  sort_order     INT NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ DEFAULT NOW(),
  updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hof_achievements (
  id           UUID PRIMARY KEY,
  person_id    UUID NOT NULL REFERENCES hof_people(id) ON DELETE CASCADE,
  title        VARCHAR(250) NOT NULL,
  category     VARCHAR(50) NOT NULL DEFAULT 'Other',
  year         INT NOT NULL,
  organization VARCHAR(200) NOT NULL DEFAULT '',
  result       VARCHAR(150) NOT NULL DEFAULT '',
  description  TEXT NOT NULL DEFAULT '',
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hof_timeline_events (
  id          UUID PRIMARY KEY,
  year_label  VARCHAR(30) NOT NULL,
  title       VARCHAR(200) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  sort_order  INT NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- migrate:down
DROP TABLE IF EXISTS hof_achievements;
DROP TABLE IF EXISTS hof_people;
DROP TABLE IF EXISTS hof_timeline_events;
DROP TABLE IF EXISTS hof_generations;
