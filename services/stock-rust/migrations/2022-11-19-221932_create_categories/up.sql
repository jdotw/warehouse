-- Your SQL goes here
CREATE TABLE categories (
  id uuid DEFAULT public.uuid_generate_v4() NOT NULL PRIMARY KEY,
  name text NOT NULL
)