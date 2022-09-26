CREATE TABLE `asset`
(
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `asset` varchar(255) NOT NULL,
    `tokenid` varchar(255) NOT NULL,
    `image` varchar(255) NOT NULL,
    `thumbnail` varchar(255),
    `timestamp` bigint(20) NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;