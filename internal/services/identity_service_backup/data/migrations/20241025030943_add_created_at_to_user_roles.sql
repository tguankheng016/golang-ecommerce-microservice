-- +goose Up
-- modify "user_roles" table
ALTER TABLE "public"."user_roles" ADD COLUMN "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP, ADD COLUMN "created_by" bigint NULL;

-- +goose Down
-- reverse: modify "user_roles" table
ALTER TABLE "public"."user_roles" DROP COLUMN "created_by", DROP COLUMN "created_at";
