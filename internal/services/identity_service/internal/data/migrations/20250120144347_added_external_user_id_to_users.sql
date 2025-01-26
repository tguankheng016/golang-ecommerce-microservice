-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."users" ADD COLUMN "external_user_id" TEXT NULL;
CREATE INDEX IF NOT EXISTS "idx_user_external_user_id" ON "public"."users" ("external_user_id" ASC NULLS LAST);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_user_external_user_id";
ALTER TABLE "public"."users" DROP COLUMN "external_user_id";
-- +goose StatementEnd
