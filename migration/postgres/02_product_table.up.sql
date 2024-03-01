
CREATE TABLE "product" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR(36) NOT NULL,
    "price" NUMERIC NOT NULL,
    "category_id" UUID REFERENCES "category"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);