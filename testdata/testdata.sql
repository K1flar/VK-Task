INSERT INTO actors (full_name, gender, birthday)
VALUES 
    ('Tom Hanks', 'male', '1956-07-09'),
    ('Meryl Streep', 'female', '1949-06-22'),
    ('Leonardo DiCaprio', 'male', '1974-11-11'),
    ('Emma Watson', 'female', '1990-04-15'),
    ('Brad Pitt', 'male', '1963-12-18'),
    ('Jennifer Lawrence', 'female', '1990-08-15'),
    ('Johnny Depp', 'male', '1963-06-09'),
    ('Natalie Portman', 'female', '1981-06-09'),
    ('Robert Downey Jr.', 'male', '1965-04-04'),
    ('Scarlett Johansson', 'female', '1984-11-22'),
    ('Cillian Murphy', 'male', '1976-05-25');

INSERT INTO films (name, description, release_date, rating)
VALUES 
    ('Forrest Gump', 'The story of a man with a low IQ who rose above his challenges.', '1994-07-06', 8),
    ('The Devil Wears Prada', 'A smart but sensible new graduate lands a job as an assistant to Miranda Priestly, the demanding editor-in-chief of a high fashion magazine.', '2006-06-30', 7),
    ('Inception', 'A thief who enters the dreams of others to steal their secrets.', '2010-07-16', 9),
    ('Harry Potter and the Philosopher''s Stone', 'A young boy discovers he is a wizard.', '2001-11-16', 8),
    ('Fight Club', 'An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.', '1999-10-15', 9),
    ('The Hunger Games', 'Katniss Everdeen voluntarily takes her younger sister''s place in the Hunger Games: a televised competition in which two teenagers from each of the twelve Districts of Panem are chosen at random to fight to the death.', '2012-03-23', 7),
    ('Pirates of the Caribbean: The Curse of the Black Pearl', 'Blacksmith Will Turner teams up with eccentric pirate "Captain" Jack Sparrow to save his love, the governor''s daughter, from Jack''s former pirate allies, who are now undead.', '2003-07-09', 8),
    ('Black Swan', 'A committed dancer struggles to maintain her sanity after winning the lead role in a production of Tchaikovsky''s "Swan Lake".', '2010-12-17', 8),
    ('Iron Man', 'After being held captive in an Afghan cave, billionaire engineer Tony Stark creates a unique weaponized suit of armor to fight evil.', '2008-05-02', 8),
    ('The Avengers', 'Earth''s mightiest heroes must come together and learn to fight as a team if they are going to stop the mischievous Loki and his alien army from enslaving humanity.', '2012-05-04', 9),
    ('Oppenheimer', 'A biographical film about J. Robert Oppenheimer, the scientist who headed the Manhattan Project.', '2023-01-01', 10);

INSERT INTO film_actor (film_id, actor_id)
VALUES 
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 4),
    (5, 5),
    (6, 6),
    (7, 7),
    (8, 8),
    (9, 9),
    (10, 10),
    (11, 11);

INSERT INTO users (login, password, role)
VALUES 
    ('admin', '$2a$10$TapsdRWZUU/26uZdj/gpwO4OPf4/0eqxOrlPBZfm3iH74aN5S8l0q', 'admin');
