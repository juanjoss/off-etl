CREATE TABLE "products" (
  "barcode" text PRIMARY KEY,
  "name" text,
  "quantity" text,
  "image_url" text
);

CREATE TABLE "brands" (
  "tag" text PRIMARY KEY
);

CREATE TABLE "product_brands" (
  "barcode" text,
  "tag" text
);

ALTER TABLE "product_brands" ADD FOREIGN KEY ("barcode") REFERENCES "products" ("barcode");

ALTER TABLE "product_brands" ADD FOREIGN KEY ("tag") REFERENCES "brands" ("tag");
