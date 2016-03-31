drop database if exists library;

create database library charset utf8;

use library;

CREATE TABLE `book` (
    `Uid` INT(10) NOT NULL AUTO_INCREMENT,
    `Username` VARCHAR(64) NULL DEFAULT NULL,
    `Bname` VARCHAR(64) NULL DEFAULT NULL,
    PRIMARY KEY (`Uid`)
);

insert into book (Username,Bname)values
("fish1","book1"),
("fish2","book2");

CREATE TABLE `session` (
    `token` VARCHAR(225) NOT NULL,
    `value` VARCHAR(225) NULL DEFAULT NULL,
    PRIMARY KEY (`token`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
