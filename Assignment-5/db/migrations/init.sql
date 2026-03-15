-- Создание таблицы users с автоинкрементным id и поддержкой soft delete
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    gender      VARCHAR(50),
    birth_date  DATE,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id <> friend_id)
);

INSERT INTO users (id, name, email, gender, birth_date) VALUES
    (1, 'Alice', 'alice@example.com', 'female', '1990-01-01'),
    (2, 'Bob', 'bob@example.com', 'male', '1985-05-10'),
    (3, 'Charlie', 'charlie@example.com', 'male', '1992-07-15'),
    (4, 'Diana', 'diana@example.com', 'female', '1988-12-20'),
    (5, 'Eve', 'eve@example.com', 'female', '1995-03-03'),
    (6, 'Frank', 'frank@example.com', 'male', '1980-11-11'),
    (7, 'Grace', 'grace@example.com', 'female', '1993-09-09'),
    (8, 'Henry', 'henry@example.com', 'male', '1987-04-04'),
    (9, 'Ivy', 'ivy@example.com', 'female', '1991-06-06'),
    (10, 'Jack', 'jack@example.com', 'male', '1984-02-02'),
    (11, 'Karen', 'karen@example.com', 'female', '1996-08-08'),
    (12, 'Leo', 'leo@example.com', 'male', '1983-10-10'),
    (13, 'Mona', 'mona@example.com', 'female', '1994-12-12'),
    (14, 'Nick', 'nick@example.com', 'male', '1982-01-01'),
    (15, 'Olivia', 'olivia@example.com', 'female', '1997-07-07'),
    (16, 'Paul', 'paul@example.com', 'male', '1981-03-03'),
    (17, 'Quinn', 'quinn@example.com', 'female', '1998-05-05'),
    (18, 'Rachel', 'rachel@example.com', 'female', '1986-09-09'),
    (19, 'Sam', 'sam@example.com', 'male', '1999-11-11'),
    (20, 'Tina', 'tina@example.com', 'female', '1989-02-02');

INSERT INTO user_friends (user_id, friend_id) VALUES
    (1, 3), (1, 4), (1, 5), (1, 6),
    (2, 3), (2, 4), (2, 5), (2, 7);

INSERT INTO user_friends (user_id, friend_id) VALUES
    (3, 8), (4, 9), (5, 10), (6, 11), (7, 12);

UPDATE users SET deleted_at = NOW() WHERE id = 9;
