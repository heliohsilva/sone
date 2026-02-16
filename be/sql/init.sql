CREATE USER IF NOT EXISTS 'golang'@'%' IDENTIFIED BY 'golang';
CREATE DATABASE IF NOT EXISTS devbook;

GRANT ALL PRIVILEGES ON devbook.* TO 'golang'@'%';

USE devbook;

DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS followers;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(50) not null,
    nickname varchar(50) not null unique,
    email varchar(50) not null unique,
    password char(64) not null,
    createdAt timestamp default current_timestamp()
) ENGINE=INNODB;


CREATE TABLE followers(
    user_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    follower_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    primary key(user_id, follower_id)

) ENGINE=INNODB;


CREATE TABLE posts(
    id int auto_increment primary key,
    title varchar(50) not null,
    content varchar(300) not null,

    user_id int not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    likes int default 0,
    createdAt timestamp default current_timestamp
)ENGINE=INNODB;