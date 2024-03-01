
CREATE TABLE "market" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "address" VARCHAR,
    "phone_number" VARCHAR,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "market_product_relation" (
    "product_id" UUID NOT NULL REFERENCES "product" ("id"),
    "market_id" UUID NOT NULL REFERENCES "market" ("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX market_product_relation_idx ON "market_product_relation"("product_id", "market_id");
