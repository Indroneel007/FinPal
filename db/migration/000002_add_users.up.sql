CREATE TABLE "users" (
  "username" varchar PRIMARY KEY NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "salary" bigint NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "locations" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" text NOT NULL,
  "address" text NOT NULL DEFAULT 'Bangalore, India',
  "latitude" double precision DEFAULT 12.97,
  "longitude" double precision DEFAULT 77.59
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

DROP INDEX IF EXISTS "accounts_owner_currency_idx";

CREATE UNIQUE INDEX IF NOT EXISTS "accounts_owner_currency_type_idx" ON "accounts" ("owner", "currency", "type");

ALTER TABLE "locations" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");