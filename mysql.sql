create database library charset utf8;

CREATE TABLE `book` (
    `Uid` INT(10) NOT NULL AUTO_INCREMENT,
    `Username` VARCHAR(64) NULL DEFAULT NULL,
    `Bname` VARCHAR(64) NULL DEFAULT NULL,
    PRIMARY KEY (`Uid`)
);
