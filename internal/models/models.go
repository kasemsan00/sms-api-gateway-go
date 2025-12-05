package models

import (
	"database/sql"
	"time"
)

// CarTrack represents the car_track table
type CarTrack struct {
	ID               int              `db:"id" json:"id"`
	UID              sql.NullString   `db:"uid" json:"uid,omitempty"`
	Status           sql.NullString   `db:"status" json:"status,omitempty"`
	Mobile           sql.NullString   `db:"mobile" json:"mobile,omitempty"`
	UserName         sql.NullString   `db:"userName" json:"userName,omitempty"`
	Room             sql.NullString   `db:"room" json:"room,omitempty"`
	Latitude         sql.NullFloat64  `db:"latitude" json:"latitude,omitempty"`
	Longitude        sql.NullFloat64  `db:"longitude" json:"longitude,omitempty"`
	Accuracy         sql.NullFloat64  `db:"accuracy" json:"accuracy,omitempty"`
	Speed            sql.NullInt32    `db:"speed" json:"speed,omitempty"`
	Heading          sql.NullInt32    `db:"heading" json:"heading,omitempty"`
	Altitude         sql.NullFloat64  `db:"altitude" json:"altitude,omitempty"`
	AltitudeAccuracy sql.NullFloat64  `db:"altitudeAccuracy" json:"altitudeAccuracy,omitempty"`
	DtmUpdated       sql.NullTime     `db:"dtmUpdated" json:"dtmUpdated,omitempty"`
	DtmCreated       sql.NullTime     `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmStarted       sql.NullTime     `db:"dtmStarted" json:"dtmStarted,omitempty"`
	DtmArrived       sql.NullTime     `db:"dtmArrived" json:"dtmArrived,omitempty"`
	DtmCanceled      sql.NullTime     `db:"dtmCanceled" json:"dtmCanceled,omitempty"`
	DtmCompleted     sql.NullTime     `db:"dtmCompleted" json:"dtmCompleted,omitempty"`
}

// CaseData represents the case_data table
type CaseData struct {
	ID              uint           `db:"id" json:"id"`
	CaseID          int            `db:"caseId" json:"caseId"`
	Service         sql.NullInt32  `db:"service" json:"service,omitempty"`
	RoomID          sql.NullInt32  `db:"roomId" json:"roomId,omitempty"`
	OperationNumber sql.NullString `db:"operationNumber" json:"operationNumber,omitempty"`
	Status          sql.NullString `db:"status" json:"status,omitempty"`
	HN              sql.NullString `db:"hn" json:"hn,omitempty"`
	PatientMobile   sql.NullString `db:"patientMobile" json:"patientMobile,omitempty"`
	MobileCreated   sql.NullString `db:"mobileCreated" json:"mobileCreated,omitempty"`
	CaseType        sql.NullString `db:"caseType" json:"caseType,omitempty"`
	UserName        sql.NullString `db:"userName" json:"userName,omitempty"`
	DtmCreated      sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
	Organization    sql.NullString `db:"organization" json:"organization,omitempty"`
}

// ChatMessage represents the chat_message table
type ChatMessage struct {
	ID               int            `db:"id" json:"id"`
	Room             sql.NullString `db:"room" json:"room,omitempty"`
	Identity         sql.NullString `db:"identity" json:"identity,omitempty"`
	ChatIdentity     sql.NullString `db:"chat_identity" json:"chatIdentity,omitempty"`
	UserName         sql.NullString `db:"userName" json:"userName,omitempty"`
	Text             sql.NullString `db:"text" json:"text,omitempty"`
	Color            sql.NullString `db:"color" json:"color,omitempty"`
	Files            sql.NullString `db:"files" json:"files,omitempty"`
	ReplyToMessageID sql.NullInt32  `db:"replyToMessageId" json:"replyToMessageId,omitempty"`
	ReplyToUserName  sql.NullString `db:"replyToUserName" json:"replyToUserName,omitempty"`
	ReplyToText      sql.NullString `db:"replyToText" json:"replyToText,omitempty"`
	DtmCreated       sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
	UserType         sql.NullString `db:"userType" json:"userType,omitempty"`
}

// ColorScheme represents the color_scheme table
type ColorScheme struct {
	ID       uint           `db:"id" json:"id"`
	ColorHex sql.NullString `db:"color_hex" json:"colorHex,omitempty"`
}

// DataLog represents the data_log table
type DataLog struct {
	ID         int            `db:"id" json:"id"`
	Data       sql.NullString `db:"data" json:"data,omitempty"`
	DtmCreated sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
}

// Files represents the files table
type Files struct {
	ID        int            `db:"id" json:"id"`
	LinkID    sql.NullString `db:"linkId" json:"linkId,omitempty"`
	ElementID sql.NullString `db:"elementId" json:"elementId,omitempty"`
	Filename  string         `db:"filename" json:"filename"`
	URL       string         `db:"url" json:"url"`
	Thumbnail sql.NullString `db:"thumbnail" json:"thumbnail,omitempty"`
	FileType  sql.NullString `db:"fileType" json:"fileType,omitempty"`
	Size      int64          `db:"size" json:"size"`
	Mimetype  sql.NullString `db:"mimetype" json:"mimetype,omitempty"`
	Width     sql.NullInt32  `db:"width" json:"width,omitempty"`
	Height    sql.NullInt32  `db:"height" json:"height,omitempty"`
	CreatedAt sql.NullTime   `db:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt sql.NullTime   `db:"updatedAt" json:"updatedAt,omitempty"`
}

