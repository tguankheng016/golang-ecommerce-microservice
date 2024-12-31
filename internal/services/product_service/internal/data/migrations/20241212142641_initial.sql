-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "first_name" character varying(64) NULL,
    "last_name" character varying(64) NULL,
    "user_name" character varying(256) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NULL,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "public"."users";
-- +goose StatementEnd
