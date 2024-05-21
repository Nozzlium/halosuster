CREATE TABLE IF NOT EXISTS "records" (
  "id" uuid NOT NULL,
  "identity_number" varchar(16) NOT NULL,
  "user_id" uuid NOT NULL,
  "symptomps" varchar(2000) NOT NULL,
  "medications" varchar(2000) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "deleted_at" timestamp,
  PRIMARY KEY("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("identity_number") REFERENCES "patients" ("identity_number") ON DELETE CASCADE
)
