CREATE TABLE IF NOT EXISTS progress (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        client_id INT NOT NULL,
                                        day_id INT NOT NULL,
                                        completed DATETIME NOT NULL,
                                        FOREIGN KEY (client_id) REFERENCES users(id),
                                        FOREIGN KEY (day_id) REFERENCES days(id)
);

use workout;