-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."categories" (
    "id" serial NOT NULL,
    "name" character varying(256) NULL,
    "normalized_name" character varying(256) NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" bigint NULL,
    "updated_at" timestamptz NULL,
    "updated_by" bigint NULL,
    "is_deleted" boolean NOT NULL,
    "deleted_at" timestamptz NULL,
    "deleted_by" bigint NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "idx_category_is_deleted" ON "public"."categories" ("is_deleted");
CREATE INDEX IF NOT EXISTS "idx_category_name" ON "public"."categories" ("normalized_name" ASC NULLS LAST);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_category_is_deleted";
DROP INDEX IF EXISTS "public"."idx_category_name";
DROP TABLE "public"."categories";
-- +goose StatementEnd
