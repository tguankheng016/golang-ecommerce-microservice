-- +goose Up
-- create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "name" character varying(64) NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "created_by" bigint NULL,
  "updated_by" bigint NULL,
  "deleted_by" bigint NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "categories" table
DROP TABLE "public"."categories";
