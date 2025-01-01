-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."products" (
    "id" serial NOT NULL,
    "name" character varying(50) NULL,
    "normalized_name" character varying(50) NULL,
    "description" character varying(150) NULL,
    "normalized_description" character varying(150) NULL,
    "price" decimal(10, 2) NOT NULL,
    "stock_quantity" int NOT NULL,
    "category_id" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" bigint NULL,
    "updated_at" timestamptz NULL,
    "updated_by" bigint NULL,
    "is_deleted" boolean NOT NULL,
    "deleted_at" timestamptz NULL,
    "deleted_by" bigint NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_products_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS "idx_product_is_deleted" ON "public"."products" ("is_deleted");
CREATE INDEX IF NOT EXISTS "idx_product_name" ON "public"."products" ("normalized_name" ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS "idx_product_category_id" ON "public"."products" ("category_id" ASC NULLS LAST);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_product_is_deleted";
DROP INDEX IF EXISTS "public"."idx_product_name";
DROP INDEX IF EXISTS "public"."idx_product_category_id";
DROP TABLE "public"."products";
-- +goose StatementEnd
