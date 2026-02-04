CREATE TABLE IF NOT EXISTS users (
  id                UUID PRIMARY KEY,
  name              VARCHAR(100),
  email             VARCHAR(150) UNIQUE,
  password_hash     TEXT,
  role_id           UUID NOT NULL,
  is_active         BOOLEAN DEFAULT true,
  last_login_at     TIMESTAMPTZ,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS roles (
  id          UUID PRIMARY KEY,
  name        VARCHAR(50),
  description TEXT
);

CREATE TABLE IF NOT EXISTS organization_profile (
  id                UUID PRIMARY KEY,
  name              VARCHAR(150),
  short_name        VARCHAR(50), -- KMH
  description       TEXT,
  vision            TEXT,
  mission           TEXT,
  history           TEXT,
  logo_media_id     UUID,
  address           TEXT,
  email             VARCHAR(150),
  phone             VARCHAR(50),
  instagram_url     TEXT,
  youtube_url       TEXT,
  website_url       TEXT,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (logo_media_id) REFERENCES media(id)
);

CREATE TABLE IF NOT EXISTS divisions (
  id                UUID PRIMARY KEY,
  name              VARCHAR(100),
  slug              VARCHAR(120) UNIQUE,
  description       TEXT,
  icon_media_id     UUID,
  coordinator_id    UUID,
  is_active         BOOLEAN DEFAULT true,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (icon_media_id) REFERENCES media(id),
  FOREIGN KEY (coordinator_id) REFERENCES members(id)
);

CREATE TABLE IF NOT EXISTS members (
  id              UUID PRIMARY KEY,
  name            VARCHAR(120),
  npm             VARCHAR(30),
  photo_media_id  UUID,
  bio             TEXT,
  email           VARCHAR(150),
  phone           VARCHAR(50),
  instagram_url   TEXT,
  period_start    YEAR,
  period_end      YEAR,
  is_active       BOOLEAN DEFAULT true,
  created_at      TIMESTAMPTZ DEFAULT NOW(),
  updated_at      TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (photo_media_id) REFERENCES media(id)
);

CREATE TABLE IF NOT EXISTS member_divisions (
  id            UUID PRIMARY KEY,
  member_id     UUID,
  division_id   UUID,
  role_title    VARCHAR(100), -- Ketua, Wakil, Staff, dll
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (member_id) REFERENCES members(id),
  FOREIGN KEY (division_id) REFERENCES divisions(id)
);

CREATE TABLE IF NOT EXISTS events (
  id                UUID PRIMARY KEY,
  title             VARCHAR(200),
  slug              VARCHAR(220) UNIQUE,
  description       TEXT,
  event_type        VARCHAR(30), -- internal | external
  start_time        TIMESTAMPTZ,
  end_time          TIMESTAMPTZ,
  location          VARCHAR(200),
  google_maps_url   TEXT,
  registration_url  TEXT,
  cover_media_id    UUID,
  status            VARCHAR(30), -- upcoming | ongoing | finished
  is_published      BOOLEAN DEFAULT false,
  created_by        UUID,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (cover_media_id) REFERENCES media(id),
  FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS blog_categories (
  id          UUID PRIMARY KEY,
  name        VARCHAR(100),
  slug        VARCHAR(120) UNIQUE
);

CREATE TABLE IF NOT EXISTS blog_posts (
  id                UUID PRIMARY KEY,
  title             VARCHAR(200),
  slug              VARCHAR(220) UNIQUE,
  excerpt           TEXT,
  content           TEXT, -- HTML / Markdown
  category_id       UUID,
  featured_media_id UUID,
  author_id         UUID,
  status            VARCHAR(20), -- draft | published
  published_at      TIMESTAMPTZ,
  created_at        TIMESTAMPTZ DEFAULT NOW(),
  updated_at        TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (category_id) REFERENCES blog_categories(id),
  FOREIGN KEY (featured_media_id) REFERENCES media(id),
  FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS blog_tags (
  id    UUID PRIMARY KEY,
  name  VARCHAR(50),
  slug  VARCHAR(60)
);

CREATE TABLE IF NOT EXISTS blog_post_tags (
  post_id UUID,
  tag_id  UUID,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES blog_posts(id),
  FOREIGN KEY (tag_id) REFERENCES blog_tags(id)
);

CREATE TABLE IF NOT EXISTS media (
  id            UUID PRIMARY KEY,
  file_name     TEXT,
  file_type     VARCHAR(30), -- image | video | document
  mime_type     VARCHAR(100),
  file_size     BIGINT,
  url           TEXT,
  alt_text      TEXT,
  caption       TEXT,
  uploaded_by   UUID,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (uploaded_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS galleries (
  id            UUID PRIMARY KEY,
  title         VARCHAR(150),
  description   TEXT,
  event_id      UUID,
  is_public     BOOLEAN DEFAULT true,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE IF NOT EXISTS gallery_items (
  id            UUID PRIMARY KEY,
  gallery_id    UUID,
  media_id      UUID,
  sort_order    INT,
  FOREIGN KEY (gallery_id) REFERENCES galleries(id),
  FOREIGN KEY (media_id) REFERENCES media(id)
);

CREATE TABLE IF NOT EXISTS homepage_banners (
  id            UUID PRIMARY KEY,
  title         VARCHAR(150),
  subtitle      TEXT,
  media_id      UUID,
  cta_text      VARCHAR(50),
  cta_url       TEXT,
  is_active     BOOLEAN DEFAULT true,
  start_date    TIMESTAMPTZ,
  end_date      TIMESTAMPTZ,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (media_id) REFERENCES media(id)
);

CREATE TABLE IF NOT EXISTS contact_messages (
  id          UUID PRIMARY KEY,
  name        VARCHAR(120),
  email       VARCHAR(150),
  subject     VARCHAR(150),
  message     TEXT,
  is_read     BOOLEAN DEFAULT false,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS activity_logs (
  id          UUID PRIMARY KEY,
  user_id     UUID,
  action      VARCHAR(100),
  entity      VARCHAR(100),
  entity_id   UUID,
  ip_address  VARCHAR(45),
  user_agent  TEXT,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
