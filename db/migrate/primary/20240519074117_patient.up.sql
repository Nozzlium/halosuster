CREATE TYPE "gender" AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS "patients" (
  "identity_number" varchar(16) NOT NULL,
  "phone_number" varchar(15) NOT NULL,
  "name" varchar(30) NOT NULL,
  "birthdate" timestamp NOT NULL,
  "gender" GENDER NOT NULL, 
  "identity_card_image_url" varchar(255) NOT NULL,
  PRIMARY KEY("identity_number")
)

