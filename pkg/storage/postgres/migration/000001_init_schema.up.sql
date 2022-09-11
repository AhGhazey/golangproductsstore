CREATE TABLE "product" (
  "id" bigserial PRIMARY KEY,
  "sku" varchar(30) NOT NULL,
  "country" varchar(5) NOT NULL,
  "name" varchar(100) NOT NULL,
  "quantity" int NOT NULL,
  "created at" varchar,
  CONSTRAINT sku_country_unique UNIQUE (sku, country)
);

CREATE INDEX ON "product" ("sku");

CREATE INDEX ON "product" ("country");

CREATE INDEX ON "product" ("sku", "country");
