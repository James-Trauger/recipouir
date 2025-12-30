INSERT INTO users (id, username, first_name, last_name, password) VALUES (47, 'asoiaf', 'ned', 'stark', '$2a$10$SaU8AfqGvUoeLIaJ2W7KY.e3ybJ5RC9mkQxiwFN3tBqu2Jj1vb.XW');
INSERT INTO recipe (id, title, created_by, servings) VALUES (99, 'Cookies', 47, 5);
INSERT INTO step (number, recipe_id, description) VALUES (1, 99, 'mix the ingredients');
INSERT INTO step (number, recipe_id, description) VALUES (2, 99, 'bake at 350');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('sugar', 99, 100, 'grams');
INSERT INTO ingredient (name, recipe_id, amount, unit, denominator) 
    VALUES ('butter', 99, 1, 'cups', 2);