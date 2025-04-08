-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    image TEXT DEFAULT 'uploads/no-image.png',
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    status INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    image TEXT DEFAULT 'uploads/no-image.png'
);

CREATE TABLE IF NOT EXISTS driver_car (
    driver_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    car_id INTEGER REFERENCES cars(id) ON DELETE CASCADE,
    PRIMARY KEY (driver_id, car_id)
);

CREATE TABLE IF NOT EXISTS trips (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    start_time TIMESTAMP DEFAULT NOW(),
    end_time TIMESTAMP,
    status TEXT DEFAULT 'active',
    rating INTEGER,
    duration INTEGER DEFAULT 0,       -- ⏱ продолжительность в секундах
    distance REAL DEFAULT 0           -- 📍 расстояние в метрах
);

CREATE TABLE IF NOT EXISTS location_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    trip_id INTEGER REFERENCES trips(id),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    recorded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS status_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    old_status INTEGER,
    new_status INTEGER,
    changed_at TIMESTAMP DEFAULT NOW()
);

-- 🔍 Индексы
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_location_trip_id ON location_logs(trip_id);
CREATE INDEX IF NOT EXISTS idx_location_user_id ON location_logs(user_id);

-- +goose Down

DROP INDEX IF EXISTS idx_location_user_id;
DROP INDEX IF EXISTS idx_location_trip_id;
DROP INDEX IF EXISTS idx_users_status;

DROP TABLE IF EXISTS status_logs;
DROP TABLE IF EXISTS location_logs;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS driver_car;
DROP TABLE IF EXISTS cars;
DROP TABLE IF EXISTS users;