// LinkConnect represents the link_connect table
type LinkConnect struct {
	ID                    uint            `db:"id" json:"id"`
	RoomUserID            sql.NullInt32   `db:"roomUserId" json:"roomUserId,omitempty"`
	SMS                   sql.NullInt32   `db:"sms" json:"sms,omitempty"`
	RecordID              sql.NullInt32   `db:"recordId" json:"recordId,omitempty"`
	Mobile                string          `db:"mobile" json:"mobile"`
	LinkID                sql.NullString  `db:"linkID" json:"linkID,omitempty"`
	DomainIndex           sql.NullInt32   `db:"domainIndex" json:"domainIndex,omitempty"`
	Share                 sql.NullInt32   `db:"share" json:"share,omitempty"`
	Enabled               sql.NullInt32   `db:"enabled" json:"enabled,omitempty"`
	UserName              sql.NullString  `db:"userName" json:"userName,omitempty"`
	Room                  sql.NullString  `db:"room" json:"room,omitempty"`
	UserType              sql.NullString  `db:"userType" json:"userType,omitempty"`
	LinkType              sql.NullString  `db:"linkType" json:"linkType,omitempty"`
	CrmSender             sql.NullString  `db:"crmSender" json:"crmSender,omitempty"`
	Accuracy              sql.NullString  `db:"accuracy" json:"accuracy,omitempty"`
	Latitude              sql.NullFloat64 `db:"latitude" json:"latitude,omitempty"`
	Longitude             sql.NullFloat64 `db:"longitude" json:"longitude,omitempty"`
	PatientLatitude       sql.NullFloat64 `db:"patientLatitude" json:"patientLatitude,omitempty"`
	PatientLongitude      sql.NullFloat64 `db:"patientLongitude" json:"patientLongitude,omitempty"`
	PatientUpdated        sql.NullTime    `db:"patientUpdated" json:"patientUpdated,omitempty"`
	Service               sql.NullInt32   `db:"service" json:"service,omitempty"`
	ErrorVideo            sql.NullString  `db:"errorVideo" json:"errorVideo,omitempty"`
	ErrorLocation         sql.NullString  `db:"errorLocation" json:"errorLocation,omitempty"`
	OS                    sql.NullString  `db:"os" json:"os,omitempty"`
	UserAgent             sql.NullString  `db:"userAgent" json:"userAgent,omitempty"`
	RequireJoinPermission sql.NullInt32   `db:"requireJoinPermission" json:"requireJoinPermission,omitempty"`
	RequireUserName       sql.NullInt32   `db:"requireUserName" json:"requireUserName,omitempty"`
	RequirePassword       sql.NullInt32   `db:"requirePassword" json:"requirePassword,omitempty"`
	OneTimeLink           sql.NullInt32   `db:"oneTimeLink" json:"oneTimeLink,omitempty"`
	Password              sql.NullString  `db:"password" json:"password,omitempty"`
	IsAdmin               sql.NullString  `db:"isAdmin" json:"isAdmin,omitempty"`
	DtmConnection         sql.NullTime    `db:"dtmConnection" json:"dtmConnection,omitempty"`
	DtmDisconnect         sql.NullTime    `db:"dtmDisconnect" json:"dtmDisconnect,omitempty"`
	DtmCreated            sql.NullTime    `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmExpired            sql.NullTime    `db:"dtmExpired" json:"dtmExpired,omitempty"`
	SyncAt                sql.NullTime    `db:"sync_at" json:"syncAt,omitempty"`
}

// NodeLivekit represents the node_livekit table
type NodeLivekit struct {
	ID               int            `db:"id" json:"id"`
	NodeName         sql.NullString `db:"nodeName" json:"nodeName,omitempty"`
	LivekitHost      sql.NullString `db:"livekitHost" json:"livekitHost,omitempty"`
	LivekitLocal     sql.NullString `db:"livekitLocal" json:"livekitLocal,omitempty"`
	LivekitApiKey    sql.NullString `db:"livekitApiKey" json:"livekitApiKey,omitempty"`
	LivekitApiSecret sql.NullString `db:"livekitApiSecret" json:"livekitApiSecret,omitempty"`
	LastHealthCheck  sql.NullTime   `db:"lastHealthCheck" json:"lastHealthCheck,omitempty"`
	Description      sql.NullString `db:"description" json:"description,omitempty"`
}

// Notification represents the notification table
type Notification struct {
	NotificationID   int            `db:"notificationId" json:"notificationId"`
	UserName         sql.NullString `db:"userName" json:"userName,omitempty"`
	Mobile           sql.NullString `db:"mobile" json:"mobile,omitempty"`
	Message          sql.NullString `db:"message" json:"message,omitempty"`
	CaseID           sql.NullInt32  `db:"caseId" json:"caseId,omitempty"`
	Read             int8           `db:"read" json:"read"`
	NotificationType sql.NullString `db:"notificationType" json:"notificationType,omitempty"`
	RelatedURL       sql.NullString `db:"relatedUrl" json:"relatedUrl,omitempty"`
	DtmRead          sql.NullTime   `db:"dtmRead" json:"dtmRead,omitempty"`
	DtmCreated       time.Time      `db:"dtmCreated" json:"dtmCreated"`
}

// RadioDevices represents the radio_devices table
type RadioDevices struct {
	ID         int            `db:"id" json:"id"`
	DeviceID   sql.NullString `db:"deviceId" json:"deviceId,omitempty"`
	DeviceName sql.NullString `db:"deviceName" json:"deviceName,omitempty"`
	DeviceType sql.NullString `db:"deviceType" json:"deviceType,omitempty"`
	Status     sql.NullString `db:"status" json:"status,omitempty"`
	LocationID sql.NullInt32  `db:"locationId" json:"locationId,omitempty"`
	Frequency  sql.NullString `db:"frequency" json:"frequency,omitempty"`
	RadioNo    sql.NullString `db:"radioNo" json:"radioNo,omitempty"`
	RadioName  sql.NullString `db:"radioName" json:"radioName,omitempty"`
	SerialNo   sql.NullString `db:"serialNo" json:"serialNo,omitempty"`
	Channel    sql.NullString `db:"channel" json:"channel,omitempty"`
	DtmCreated sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmUpdated sql.NullTime   `db:"dtmUpdated" json:"dtmUpdated,omitempty"`
}

// RadioLocations represents the radio_locations table
type RadioLocations struct {
	LogID        sql.NullInt32   `db:"logId" json:"logId,omitempty"`
	ID           int             `db:"id" json:"id"`
	LocationName sql.NullString  `db:"locationName" json:"locationName,omitempty"`
	Latitude     sql.NullFloat64 `db:"latitude" json:"latitude,omitempty"`
	Longitude    sql.NullFloat64 `db:"longitude" json:"longitude,omitempty"`
	Address      sql.NullString  `db:"address" json:"address,omitempty"`
	Description  sql.NullString  `db:"description" json:"description,omitempty"`
	Status       sql.NullString  `db:"status" json:"status,omitempty"`
	DtmCreated   sql.NullTime    `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmUpdated   sql.NullTime    `db:"dtmUpdated" json:"dtmUpdated,omitempty"`
}

