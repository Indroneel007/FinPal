CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "type" varchar NOT NULL,
  "group_id" bigint,
  "has_accepted" bool,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "groups" (
  "id" BIGSERIAL PRIMARY KEY,
  "group_name" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "notifications" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "message" text NOT NULL,
  "sent" bool NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE UNIQUE INDEX ON "accounts" ("group_id", "owner");

CREATE INDEX ON "groups" ("currency");

CREATE INDEX ON "groups" ("type");

CREATE INDEX ON "notifications" ("account_id");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'It can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "notifications" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "transfers" ADD COLUMN "group_id" bigint REFERENCES groups(id);