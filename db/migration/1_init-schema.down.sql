-- 1. Drop foreign key constraints
ALTER TABLE "notifications" DROP CONSTRAINT "notifications_account_id_fkey";
ALTER TABLE "entries" DROP CONSTRAINT "entries_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_from_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT "transfers_to_account_id_fkey";
ALTER TABLE "accounts" DROP CONSTRAINT "accounts_group_id_fkey";

-- 2. Drop indexes
DROP INDEX IF EXISTS "accounts_owner_idx";
DROP INDEX IF EXISTS "accounts_group_id_owner_idx";
DROP INDEX IF EXISTS "groups_currency_idx";
DROP INDEX IF EXISTS "groups_type_idx";
DROP INDEX IF EXISTS "notifications_account_id_idx";
DROP INDEX IF EXISTS "entries_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_idx";
DROP INDEX IF EXISTS "transfers_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_to_account_id_idx";

-- 3. Remove column comments
COMMENT ON COLUMN "entries"."amount" IS NULL;
COMMENT ON COLUMN "transfers"."amount" IS NULL;

-- 4. Drop tables in dependency-safe order
DROP TABLE IF EXISTS "notifications";
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "accounts";
DROP TABLE IF EXISTS "groups";