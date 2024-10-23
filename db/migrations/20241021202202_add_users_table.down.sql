-- Step 1: Drop the foreign key constraint before dropping the table
-- Step 2: Now drop the users table
DROP TABLE IF EXISTS `users`;

-- Optionally, if you want to drop orders as well
DROP TABLE IF EXISTS `orders`;
