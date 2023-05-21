-- DropForeignKey
ALTER TABLE "credentials" DROP CONSTRAINT "credentials_identity_id_fkey";

-- DropForeignKey
ALTER TABLE "sessions" DROP CONSTRAINT "sessions_credential_id_fkey";

-- DropTable
DROP TABLE "credentials";

-- DropTable
DROP TABLE "identities";

-- DropTable
DROP TABLE "sessions";

-- DropEnum
DROP TYPE "CredentialMethod";
