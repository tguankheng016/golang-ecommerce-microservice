-- +goose Up
-- modify "users" table
-- modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "external_user_id" uuid;

-- Populate the new column with default values
UPDATE "public"."users" SET "external_user_id" = gen_random_uuid();

-- Alter the column to add the NOT NULL constraint
ALTER TABLE "public"."users" ALTER COLUMN "external_user_id" SET NOT NULL;

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "external_user_id";
