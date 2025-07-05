CREATE TABLE IF NOT EXISTS verification_codes
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    code       VARCHAR(255) NOT NULL,
    created_at DATETIME     NOT NULL
);

use workout;