CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name NVARCHAR(255),
    email NVARCHAR(255) UNIQUE,
    password NVARCHAR(255),
    phone NVARCHAR(255),
    role NVARCHAR(50) DEFAULT 'user'
);

CREATE TABLE movies (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    title NVARCHAR(255),
    poster NVARCHAR(MAX),
    duration INT,
    description NVARCHAR(MAX),
    genre NVARCHAR(255)
);

CREATE TABLE studios (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name NVARCHAR(255),
    capacity INT
);

CREATE TABLE seats (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    studio_id INT FOREIGN KEY REFERENCES studios(id),
    seat_number NVARCHAR(50)
);

CREATE TABLE schedules (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    movie_id INT FOREIGN KEY REFERENCES movies(id),
    studio_id INT FOREIGN KEY REFERENCES studios(id),
    start_time DATETIME,
    price FLOAT
);

CREATE TABLE transactions (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    user_id INT FOREIGN KEY REFERENCES users(id),
    total_amount FLOAT,
    status NVARCHAR(50),
    payment_time DATETIME
);

CREATE TABLE tickets (
    id INT IDENTITY(1,1) PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    transaction_id INT FOREIGN KEY REFERENCES transactions(id),
    schedule_id INT FOREIGN KEY REFERENCES schedules(id),
    seat_id INT FOREIGN KEY REFERENCES seats(id)
);