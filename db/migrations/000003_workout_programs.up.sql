CREATE TABLE IF NOT EXISTS workout_programs (
                                                id INT AUTO_INCREMENT PRIMARY KEY,
                                                trainer_id INT NOT NULL,
                                                name VARCHAR(255) NOT NULL,
                                                days INT NOT NULL,
                                                description TEXT,
                                                created_at DATETIME NOT NULL,
                                                updated_at DATETIME,
                                                FOREIGN KEY (trainer_id) REFERENCES users(id)
);

use workout;