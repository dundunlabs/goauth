-- CreateEnum
CREATE TYPE "CredentialMethod" AS ENUM ('PASSWORD');

-- CreateTable
CREATE TABLE "identities" (
    "id" SERIAL NOT NULL,
    "traits" JSONB NOT NULL DEFAULT '{}',

    CONSTRAINT "identities_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "credentials" (
    "id" SERIAL NOT NULL,
    "method" "CredentialMethod" NOT NULL DEFAULT 'PASSWORD',
    "provider" VARCHAR(255),
    "secret" VARCHAR(255) NOT NULL,
    "identity_id" INTEGER NOT NULL,

    CONSTRAINT "credentials_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "sessions" (
    "id" SERIAL NOT NULL,
    "authenticated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expires_at" TIMESTAMP(3),
    "credential_id" INTEGER NOT NULL,

    CONSTRAINT "sessions_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "identities_traits_username_key" ON "identities"(("traits"->>'username'));
CREATE UNIQUE INDEX "identities_traits_email_key" ON "identities"(("traits"->>'email'));
CREATE UNIQUE INDEX "identities_traits_phone_key" ON "identities"(("traits"->>'phone'));

-- CreateIndex
CREATE UNIQUE INDEX "credentials_provider_secret_key" ON "credentials"("provider", "secret");

-- CreateIndex
CREATE UNIQUE INDEX "credentials_identity_id_method_key" ON "credentials"("identity_id", "method") WHERE "method" = 'PASSWORD';

-- AddForeignKey
ALTER TABLE "credentials" ADD CONSTRAINT "credentials_identity_id_fkey" FOREIGN KEY ("identity_id") REFERENCES "identities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "sessions" ADD CONSTRAINT "sessions_credential_id_fkey" FOREIGN KEY ("credential_id") REFERENCES "credentials"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
