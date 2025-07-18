-- 1. Drop foreign key constraints if tables exist
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'accounts') THEN
    ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
    DROP INDEX IF EXISTS "accounts_owner_currency_idx";
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'locations') THEN
    ALTER TABLE "locations" DROP CONSTRAINT IF EXISTS "locations_username_fkey";
  END IF;
END $$;

-- 2. Drop tables (dependents first)
DROP TABLE IF EXISTS "locations";
DROP TABLE IF EXISTS "users";
