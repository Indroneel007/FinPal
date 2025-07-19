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

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

ALTER TABLE "locations" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");