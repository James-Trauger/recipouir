INSERT INTO users (id, username, first_name, last_name) VALUES (47, 'asoiaf', 'ned', 'stark');
INSERT INTO recipe (id, title, created_by) VALUES (99, 'Cookies', 47);
INSERT INTO step (number, recipe_id, description) VALUES (1, 99, 'mix the ingredients');
INSERT INTO step (number, recipe_id, description) VALUES (2, 99, 'bake at 350');
INSERT INTO ingredient (name, recipe_id, amount, unit) VALUES ('sugar', 99, 100, 'grams');
INSERT INTO ingredient (name, recipe_id, amount, unit, denominator) 
    VALUES ('butter', 99, 1, 'cups', 2);