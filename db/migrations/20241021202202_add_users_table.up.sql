-- Create the users table
CREATE TABLE `users` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `is_admin` boolean NOT NULL DEFAULT false
);

-- Alter the orders table to add a user_id column and a foreign key constraint
ALTER TABLE `orders` 
    ADD COLUMN `user_id` int NOT NULL,
    ADD CONSTRAINT `user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
