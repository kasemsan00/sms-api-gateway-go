package socket

// SocketEvents contains all socket event names
var SocketEvents = struct {
	// Connection events
	Connection        string
	Disconnect        string
	ConnectionSuccess string
	UserConnection    string
	UserDisconnect    string

	// Position events
	UserPosition         string
	GetUserPositionGroup string
	UserPositionGroup    string
	CarTracking          string
	Position             string
	AgentCar             string
	Destination          string

	// Chat events
	JoinChat       string
	ChatMessage    string
	GetChatHistory string
	ChatHistory    string

	// User events
	UserConference         string
	CameraMicrophoneStatus string
	GetUserDetail          string
	UserDetail             string
	GetUserList            string
	UserList               string
	UpdateUsername         string
	UpdateAdminUser        string
	UpdateMicrophoneStatus string
	UpdateCameraStatus     string
	TrackMutedUnmuted      string

	// Conference events
	AuthJoinConference       string
	AuthJoinConferenceAnswer string
	AgentCloseRoom           string
	ForceLeaveConference     string
	ForceStopSharescreen     string

	// Room events
	RoomRecord   string
	CaseData     string
	CaseListData string
	CaseListPage string

	// Queue events
	Queue string

	// Common events
	Error   string
	Log     string
	Message string
	Status  string
}{
	// Connection events
	Connection:        "connection",
	Disconnect:        "disconnect",
	ConnectionSuccess: "connection-success",
	UserConnection:    "user-connection",
	UserDisconnect:    "user-disconnect",

	// Position events
	UserPosition:         "user-position",
	GetUserPositionGroup: "get-user-position-group",
	UserPositionGroup:    "user-position-group",
	CarTracking:          "carTracking",
	Position:             "position",
	AgentCar:             "agentCar",
	Destination:          "destination",

	// Chat events
	JoinChat:       "joinChat",
	ChatMessage:    "chat-message",
	GetChatHistory: "get-chat-history",
	ChatHistory:    "chat-history",

	// User events
	UserConference:         "user-conference",
	CameraMicrophoneStatus: "camera-microphone-status",
	GetUserDetail:          "get-user-detail",
	UserDetail:             "user-detail",
	GetUserList:            "get-user-list",
	UserList:               "user-list",
	UpdateUsername:         "update-username",
	UpdateAdminUser:        "update-admin-user",
	UpdateMicrophoneStatus: "update-microphone-status",
	UpdateCameraStatus:     "update-camera-status",
	TrackMutedUnmuted:      "track-muted-unmuted",

	// Conference events
	AuthJoinConference:       "auth-join-conference",
	AuthJoinConferenceAnswer: "auth-join-conference-answer",
	AgentCloseRoom:           "agent-close-room",
	ForceLeaveConference:     "force-leave-conference",
	ForceStopSharescreen:     "force-stop-sharescreen",

	// Room events
	RoomRecord:   "room-record",
	CaseData:     "case-data",
	CaseListData: "case-list-data",
	CaseListPage: "case-list-page",

	// Queue events
	Queue: "queue",

	// Common events
	Error:   "error",
	Log:     "log",
	Message: "message",
	Status:  "status",
}

// SocketStatus contains socket status values
var SocketStatus = struct {
	Connection string
	Disconnect string
	Wait       string
	Open       string
	Close      string
	Success    string
	Failed     string
	Waiting    string
	Admit      string
	Deny       string
}{
	Connection: "connection",
	Disconnect: "disconnect",
	Wait:       "wait",
	Open:       "open",
	Close:      "close",
	Success:    "success",
	Failed:     "failed",
	Waiting:    "waiting",
	Admit:      "admit",
	Deny:       "deny",
}

