CREATE TABLE `products` (
   `id` int AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `name` varchar(255) NOT NULL,
  `image` varchar(255) NOT NULL,
  `category` varchar(255) NOT NULL,
  `description` text,
  `rating` int NOT NULL,
  `num_reviews` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `count_in_stock` int NOT NULL,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `order_items` (
  `id` int PRIMARY KEY NOT NULL,
  `order_id` int NOT NULL,
  `product_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `quantity` int NOT NULL,
  `image` varchar(255) NOT NULL,
  `price` decimal(10,2) NOT NULL
);

CREATE TABLE `orders` (
  `id` int PRIMARY KEY NOT NULL,
  `payment_method` varchar(255) NOT NULL,
  `tax_price` decimal(10,2) NOT NULL,
  `shipping_price` decimal(10,2) NOT NULL,
  `total_price` decimal(10,2) NOT NULL,
  `status` varchar(255) COMMENT 'orderStatus',
  `created_at` datetime,
  `updated_at` datetime
);

ALTER TABLE `order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);