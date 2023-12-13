-- +migrate Up

CREATE TABLE `customers` (
    `uuid` varchar(255) NOT NULL,
    `first_name` varchar(255) NOT NULL,
    `last_name` varchar(255) NOT NULL,
    `email` varchar(300) NOT NULL,
    `phone_number` varchar(100) NOT NULL,
    `registration_date` datetime NOT NULL,
    PRIMARY KEY (`uuid`),
    KEY `customers_uuid_idx` (`uuid`),
    KEY `customers_email_idx` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `orders` (
    `uuid` varchar(255) NOT NULL,
    `customer_uuid` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `placed_date` datetime NOT NULL,
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (customer_uuid) REFERENCES customers(uuid),
    KEY `orders_uuid_idx` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `products` (
    `uuid` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `description` text DEFAULT NULL,
    `image` varchar(255) DEFAULT NULL,
    `availability_status` varchar(255) DEFAULT NULL,
    `available_items` int,
    PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `order_products` (
    `uuid` varchar(255) NOT NULL,
    `order_uuid` varchar(255) NOT NULL,
    `product_uuid` varchar(255) NOT NULL,
    `items` int,
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (order_uuid) REFERENCES orders(uuid),
    FOREIGN KEY (product_uuid) REFERENCES products(uuid),
    KEY `order_products_uuid_idx` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `order_product_reviews` (
    `uuid` varchar(255) NOT NULL,
    `order_product_uuid` varchar(255) NOT NULL,
    `score` int,
    PRIMARY KEY (`uuid`),
    FOREIGN KEY (order_product_uuid) REFERENCES order_products(uuid),
    KEY `order_product_reviews_uuid_idx` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- Pre-populate tables with dummy date
LOCK TABLES `customers` WRITE;
INSERT INTO `customers` (`uuid`, `first_name`, `last_name`, `email`, `phone_number`, `registration_date`)
VALUES
    ('cus1','fnam','lname','email@mail.com','1234567890','2023-12-12 12:12:51'),
    ('cus2','Name1','Name2','mail@mail.gr','0987654321','2022-11-10 12:12:51');
UNLOCK TABLES;

LOCK TABLES `products` WRITE;
INSERT INTO `products` (`uuid`, `name`, `description`, `image`, `availability_status`, `available_items`)
VALUES
    ('prod1','iPhone 23','some descript','An image URL','available',129),
    ('prod2','iSpoon','Desc','img','available',973),
    ('prod3','iTable','A new kind of table','Some_image','available',2);
UNLOCK TABLES;

LOCK TABLES `orders` WRITE;
INSERT INTO `orders` (`uuid`, `customer_uuid`, `status`, `placed_date`)
VALUES
    ('ord1','cus1','placed','2023-12-12 12:12:51'),
    ('ord2','cus2','completed','2022-12-12 12:12:51'),
    ('ord3','cus2','sending','2023-12-12 13:12:51');
UNLOCK TABLES;

LOCK TABLES `order_products` WRITE;
INSERT INTO `order_products` (`uuid`, `order_uuid`, `product_uuid`, `items`)
VALUES
    ('op1','ord1','prod1',2),
    ('op2','ord1','prod2',12),
    ('op3','ord2','prod3',1),
    ('op4','ord3','prod1',1);
UNLOCK TABLES;

-- +migrate Down
DROP TABLE IF EXISTS `order_products`;
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `orders`;
DROP TABLE IF EXISTS `customers`;