-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create index on role for filtering
CREATE INDEX idx_users_role ON users(role);

-- Create index on is_active for filtering
CREATE INDEX idx_users_is_active ON users(is_active);

-- Create index on created_at for sorting
CREATE INDEX idx_users_created_at ON users(created_at);

-- Add constraint to ensure role is valid
ALTER TABLE users ADD CONSTRAINT check_valid_role 
    CHECK (role IN ('admin', 'user', 'guest'));

-- Add constraint to ensure email is not empty
ALTER TABLE users ADD CONSTRAINT check_email_not_empty 
    CHECK (email != '');

-- Add constraint to ensure name is not empty
ALTER TABLE users ADD CONSTRAINT check_name_not_empty 
    CHECK (name != '');

-- Add constraint to ensure password is not empty
ALTER TABLE users ADD CONSTRAINT check_password_not_empty 
    CHECK (password != ''); 