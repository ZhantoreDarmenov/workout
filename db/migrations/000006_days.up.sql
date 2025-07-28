CREATE TABLE IF NOT EXISTS days
(
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    work_out_program_id INT      NOT NULL,
    day_number          INT      NOT NULL,
    exercises_id        INT      NOT NULL,
    food_id             INT      NOT NULL,
    note                TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (work_out_program_id) REFERENCES workout_programs (id),
    FOREIGN KEY (exercises_id) REFERENCES exercises (id),
    FOREIGN KEY (food_id) REFERENCES food (id)
);

use workout;