DROP SCHEMA IF EXISTS `vald` ;
CREATE SCHEMA IF NOT EXISTS `vald` DEFAULT CHARACTER SET utf8 ;

USE `vald` ;

CREATE TABLE IF NOT EXISTS `vald`.`meta_vector` (
    `uuid` VARCHAR(255) NOT NULL,
    `vector` TEXT NOT NULL,
    `meta` VARCHAR(1024) NOT NULL,
    PRIMARY KEY (`uuid`),
    UNIQUE INDEX `meta_unique` (`meta` ASC)
)
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `vald`.`pod_ip` (
    `uuid` VARCHAR(255) NOT NULL,
    `ip` VARCHAR(64) NOT NULL,
    PRIMARY KEY (`uuid`),
    INDEX `ip_index` (`ip` ASC)
)
ENGINE = InnoDB;
