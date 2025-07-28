CREATE TABLE IF NOT EXISTS exercises (
                                         id INT AUTO_INCREMENT PRIMARY KEY,
                                         name VARCHAR(255) NOT NULL,
                                         description TEXT,
                                         media_url TEXT,
                                         sets VARCHAR(50),
                                         repetitions VARCHAR(50),
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

use workout;