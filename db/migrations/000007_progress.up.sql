CREATE TABLE IF NOT EXISTS progress (
    id INT AUTO_INCREMENT PRIMARY KEY,
    client_id INT NOT NULL,
    day_id INT NOT NULL,
    food_completed BOOLEAN NOT NULL DEFAULT FALSE,
    exercise_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed DATETIME,
    UNIQUE KEY client_day_unique (client_id, day_id),
    FOREIGN KEY (client_id) REFERENCES users(id),
    FOREIGN KEY (day_id) REFERENCES days(id)
);
USE workout;