CREATE TYPE users_role AS enum ('Admin', 'Worker', 'HeadOfDepartment');

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role users_role NOT NULL DEFAULT 'Worker',
    is_confirmed BOOLEAN DEFAULT FALSE,
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);
