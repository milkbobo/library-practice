drop database if exists library;

create database library charset utf8;

use library;

CREATE TABLE `book` (
    `Bid` INT(10) NOT NULL AUTO_INCREMENT,
    `Username` VARCHAR(64) NULL DEFAULT NULL,
    `Bname` VARCHAR(64) NULL DEFAULT NULL,
    PRIMARY KEY (`Bid`)
);

insert into book (Username,Bname)values
("fish1","book1"),
("fish2","book2");

CREATE TABLE `session` (
    `token` VARCHAR(225) NOT NULL,
    `value` VARCHAR(225) NULL DEFAULT NULL,
    PRIMARY KEY (`token`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `userinfo` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(225) NOT NULL,
    `password` VARCHAR(225) NOT NULL,
    index username (username),
    PRIMARY KEY (`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
