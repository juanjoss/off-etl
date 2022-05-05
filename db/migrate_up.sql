CREATE TABLE "products" (
  "barcode" text PRIMARY KEY,
  "name" text,
  "quantity" text,
  "image_url" text,
  "energy_100g" float,
  "energy_serving" float,
  "nutrient_levels_id" int NOT NULL,
  "nova_group" int NOT NULL,
  "nutriscore_score" int,
  "nutriscore_grade" text NOT NULL
);

CREATE TABLE "brands" (
  "tag" text PRIMARY KEY
);

CREATE TABLE "product_brands" (
  "barcode" text,
  "tag" text,
  PRIMARY KEY ("barcode", "tag")
);

CREATE TABLE "nutrient_levels" (
  "id" SERIAL PRIMARY KEY,
  "fat" text NOT NULL,
  "saturated_fat" text NOT NULL,
  "sugar" text NOT NULL,
  "salt" text NOT NULL
);

CREATE TABLE "nova_groups" (
  "group" int PRIMARY KEY,
  "description" text
);

CREATE TABLE "nutriscore" (
  "grade" text PRIMARY KEY,
  "description" text
);

ALTER TABLE "products" ADD FOREIGN KEY ("nutrient_levels_id") REFERENCES "nutrient_levels" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("nova_group") REFERENCES "nova_groups" ("group");

ALTER TABLE "products" ADD FOREIGN KEY ("nutriscore_grade") REFERENCES "nutriscore" ("grade");

ALTER TABLE "product_brands" ADD FOREIGN KEY ("barcode") REFERENCES "products" ("barcode");

ALTER TABLE "product_brands" ADD FOREIGN KEY ("tag") REFERENCES "brands" ("tag");

INSERT INTO "nova_groups" ("group", "description") VALUES (1, 'Unprocessed or minimally processed foods');
INSERT INTO "nova_groups" ("group", "description") VALUES (2, 'Processed culinary ingredients');
INSERT INTO "nova_groups" ("group", "description") VALUES (3, 'Processed foods');
INSERT INTO "nova_groups" ("group", "description") VALUES (4, 'Ultra-processed food and drink products');

INSERT INTO "nutriscore" ("grade", "description") VALUES ('A', '-15 to -1 points for solid foods, water for beverages');
INSERT INTO "nutriscore" ("grade", "description") VALUES ('B', '0 to 2 points, <= 1 for beverages');
INSERT INTO "nutriscore" ("grade", "description") VALUES ('C', '3 to 10 points, 2 to 5 points for beverages');
INSERT INTO "nutriscore" ("grade", "description") VALUES ('D', '11 to 18 points, 6 to 9 points for beverages');
INSERT INTO "nutriscore" ("grade", "description") VALUES ('E', '19 to 40 points for solid foods, 10 to 40 for beverages');