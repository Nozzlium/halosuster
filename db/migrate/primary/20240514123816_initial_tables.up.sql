CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid NOT NULL,
  "employee_id" bigint not null,
  "name" varchar(50),
  -- "username" varchar(255) NOT NULL,
  "password" varchar(100),
  -- "email" varchar(255) NULL DEFAULT NULL,
  -- "email_verified_at" timestamp NULL DEFAULT NULL,
  "identity_card_image_url" varchar(255),
  "created_by" uuid NULL DEFAULT NULL, 
  "updated_by" uuid NULL DEFAULT NULL,
  "deleted_by" uuid NULL DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp NULL DEFAULT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("employee_id")
);
