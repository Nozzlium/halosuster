CREATE TYPE "gender" AS ENUM ('male', 'female');

CREATE TABLE IF NOT EXISTS "patients" (
  "identity_number" varchar(16) NOT NULL,
  "user_id" uuid NOT NULL,
  "phone_number" varchar(15) NOT NULL,
  "name" varchar(30) NOT NULL,
  "birthdate" timestamp NOT NULL,
  "gender" GENDER NOT NULL, 
  "identity_card_image_url" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp,
  PRIMARY KEY("identity_number"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
)

