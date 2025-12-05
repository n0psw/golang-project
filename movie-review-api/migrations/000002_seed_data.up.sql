INSERT INTO genres (name) VALUES 
('Action'),
('Comedy'),
('Drama'),
('Horror'),
('Sci-Fi'),
('Romance'),
('Thriller'),
('Adventure'),
('Fantasy'),
('Animation');

INSERT INTO users (email, username, password_hash, role) VALUES 
('admin@example.com', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin');

INSERT INTO movies (title, description, release_year, director, duration_minutes) VALUES 
('The Matrix', 'A computer hacker learns about the true nature of reality', 1999, 'Lana Wachowski', 136),
('Inception', 'A thief who steals corporate secrets through dream-sharing technology', 2010, 'Christopher Nolan', 148),
('Pulp Fiction', 'The lives of two mob hitmen, a boxer, a gangster and his wife intertwine', 1994, 'Quentin Tarantino', 154);

INSERT INTO movie_genres (movie_id, genre_id) VALUES 
((SELECT id FROM movies WHERE title = 'The Matrix'), (SELECT id FROM genres WHERE name = 'Action')),
((SELECT id FROM movies WHERE title = 'The Matrix'), (SELECT id FROM genres WHERE name = 'Sci-Fi')),
((SELECT id FROM movies WHERE title = 'Inception'), (SELECT id FROM genres WHERE name = 'Action')),
((SELECT id FROM movies WHERE title = 'Inception'), (SELECT id FROM genres WHERE name = 'Sci-Fi')),
((SELECT id FROM movies WHERE title = 'Pulp Fiction'), (SELECT id FROM genres WHERE name = 'Drama')),
((SELECT id FROM movies WHERE title = 'Pulp Fiction'), (SELECT id FROM genres WHERE name = 'Thriller'));

