DROP TABLE IF EXISTS User;

CREATE TABLE User (
       id bigint(20) not null auto_increment,
       name varchar(255) default null,
       primary key (id)
) engine=innodb charset=utf8;

INSERT INTO User values (null, 'liangzuobin');
