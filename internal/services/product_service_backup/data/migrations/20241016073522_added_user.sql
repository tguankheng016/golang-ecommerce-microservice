-- +goose Up
-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "first_name" character varying(64) NULL,
  "last_name" character varying(64) NULL,
  "user_name" character varying(256) NOT NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "users" table
DROP TABLE "public"."users";
