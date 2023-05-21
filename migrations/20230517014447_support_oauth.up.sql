-- AlterEnum
ALTER TYPE "CredentialMethod" ADD VALUE 'OAUTH';

-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CreateTable
CREATE TABLE "oauth_states" (
    "id" uuid DEFAULT uuid_generate_v4 (),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expires_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "oauth_states_pkey" PRIMARY KEY ("id")
);