// RecordMedia represents the record_media table
type RecordMedia struct {
	ID           int            `db:"id" json:"id"`
	EgressID     sql.NullString `db:"egressId" json:"egressId,omitempty"`
	Room         sql.NullString `db:"room" json:"room,omitempty"`
	FileName     sql.NullString `db:"fileName" json:"fileName,omitempty"`
	FilePath     sql.NullString `db:"filePath" json:"filePath,omitempty"`
	FileSize     sql.NullInt32  `db:"fileSize" json:"fileSize,omitempty"`
	Duration     sql.NullInt32  `db:"duration" json:"duration,omitempty"`
	RecordType   sql.NullString `db:"recordType" json:"recordType,omitempty"`
	Status       sql.NullString `db:"status" json:"status,omitempty"`
	DtmCreated   sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmCompleted sql.NullTime   `db:"dtmCompleted" json:"dtmCompleted,omitempty"`
	HLS          sql.NullString `db:"hls" json:"hls,omitempty"`
	Encode       sql.NullInt32  `db:"encode" json:"encode,omitempty"`
	Uploader     sql.NullString `db:"uploader" json:"uploader,omitempty"`
	StartRecord  sql.NullTime   `db:"startRecord" json:"startRecord,omitempty"`
	EndRecord    sql.NullTime   `db:"endRecord" json:"endRecord,omitempty"`
	DtmUpdated   sql.NullTime   `db:"dtmUpdated" json:"dtmUpdated,omitempty"`
}

