drop database if exists library;

create database library charset utf8;

use library;

create table t_client(
	clientId integer not null auto_increment,
	username varchar(256) not null,
	password varchar(256) not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
	primary key( clientId )
)engine=innodb default charset=utf8mb4 auto_increment = 10001;
alter table t_client add unique index nameIndex(username);


create table t_book(
	bid integer not null auto_increment,
	bname varchar(256) not null,
	username varchar(256) not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
	primary key( bid )
)engine=innodb default charset=utf8mb4 auto_increment = 10001;