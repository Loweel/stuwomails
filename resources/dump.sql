CREATE TABLE `check_log` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `input` VARCHAR(100) NOT NULL COLLATE 'utf8_unicode_ci',
    `password` VARCHAR(100) NOT NULL COLLATE 'utf8_unicode_ci',
    `ip_address` VARCHAR(45) NOT NULL COLLATE 'utf8_unicode_ci',
    `creation_date` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY (`id`)
)
COLLATE='utf8_unicode_ci'
ENGINE=InnoDB;


CREATE TABLE `mailing_list` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    `address` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    `subscribe_address` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    PRIMARY KEY (`id`)
)
COLLATE='utf8_unicode_ci'
ENGINE=InnoDB;


CREATE TABLE `request` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `mailing_list_id` INT(10) UNSIGNED NOT NULL,
    `first_name` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    `last_name` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    `room` VARCHAR(50) NOT NULL COLLATE 'utf8_unicode_ci',
    `email` VARCHAR(100) NOT NULL COLLATE 'utf8_unicode_ci',
    `ip_address` VARCHAR(45) NOT NULL COLLATE 'utf8_unicode_ci',
    `creation_date` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY (`id`),
    INDEX `FK_request_mailing_list` (`mailing_list_id`),
    CONSTRAINT `FK_request_mailing_list` FOREIGN KEY (`mailing_list_id`) REFERENCES `mailing_list` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
)
COLLATE='utf8_unicode_ci'
ENGINE=InnoDB;
