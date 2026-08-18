INSERT INTO users (username, password_hash, display_name, created_at, updated_at)
VALUES (
    'admin',
    '$2a$10$b5fIC0PzsNX9FN4DkrAe/uR6frQ/DqTSU23dvPp6QFTCmRSSmZb.a',
    '系统管理员',
    NOW(),
    NOW()
)
ON CONFLICT (username) DO NOTHING;
