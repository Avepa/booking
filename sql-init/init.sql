CREATE DATABASE `booking`;
USE `booking`;

CREATE TABLE `room` (
  `id` 					INT NOT NULL AUTO_INCREMENT,
  `description` 		VARCHAR(1024) NOT NULL,
  `price` 				FLOAT NOT NULL,
  `date` 				DATE NOT NULL,
  
  PRIMARY KEY (`id`),
  INDEX `SERCH` (`price` ASC, `date` ASC) INVISIBLE
 );
  
  CREATE TABLE `bookings` (
  `id` 					INT NOT NULL AUTO_INCREMENT,
  `room_id` 			INT NOT NULL,
  `date_start` 			DATE NOT NULL,
  `date_end` 			DATE NOT NULL,
  
  PRIMARY KEY (`id`),
  INDEX `SERCH` (`date_start` ASC, `room_id` ASC) INVISIBLE,
  FOREIGN KEY (`room_id`)   REFERENCES `room` (`id`) ON DELETE CASCADE
);