// ErrorCodes contains socket error codes
var ErrorCodes = struct {
	AuthenticationError   string
	ValidationError       string
	DatabaseError         string
	NetworkError          string
	PermissionError       string
	ResourceNotFound      string
	InternalError         string
	PositionUpdateError   string
	PositionGroupError    string
	CarTrackingError      string
	PositionError         string
	ChatJoinError         string
	ChatMessageError      string
	ChatHistoryError      string
	ConferenceUpdateError string
	DeviceStatusError     string
	UserDetailError       string
	UserListError         string
	UsernameUpdateError   string
	AdminUpdateError      string
	MicrophoneUpdateError string
	CameraUpdateError     string
	TrackStatusError      string
	RoomRecordError       string
	CaseDataError         string
	CaseListError         string
	CasePageError         string
	ConferenceAccessError string
	ConferenceAnswerError string
	RoomCloseError        string
	ForceLeaveError       string
	ScreenShareError      string
	QueueError            string
}{
	AuthenticationError:   "AUTHENTICATION_ERROR",
	ValidationError:       "VALIDATION_ERROR",
	DatabaseError:         "DATABASE_ERROR",
	NetworkError:          "NETWORK_ERROR",
	PermissionError:       "PERMISSION_ERROR",
	ResourceNotFound:      "RESOURCE_NOT_FOUND",
	InternalError:         "INTERNAL_ERROR",
	PositionUpdateError:   "POSITION_UPDATE_ERROR",
	PositionGroupError:    "POSITION_GROUP_ERROR",
	CarTrackingError:      "CAR_TRACKING_ERROR",
	PositionError:         "POSITION_ERROR",
	ChatJoinError:         "CHAT_JOIN_ERROR",
	ChatMessageError:      "CHAT_MESSAGE_ERROR",
	ChatHistoryError:      "CHAT_HISTORY_ERROR",
	ConferenceUpdateError: "CONFERENCE_UPDATE_ERROR",
	DeviceStatusError:     "DEVICE_STATUS_ERROR",
	UserDetailError:       "USER_DETAIL_ERROR",
	UserListError:         "USER_LIST_ERROR",
	UsernameUpdateError:   "USERNAME_UPDATE_ERROR",
	AdminUpdateError:      "ADMIN_UPDATE_ERROR",
	MicrophoneUpdateError: "MICROPHONE_UPDATE_ERROR",
	CameraUpdateError:     "CAMERA_UPDATE_ERROR",
	TrackStatusError:      "TRACK_STATUS_ERROR",
	RoomRecordError:       "ROOM_RECORD_ERROR",
	CaseDataError:         "CASE_DATA_ERROR",
	CaseListError:         "CASE_LIST_ERROR",
	CasePageError:         "CASE_PAGE_ERROR",
	ConferenceAccessError: "CONFERENCE_ACCESS_ERROR",
	ConferenceAnswerError: "CONFERENCE_ANSWER_ERROR",
	RoomCloseError:        "ROOM_CLOSE_ERROR",
	ForceLeaveError:       "FORCE_LEAVE_ERROR",
	ScreenShareError:      "SCREEN_SHARE_ERROR",
	QueueError:            "QUEUE_ERROR",
}

// UserTypes contains user type values
var UserTypes = struct {
	User  string
	Admin string
}{
	User:  "user",
	Admin: "admin",
}

// MobileDefaults contains default values for mobile namespace
var MobileDefaults = struct {
	Latitude   float64
	Longitude  float64
	Accuracy   int
	Status     string
	DtmCreated string
	DtmUpdated string
}{
	Latitude:   13.734,
	Longitude:  100.567,
	Accuracy:   100,
	Status:     "open || arrive || cancel || complete",
	DtmCreated: "2022-01-01 00:00:00",
	DtmUpdated: "2022-01-01 00:00:00",
}

// AllowedNamespaces contains the allowed namespace names
var AllowedNamespaces = []string{"queue", "newqueue", "mobile", "notification"}

// IsAllowedNamespace checks if a namespace is in the allowed list
func IsAllowedNamespace(namespace string) bool {
	for _, ns := range AllowedNamespaces {
		if ns == namespace {
			return true
		}
	}
	return false
}
