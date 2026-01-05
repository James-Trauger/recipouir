INSERT INTO users (id, username, first_name, last_name, password) VALUES (47, 'asoiaf', 'ned', 'stark', '$2a$10$SaU8AfqGvUoeLIaJ2W7KY.e3ybJ5RC9mkQxiwFN3tBqu2Jj1vb.XW');

INSERT INTO recipe (id, title, created_by, servings) VALUES (99, 'Cookies', 47, 5);
INSERT INTO step (number, recipe_id, description) VALUES (1, 99, 'mix the ingredients');
INSERT INTO step (number, recipe_id, description) VALUES (2, 99, 'bake at 350');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('sugar', 99, 100, 'grams');
INSERT INTO ingredient (name, recipe_id, amount, unit, denominator) 
    VALUES ('butter', 99, 1, 'cups', 2);


INSERT INTO recipe (id, title, created_by, servings) VALUES (100, 'Cake', 47, 10);
INSERT INTO step (number, recipe_id, description) VALUES (1, 100, 'combine sugar, flour, salt');
INSERT INTO step (number, recipe_id, description) VALUES (2, 100, 'mix wet ingredients');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('flour', 100, 2, 'cups');
INSERT INTO ingredient (name, recipe_id, amount, unit, numerator, denominator) 
    VALUES ('salt', 100, 1, 'tsp', 1, 2);

INSERT INTO recipe (id, title, created_by, servings) VALUES (101, 'Apple Pie', 47, 6);
INSERT INTO step (number, recipe_id, description) VALUES (1, 101, 'combine flour and water');
INSERT INTO step (number, recipe_id, description) VALUES (2, 101, 'cut apples into pieces');
INSERT INTO step (number, recipe_id, description) VALUES (3, 101, 'mix everything and bake');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('flour', 101, 300, 'cgrams');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('water', 101, 150, 'ml');
INSERT INTO ingredient (name, recipe_id, amount, unit, numerator, denominator) 
    VALUES ('apples', 101, 2, 'cups', 2, 3);