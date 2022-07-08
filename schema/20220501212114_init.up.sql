CREATE TABLE "public"."foo"
(
    "id"         text        NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "name"       text        NOT NULL,
    PRIMARY KEY ("id")
);