// RoomConference represents the room_conference table
type RoomConference struct {
	ID                    uint           `db:"id" json:"id"`
	NodeLivekitID         sql.NullInt32  `db:"nodeLivekitId" json:"nodeLivekitId,omitempty"`
	Status                sql.NullString `db:"status" json:"status,omitempty"`
	RoomType              sql.NullString `db:"roomType" json:"roomType,omitempty"`
	Room                  sql.NullString `db:"room" json:"room,omitempty"`
	Service               sql.NullInt32  `db:"service" json:"service,omitempty"`
	RecordStatus          sql.NullInt32  `db:"recordStatus" json:"recordStatus,omitempty"`
	RecordID              sql.NullString `db:"recordId" json:"recordId,omitempty"`
	AutoRecord            sql.NullInt32  `db:"autoRecord" json:"autoRecord,omitempty"`
	RecordType            sql.NullString `db:"recordType" json:"recordType,omitempty"`
	EncodingOptionsPreset sql.NullString `db:"encodingOptionsPreset" json:"encodingOptionsPreset,omitempty"`
	ChatEnabled           sql.NullInt32  `db:"chatEnabled" json:"chatEnabled,omitempty"`
	MessageUnread         sql.NullInt32  `db:"messageUnread" json:"messageUnread,omitempty"`
	AgentSeen             sql.NullTime   `db:"agentSeen" json:"agentSeen,omitempty"`
	UserAgent             sql.NullString `db:"userAgent" json:"userAgent,omitempty"`
	WebSocketURL          sql.NullString `db:"webSocketURL" json:"webSocketURL,omitempty"`
	DtmCreated            sql.NullTime   `db:"dtmCreated" json:"dtmCreated,omitempty"`
	DtmUpdated            sql.NullTime   `db:"dtmUpdated" json:"dtmUpdated,omitempty"`
	DtmClosed             sql.NullTime   `db:"dtmClosed" json:"dtmClosed,omitempty"`
	DtmExpired            sql.NullTime   `db:"dtmExpired" json:"dtmExpired,omitempty"`
	DtmRoomStarted        sql.NullTime   `db:"dtmRoomStarted" json:"dtmRoomStarted,omitempty"`
	DtmRoomFinished       sql.NullTime   `db:"dtmRoomFinished" json:"dtmRoomFinished,omitempty"`
	DtmStartRecord        sql.NullTime   `db:"dtmStartRecord" json:"dtmStartRecord,omitempty"`
	DtmStopRecord         sql.NullTime   `db:"dtmStopRecord" json:"dtmStopRecord,omitempty"`
	SyncAt                sql.NullTime   `db:"sync_at" json:"syncAt,omitempty"`
}

