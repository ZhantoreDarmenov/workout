CREATE TABLE IF NOT EXISTS program_invites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    program_id INT NOT NULL,
    email VARCHAR(255) NOT NULL,
    message TEXT,
    access_days INT NOT NULL,
    token VARCHAR(64) NOT NULL UNIQUE,
    client_id INT,
    accepted_at DATETIME,
    access_expires DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    FOREIGN KEY (program_id) REFERENCES workout_programs(id),
    FOREIGN KEY (client_id) REFERENCES users(id)
);

USE workout;
