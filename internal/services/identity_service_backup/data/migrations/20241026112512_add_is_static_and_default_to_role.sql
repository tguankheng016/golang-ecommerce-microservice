-- +goose Up
-- modify "roles" table
ALTER TABLE "public"."roles" ADD COLUMN "is_static" boolean NOT NULL DEFAULT false, ADD COLUMN "is_default" boolean NOT NULL DEFAULT false;

-- +goose Down
-- reverse: modify "roles" table
ALTER TABLE "public"."roles" DROP COLUMN "is_default", DROP COLUMN "is_static";