// RoomUser represents the room_user table
type RoomUser struct {
	ID                     int            `db:"id" json:"id"`
	Room                   sql.NullString `db:"room" json:"room,omitempty"`
	Identity               sql.NullString `db:"identity" json:"identity,omitempty"`
	Color                  sql.NullString `db:"color" json:"color,omitempty"`
	UserName               sql.NullString `db:"userName" json:"userName,omitempty"`
	UserType               sql.NullString `db:"userType" json:"userType,omitempty"`
	Status                 sql.NullString `db:"status" json:"status,omitempty"`
	SocketID               sql.NullString `db:"socketId" json:"socketId,omitempty"`
	Conference             sql.NullInt32  `db:"conference" json:"conference,omitempty"`
	CameraMicrophoneStatus sql.NullString `db:"cameraMicrophoneStatus" json:"cameraMicrophoneStatus,omitempty"`
	Camera                 sql.NullBool   `db:"camera" json:"camera,omitempty"`
	Microphone             sql.NullBool   `db:"microphone" json:"microphone,omitempty"`
	Latitude               sql.NullString `db:"latitude" json:"latitude,omitempty"`
	Longitude              sql.NullString `db:"longitude" json:"longitude,omitempty"`
	Accuracy               sql.NullString `db:"accuracy" json:"accuracy,omitempty"`
	UserAgent              sql.NullString `db:"userAgent" json:"userAgent,omitempty"`
	DtmCreated             sql.NullTime   `db:"dtmcreated" json:"dtmCreated,omitempty"`
	DtmUpdated             time.Time      `db:"dtmupdated" json:"dtmUpdated"`
}

// Services represents the services table
type Services struct {
	ID                      uint            `db:"id" json:"id"`
	Name                    sql.NullString  `db:"name" json:"name,omitempty"`
	WebTitle                sql.NullString  `db:"webTitle" json:"webTitle,omitempty"`
	PrefixHLSRecordVideoSMS sql.NullString  `db:"prefixHLSRecordVideoSMS" json:"prefixHLSRecordVideoSMS,omitempty"`
	PrefixTextVideoSMS      sql.NullString  `db:"prefixTextVideoSMS" json:"prefixTextVideoSMS,omitempty"`
	PrefixTextLocationSMS   sql.NullString  `db:"prefixTextLocationSMS" json:"prefixTextLocationSMS,omitempty"`
	DomainsVideo            sql.NullString  `db:"domainsVideo" json:"domainsVideo,omitempty"`
	DomainsLocation         sql.NullString  `db:"domainsLocation" json:"domainsLocation,omitempty"`
	SmsSenderName           sql.NullString  `db:"smsSenderName" json:"smsSenderName,omitempty"`
	Logo                    sql.NullString  `db:"logo" json:"logo,omitempty"`
	TitleColor              sql.NullString  `db:"titleColor" json:"titleColor,omitempty"`
	Latitude                sql.NullFloat64 `db:"latitude" json:"latitude,omitempty"`
	Longitude               sql.NullFloat64 `db:"longitude" json:"longitude,omitempty"`
}

