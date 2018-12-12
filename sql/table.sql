CREATE TABLE `t_user` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `name` varchar(50) DEFAULT NULL,
      `password` varchar(100) DEFAULT '123',
      `createDate` date DEFAULT NULL,
      PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8




CREATE TABLE `t_user_message` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `userId` int(11) NOT NULL ,
      `messageKey` varchar(1000) DEFAULT NULL,
      `messageValue` varchar(1000) DEFAULT '123',
      `createDate` date DEFAULT NULL,
      PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8
ALTER TABLE `t_user_message` ADD INDEX index_userId ( `userId` )


CREATE TABLE `t_user_quota` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` int(11) NOT NULL ,
    `writeSpeed` int(11) NOT NULL ,
    `readSpeed` int(11)  NOT NULL ,
    `createDate` date DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
ALTER TABLE `t_user_quota` ADD INDEX index_userId ( `userId` )
