CREATE TABLE IF NOT EXISTS food (
                                    id INT AUTO_INCREMENT PRIMARY KEY,
                                    name VARCHAR(255) NOT NULL,
                                    description TEXT,
                                    calories DOUBLE,
                                    protein DOUBLE,
                                    fats DOUBLE,
                                    carbohydrates DOUBLE,
                                    created_at DATETIME NOT NULL,
                                    updated_at DATETIME
);

use workout;