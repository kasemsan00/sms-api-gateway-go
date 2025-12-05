-- --------------------------------------------------------
-- Host:                         171.103.89.169
-- Server version:               9.4.0 - MySQL Community Server - GPL
-- Server OS:                    Linux
-- HeidiSQL Version:             12.13.0.7147
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for conference
CREATE DATABASE IF NOT EXISTS `conference` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `conference`;

-- Dumping structure for table conference.car_track
CREATE TABLE IF NOT EXISTS `car_track` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'open',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `room` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latitude` decimal(12,9) DEFAULT NULL,
  `longitude` decimal(12,9) DEFAULT NULL,
  `accuracy` decimal(20,6) DEFAULT NULL,
  `speed` int DEFAULT NULL,
  `heading` int DEFAULT NULL,
  `altitude` float DEFAULT NULL,
  `altitudeAccuracy` float DEFAULT NULL,
  `dtmUpdated` datetime DEFAULT (now()),
  `dtmCreated` datetime DEFAULT (now()),
  `dtmStarted` datetime DEFAULT NULL,
  `dtmArrived` datetime DEFAULT NULL,
  `dtmCanceled` datetime DEFAULT NULL,
  `dtmCompleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `room` (`room`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.case_data
CREATE TABLE IF NOT EXISTS `case_data` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `caseId` int NOT NULL,
  `service` int DEFAULT NULL,
  `roomId` int DEFAULT NULL,
  `operationNumber` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `hn` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `patientMobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mobileCreated` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `caseType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `dtmCreated` datetime DEFAULT NULL,
  `organization` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=378 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.chat_message
CREATE TABLE IF NOT EXISTS `chat_message` (
  `id` int NOT NULL AUTO_INCREMENT,
  `room` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `identity` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chat_identity` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userName` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `color` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `files` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `replyToMessageId` int DEFAULT NULL,
  `replyToUserName` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `replyToText` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `dtmCreated` datetime DEFAULT NULL,
  `userType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=232 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.color_scheme
CREATE TABLE IF NOT EXISTS `color_scheme` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `color_hex` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.data_log
CREATE TABLE IF NOT EXISTS `data_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `dtmCreated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27156 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.files
CREATE TABLE IF NOT EXISTS `files` (
  `id` int NOT NULL AUTO_INCREMENT,
  `linkId` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `elementId` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `filename` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `thumbnail` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `fileType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `size` bigint NOT NULL,
  `mimetype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `width` int DEFAULT NULL,
  `height` int DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT (now()),
  `updatedAt` timestamp NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=140 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.link_connect
CREATE TABLE IF NOT EXISTS `link_connect` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `roomUserId` int DEFAULT NULL,
  `sms` int DEFAULT '1',
  `recordId` int DEFAULT NULL,
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `linkID` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `domainIndex` int DEFAULT '0',
  `share` int DEFAULT '0',
  `enabled` int DEFAULT '1',
  `userName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `room` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userType` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `linkType` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `crmSender` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `accuracy` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latitude` decimal(12,9) DEFAULT NULL,
  `longitude` decimal(12,9) DEFAULT NULL,
  `patientLatitude` decimal(12,9) unsigned DEFAULT NULL,
  `patientLongitude` decimal(12,9) unsigned DEFAULT NULL,
  `patientUpdated` datetime DEFAULT NULL,
  `service` int DEFAULT NULL,
  `errorVideo` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `errorLocation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `os` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `userAgent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `requireJoinPermission` int DEFAULT '0',
  `requireUserName` int DEFAULT '0',
  `requirePassword` int DEFAULT '0',
  `oneTimeLink` int DEFAULT '0',
  `password` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `isAdmin` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '0',
  `dtmConnection` datetime DEFAULT NULL,
  `dtmDisconnect` datetime DEFAULT NULL,
  `dtmCreated` datetime DEFAULT NULL,
  `dtmExpired` datetime DEFAULT NULL,
  `sync_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_link_connect_expired_enabled` (`dtmExpired`,`enabled`),
  KEY `idx_link_connect_room` (`room`),
  KEY `idx_link_connect_linkID` (`linkID`),
  KEY `FK_link_connect_room_user` (`roomUserId`),
  CONSTRAINT `FK_link_connect_room_user` FOREIGN KEY (`roomUserId`) REFERENCES `room_user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=111528 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.node_livekit
CREATE TABLE IF NOT EXISTS `node_livekit` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nodeName` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `livekitHost` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `livekitLocal` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `livekitApiKey` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `livekitApiSecret` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `lastHealthCheck` datetime DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.notification
CREATE TABLE IF NOT EXISTS `notification` (
  `notificationId` int NOT NULL AUTO_INCREMENT,
  `userName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'à¸œà¸¹à¹‰à¹à¸ˆà¹‰à¸‡',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `message` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'à¸‚à¹‰à¸­à¸„à¸§à¸²à¸¡à¹€à¸žà¸´à¹ˆà¸¡à¹€à¸•à¸´à¸¡',
  `caseId` int DEFAULT NULL COMMENT 'à¸«à¸¡à¸²à¸¢à¹€à¸¥à¸‚à¹€à¸„à¸ª',
  `read` tinyint NOT NULL DEFAULT '0' COMMENT 'à¸ªà¸–à¸²à¸™à¸°à¸à¸²à¸£à¸­à¹ˆà¸²à¸™: 0 = à¸¢à¸±à¸‡à¹„à¸¡à¹ˆà¹„à¸”à¹‰à¸­à¹ˆà¸²à¸™, 1 = à¸­à¹ˆà¸²à¸™à¹à¸¥à¹‰à¸§',
  `notificationType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'à¸›à¸£à¸°à¹€à¸ à¸—à¸‚à¸­à¸‡à¸à¸²à¸£à¹à¸ˆà¹‰à¸‡à¹€à¸•à¸·à¸­à¸™',
  `relatedUrl` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'URL à¸«à¸£à¸·à¸­ Path à¸—à¸µà¹ˆà¹€à¸à¸µà¹ˆà¸¢à¸§à¸‚à¹‰à¸­à¸‡à¸à¸±à¸šà¸à¸²à¸£à¹à¸ˆà¹‰à¸‡à¹€à¸•à¸·à¸­à¸™ (à¸ªà¸³à¸«à¸£à¸±à¸šà¸„à¸¥à¸´à¸à¹„à¸›à¸«à¸™à¹‰à¸²à¸£à¸²à¸¢à¸¥à¸°à¹€à¸­à¸µà¸¢à¸”)',
  `dtmRead` datetime DEFAULT NULL COMMENT 'à¸§à¸±à¸™à¹€à¸§à¸¥à¸²à¸—à¸µà¹ˆà¸­à¹ˆà¸²à¸™à¸à¸²à¸£à¹à¸ˆà¹‰à¸‡à¹€à¸•à¸·à¸­à¸™',
  `dtmCreated` datetime NOT NULL DEFAULT (now()) COMMENT 'à¸§à¸±à¸™à¹€à¸§à¸¥à¸²à¸—à¸µà¹ˆà¸ªà¸£à¹‰à¸²à¸‡à¸à¸²à¸£à¹à¸ˆà¹‰à¸‡à¹€à¸•à¸·à¸­à¸™',
  PRIMARY KEY (`notificationId`) USING BTREE,
  KEY `notification_id` (`notificationId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.radio_devices
CREATE TABLE IF NOT EXISTS `radio_devices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `deviceId` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deviceName` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deviceType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `locationId` int DEFAULT NULL,
  `frequency` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `radioNo` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `radioName` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `serialNo` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `channel` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `dtmCreated` datetime DEFAULT NULL,
  `dtmUpdated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.radio_locations
CREATE TABLE IF NOT EXISTS `radio_locations` (
  `logId` int DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  `locationName` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latitude` decimal(12,9) DEFAULT NULL,
  `longitude` decimal(12,9) DEFAULT NULL,
  `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `dtmCreated` datetime DEFAULT NULL,
  `dtmUpdated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.record_media
CREATE TABLE IF NOT EXISTS `record_media` (
  `id` int NOT NULL AUTO_INCREMENT,
  `egressId` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `room` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fileName` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `filePath` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fileSize` int DEFAULT NULL,
  `duration` int DEFAULT NULL,
  `recordType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'completed',
  `dtmCreated` datetime DEFAULT NULL,
  `dtmCompleted` datetime DEFAULT NULL,
  `hls` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `encode` int DEFAULT NULL COMMENT '0=à¸£à¸­ Encode, 1=à¸ªà¸³à¹€à¸£à¹‡à¸ˆ, 2=à¸¥à¹‰à¸¡à¹€à¸«à¸¥à¸§, 3=à¸žà¸±à¸à¹„à¸Ÿà¸¥à¹Œ',
  `uploader` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `startRecord` datetime DEFAULT NULL,
  `endRecord` datetime DEFAULT NULL,
  `dtmUpdated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=542 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.room_conference
CREATE TABLE IF NOT EXISTS `room_conference` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `nodeLivekitId` int DEFAULT NULL,
  `status` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `roomType` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `room` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `service` int DEFAULT NULL,
  `recordStatus` int DEFAULT '0',
  `recordId` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `autoRecord` int DEFAULT '0',
  `recordType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `encodingOptionsPreset` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chatEnabled` int DEFAULT '0',
  `messageUnread` int DEFAULT '0',
  `agentSeen` datetime DEFAULT NULL,
  `userAgent` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `webSocketURL` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `dtmCreated` datetime DEFAULT NULL,
  `dtmUpdated` datetime DEFAULT NULL,
  `dtmClosed` datetime DEFAULT NULL,
  `dtmExpired` datetime DEFAULT NULL,
  `dtmRoomStarted` datetime DEFAULT NULL,
  `dtmRoomFinished` datetime DEFAULT NULL,
  `dtmStartRecord` datetime DEFAULT NULL,
  `dtmStopRecord` datetime DEFAULT NULL,
  `sync_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room_conference_status_dtmexpired` (`status`,`dtmExpired`),
  KEY `FK_room_conference_node_livekit` (`nodeLivekitId`) USING BTREE,
  KEY `idx_room_conference_room` (`room`),
  KEY `idx_room_conference_service` (`service`),
  CONSTRAINT `FK_room_conference_node_livekit` FOREIGN KEY (`nodeLivekitId`) REFERENCES `node_livekit` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28648 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.room_user
CREATE TABLE IF NOT EXISTS `room_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `room` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `identity` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `color` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `socketId` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `conference` int DEFAULT '1',
  `cameraMicrophoneStatus` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `camera` tinyint(1) DEFAULT '1',
  `microphone` tinyint(1) DEFAULT '1',
  `latitude` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `longitude` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `accuracy` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `userAgent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `dtmcreated` datetime DEFAULT NULL,
  `dtmupdated` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room_user_identity` (`identity`),
  KEY `idx_room_user_socketId` (`socketId`)
) ENGINE=InnoDB AUTO_INCREMENT=1563 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.services
CREATE TABLE IF NOT EXISTS `services` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `webTitle` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `prefixHLSRecordVideoSMS` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `prefixTextVideoSMS` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `prefixTextLocationSMS` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `domainsVideo` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `domainsLocation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `smsSenderName` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `logo` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `titleColor` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latitude` decimal(12,9) DEFAULT NULL,
  `longitude` decimal(12,9) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table conference.usage_status_log
CREATE TABLE IF NOT EXISTS `usage_status_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `linkID` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
  `room` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
  `mobile` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `linkType` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `latitude` decimal(12,9) DEFAULT NULL,
  `longitude` decimal(12,9) DEFAULT NULL,
  `identity` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
  `userName` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `userType` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `userAgent` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `data` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `dtmCreated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_usage_status_log_linkid` (`linkID`),
  KEY `idx_usage_status_log_linkid_status_userType_dtmCreated` (`linkID`,`status`,`userType`,`dtmCreated`),
  KEY `idx_usage_status_log_room` (`room`)
) ENGINE=InnoDB AUTO_INCREMENT=8515 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Data exporting was unselected.

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
