-- +migrate Up

CREATE TABLE `customers` (
    `uuid` varchar(255) NOT NULL,
    `first_name` varchar(255) DEFAULT NULL,
    `last_name` varchar(255) DEFAULT NULL,
    `email` varchar(300) DEFAULT NULL,
    `phone_number` varchar(100) DEFAULT NULL,
    `registration_date` datetime DEFAULT NULL,
    PRIMARY KEY (`uuid`),
    KEY `customers_uuid_idx` (`uuid`),
    KEY `customers_email_idx` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `orders` (
    `uuid` varchar(255) NOT NULL,
    `customer_uuid` varchar(255) DEFAULT NULL,
    `status` varchar(255) DEFAULT NULL,
    `placed_date` datetime DEFAULT NULL,
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (customer_uuid) REFERENCES customers(uuid),
    KEY `orders_uuid_idx` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `order_products` (
    `uuid` varchar(255) NOT NULL,
    `order_uuid` varchar(255) DEFAULT NULL,
    `items` int,
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (order_uuid) REFERENCES orders(uuid),
    KEY `order_products_uuid_idx` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- +migrate Down
DROP TABLE IF EXISTS `order_products`;
DROP TABLE IF EXISTS `orders`;
DROP TABLE IF EXISTS `customers`;