// UsageStatusLog represents the usage_status_log table
type UsageStatusLog struct {
	ID         uint            `db:"id" json:"id"`
	LinkID     sql.NullString  `db:"linkID" json:"linkID,omitempty"`
	Room       sql.NullString  `db:"room" json:"room,omitempty"`
	Mobile     sql.NullString  `db:"mobile" json:"mobile,omitempty"`
	LinkType   sql.NullString  `db:"linkType" json:"linkType,omitempty"`
	Latitude   sql.NullFloat64 `db:"latitude" json:"latitude,omitempty"`
	Longitude  sql.NullFloat64 `db:"longitude" json:"longitude,omitempty"`
	Identity   sql.NullString  `db:"identity" json:"identity,omitempty"`
	UserName   sql.NullString  `db:"userName" json:"userName,omitempty"`
	UserType   sql.NullString  `db:"userType" json:"userType,omitempty"`
	Status     sql.NullString  `db:"status" json:"status,omitempty"`
	UserAgent  sql.NullString  `db:"userAgent" json:"userAgent,omitempty"`
	Data       sql.NullString  `db:"data" json:"data,omitempty"`
	DtmCreated sql.NullTime    `db:"dtmCreated" json:"dtmCreated,omitempty"`
}

// RoomDetailResponse is the response for room detail
type RoomDetailResponse struct {
	ID                    uint      `json:"id"`
	Status                string    `json:"status"`
	Room                  string    `json:"room"`
	RoomType              string    `json:"roomType"`
	RecordID              string    `json:"recordId"`
	AutoRecord            int       `json:"autoRecord"`
	RecordType            string    `json:"recordType"`
	EncodingOptionsPreset string    `json:"encodingOptionsPreset"`
	ChatEnabled           int       `json:"chatEnabled"`
	MessageUnread         int       `json:"messageUnread"`
	UserAgent             string    `json:"userAgent"`
	DtmCreated            time.Time `json:"dtmCreated"`
	DtmStartRecord        time.Time `json:"dtmStartRecord,omitempty"`
	DtmStopRecord         time.Time `json:"dtmStopRecord,omitempty"`
	WebSocketURL          string    `json:"webSocketURL"`
	RecordDuration        int64     `json:"recordDuration,omitempty"`
}

// UserDetailResponse is the response for user detail
type UserDetailResponse struct {
	Room       string `json:"room"`
	Identity   string `json:"identity"`
	Color      string `json:"color"`
	UserName   string `json:"userName"`
	Status     string `json:"status"`
	Camera     bool   `json:"camera"`
	Microphone bool   `json:"microphone"`
	UserType   string `json:"userType"`
	Conference int    `json:"conference"`
}

// LinkDetailResponse is the response for link detail
type LinkDetailResponse struct {
	LinkID                string    `json:"linkID"`
	Room                  string    `json:"room"`
	Enabled               int       `json:"enabled"`
	Mobile                string    `json:"mobile"`
	IsAdmin               string    `json:"isAdmin"`
	UserName              string    `json:"userName"`
	UserType              string    `json:"userType"`
	LinkType              string    `json:"linkType"`
	RequireJoinPermission int       `json:"requireJoinPermission"`
	CrmSender             string    `json:"crmSender"`
	RequireUserName       int       `json:"requireUserName"`
	RequirePassword       int       `json:"requirePassword"`
	OneTimeLink           int       `json:"oneTimeLink"`
	DtmCreated            time.Time `json:"dtmCreated"`
	DtmExpired            time.Time `json:"dtmExpired"`
	Latitude              float64   `json:"latitude,omitempty"`
	Longitude             float64   `json:"longitude,omitempty"`
	PatientLatitude       float64   `json:"patientLatitude,omitempty"`
	PatientLongitude      float64   `json:"patientLongitude,omitempty"`
}
