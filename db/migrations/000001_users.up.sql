CREATE TABLE IF NOT EXISTS users
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    phone         VARCHAR(50),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    role          VARCHAR(50) DEFAULT 'client',
    refresh_token VARCHAR(255),
    expires_at    DATETIME,
    created_at    DATETIME     NOT NULL,
    updated_at    DATETIME
);

use workout;