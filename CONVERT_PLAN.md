# API Gateway: Node.js to Golang Conversion Plan

> แผนการแปลง API Gateway จาก Node.js Express เป็น Golang (Fiber)  
> **วันที่สร้าง**: 2025-12-05  
> **อัพเดทล่าสุด**: 2025-12-05  
> **สถานะ**: 🚧 In Progress (Phase 1-4 Complete)

---

## 📊 Implementation Progress Summary

| Phase       | Description                                                                | Status      |
| ----------- | -------------------------------------------------------------------------- | ----------- |
| **Phase 1** | Foundation Setup (Config, DB, Redis, LiveKit, Models, Utils, Repositories) | ✅ Complete |
| **Phase 2** | Core Services (14 services)                                                | ✅ Complete |
| **Phase 3** | HTTP Handlers (15 handlers)                                                | ✅ Complete |
| **Phase 4** | Middleware & Routing (Auth, CORS, Logger, Router)                          | ✅ Complete |
| **Phase 5** | WebSocket/Socket.IO Implementation                                         | ⏳ Pending  |
| **Phase 6** | LiveKit Webhook Processing                                                 | ⏳ Pending  |
| **Phase 7** | Testing & Integration                                                      | ⏳ Pending  |
| **Phase 8** | Documentation & Deployment                                                 | ⏳ Pending  |

### Files Created (Summary):

**Configuration (4 files):**

- `internal/config/config.go`
- `internal/config/database.go`
- `internal/config/redis.go`
- `internal/config/livekit.go`

**Repositories (11 files):**

- `room_repository.go`, `user_repository.go`, `link_repository.go`, `chat_repository.go`
- `notification_repository.go`, `record_repository.go`, `car_repository.go`
- `case_repository.go`, `radio_repository.go`, `stats_repository.go`, `usage_log_repository.go`

**Services (14 files):**

- `auth_service.go`, `room_service.go`, `user_service.go`, `link_service.go`
- `chat_service.go`, `notification_service.go`, `record_service.go`, `car_service.go`
- `case_service.go`, `radio_service.go`, `stats_service.go`, `sms_service.go`
- `file_service.go`, `crontab_service.go`

**Handlers (15 files):**

- `auth_handler.go`, `room_handler.go`, `user_handler.go`, `link_handler.go`
- `system_handler.go`, `chat_handler.go`, `notification_handler.go`, `record_handler.go`
- `car_handler.go`, `case_handler.go`, `radio_handler.go`, `stats_handler.go`
- `upload_handler.go`, `webhook_handler.go`, `test_handler.go`

**Middleware (3 files):**

- `auth_middleware.go`, `cors_middleware.go`, `logger_middleware.go`

**Router & Utils:**

- `internal/router/router.go`
- `pkg/logger/logger.go`
- `pkg/utils/response.go`, `pkg/utils/helpers.go`

---

## 📌 สรุปภาพรวมโปรเจค

### Node.js Original Stack

| Component          | Technology                    |
| ------------------ | ----------------------------- |
| Framework          | Express.js 5.x                |
| Real-time          | Socket.IO 4.8 + Redis Adapter |
| Database           | MySQL (mysql2)                |
| Cache              | Redis (ioredis)               |
| Video Conferencing | LiveKit Server SDK            |
| Authentication     | JWT (jsonwebtoken)            |
| File Upload        | Multer                        |
| Cron Jobs          | cron (node-cron)              |
| Logging            | Winston + Morgan              |

### Golang Target Stack

| Component          | Technology                           |
| ------------------ | ------------------------------------ |
| Framework          | Fiber v2                             |
| Real-time          | Socket.IO Go หรือ Melody (WebSocket) |
| Database           | sqlx + MySQL driver                  |
| Cache              | go-redis/v9                          |
| Video Conferencing | livekit-server-sdk-go                |
| Authentication     | golang-jwt/jwt/v5                    |
| File Upload        | Fiber built-in                       |
| Cron Jobs          | robfig/cron/v3                       |
| Logging            | zerolog หรือ zap                     |

---

## 📁 โครงสร้างโปรเจค Golang ที่แนะนำ

```
api-gateway-go/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration loader
│   │   ├── database.go             # Database connection
│   │   ├── redis.go                # Redis connection
│   │   └── livekit.go              # LiveKit configuration
│   ├── models/
│   │   └── models.go               # Database models (✅ มีแล้ว)
│   ├── repository/
│   │   ├── room_repository.go
│   │   ├── user_repository.go
│   │   ├── link_repository.go
│   │   ├── chat_repository.go
│   │   ├── notification_repository.go
│   │   ├── record_repository.go
│   │   ├── car_repository.go
│   │   ├── case_repository.go
│   │   ├── radio_repository.go
│   │   ├── stats_repository.go
│   │   └── usage_log_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── room_service.go
│   │   ├── user_service.go
│   │   ├── link_service.go
│   │   ├── chat_service.go
│   │   ├── notification_service.go
│   │   ├── record_service.go
│   │   ├── car_service.go
│   │   ├── case_service.go
│   │   ├── radio_service.go
│   │   ├── stats_service.go
│   │   ├── usage_log_service.go
│   │   ├── livekit_service.go
│   │   ├── sms_service.go
│   │   ├── file_service.go
│   │   └── crontab_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── room_handler.go
│   │   ├── user_handler.go
│   │   ├── link_handler.go
│   │   ├── chat_handler.go
│   │   ├── notification_handler.go
│   │   ├── record_handler.go
│   │   ├── car_handler.go
│   │   ├── case_handler.go
│   │   ├── radio_handler.go
│   │   ├── stats_handler.go
│   │   ├── upload_handler.go
│   │   ├── webhook_handler.go
│   │   ├── service_handler.go
│   │   └── test_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   └── logger_middleware.go
│   ├── socket/
│   │   ├── hub.go                  # WebSocket hub manager
│   │   ├── client.go               # WebSocket client
│   │   ├── room_socket.go          # Room namespace handlers
│   │   ├── mobile_socket.go        # Mobile namespace handlers
│   │   ├── notification_socket.go  # Notification namespace handlers
│   │   ├── queue_socket.go         # Queue namespace handlers
│   │   └── handlers/
│   │       ├── chat_handler.go
│   │       ├── position_handler.go
│   │       ├── conference_handler.go
│   │       └── user_handler.go
│   └── router/
│       └── router.go               # Route definitions
├── pkg/
│   ├── utils/
│   │   ├── response.go
│   │   ├── validator.go
│   │   └── helpers.go
│   └── logger/
│       └── logger.go
├── .env
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## 🚀 Phase 1: Foundation Setup (Week 1) ✅ COMPLETE

### 1.1 Project Configuration

- [x] สร้าง `internal/config/config.go` - โหลด environment variables
- [x] สร้าง `internal/config/database.go` - MySQL connection pool with sqlx
- [x] สร้าง `internal/config/redis.go` - Redis connection manager
- [x] สร้าง `internal/config/livekit.go` - LiveKit client configuration
- [x] สร้าง `pkg/logger/logger.go` - Structured logging (zerolog)

### 1.2 Database Layer

- [x] ตรวจสอบ/อัพเดท `internal/models/models.go` (มีอยู่แล้ว)
- [x] สร้าง All Repositories (11 files implemented)

### 1.3 Dependencies ที่ต้องเพิ่มใน go.mod

```go
require (
    github.com/gofiber/fiber/v2 v2.52.10
    github.com/gofiber/contrib/websocket v1.3.0
    github.com/jmoiron/sqlx v1.4.0
    github.com/go-sql-driver/mysql v1.8.1
    github.com/redis/go-redis/v9 v9.7.0
    github.com/livekit/server-sdk-go/v2 v2.4.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/joho/godotenv v1.5.1
    github.com/rs/zerolog v1.33.0
    github.com/robfig/cron/v3 v3.0.1
    github.com/google/uuid v1.6.0
    github.com/shopspring/decimal v1.4.0
)
```

---

## 🔌 Phase 2: Core Services Implementation (Week 2-3) ✅ COMPLETE

> **สถานะ**: ✅ All 14 services implemented
>
> - auth_service.go
> - room_service.go
> - user_service.go
> - link_service.go
> - chat_service.go
> - notification_service.go
> - record_service.go
> - car_service.go
> - case_service.go
> - radio_service.go
> - stats_service.go
> - sms_service.go
> - file_service.go
> - crontab_service.go

### 2.1 Repository Layer (Data Access)

#### 2.1.1 Room Repository

```go
// Node.js functions to convert:
// - createRoom
// - closeRoomById
// - closeRoomAll
// - updateRoomStatus
// - updateExpired
// - updateRoomType
// - deleteRoom
// - getRoomDetail
// - getRoomConferenceList
// - autoRoomExpiredClose
// - checkRoomExpired
// - updateRecordStatus
// - getRoomRecordStatus
// - updateRecordID
// - updateRoomStartedFinished
```

#### 2.1.2 User Repository

```go
// Node.js functions to convert:
// - addUser
// - updateUser
// - updateUserStatus
// - updateUserType
// - getUserDetail
// - getUserAlreadyInRoom
// - listParticipants
// - removeParticipant
// - updateSocketIOUser
// - getRoomUserId
// - agentList
```

#### 2.1.3 Link Repository

```go
// Node.js functions to convert:
// - createLink
// - getLinkDetail
// - updateLatLngLinkDetail
// - getSMSLinkHistory
// - getOneTimeLinkStatus
// - updateOneTimeLink
// - getLinkIdList
// - updateLinkEnabled
// - autoLinkExpiredClose
// - updateLinkConnectTime
```

#### 2.1.4 Chat Repository

```go
// Node.js functions to convert:
// - getChatHistory
// - addChatMessage
// - getChatNotification
```

#### 2.1.5 Notification Repository

```go
// Node.js functions to convert:
// - getAllNotifications
// - getUnreadNotifications
// - createNotification
// - updateNotificationReadStatus
// - getNotificationById
```

#### 2.1.6 Record Repository

```go
// Node.js functions to convert:
// - addRecordMedia
// - updateRecordMedia
// - getFileHistory
// - getRecordDetail
// - checkEgressAvailable
// - getRecordQueue
```

#### 2.1.7 Car Repository

```go
// Node.js functions to convert:
// - createCarTask
// - getTaskDetail
// - updateCarTask
// - updateCarPosition
// - getCarTaskList
```

#### 2.1.8 Case Repository

```go
// Node.js functions to convert:
// - createCase
// - getCaseById
// - updateCase
// - getCaseHistory
```

#### 2.1.9 Stats Repository

```go
// Node.js functions to convert:
// - getStatsSummary
// - getDeviceStats
// - getTypeStats
// - generateStats
// - getUserStats
// - getCaseStats
```

---

### 2.2 Service Layer (Business Logic)

#### 2.2.1 LiveKit Service (Critical)

```go
// internal/service/livekit_service.go

type LiveKitService interface {
    GetLiveKitNode(ctx context.Context, nodeName, room string, roomId int) (*NodeLivekit, error)
    GetLiveKitNodeAll(ctx context.Context) ([]NodeLivekit, error)
    GetRoomServiceClient(ctx context.Context, room string) (*lksdk.RoomServiceClient, error)
    UpdateHealthCheck(ctx context.Context, nodeName string) error
    CreateRoom(ctx context.Context, room string, opts ...lksdk.CreateRoomOption) (*livekit.Room, error)
    DeleteRoom(ctx context.Context, room string) error
}

// Node.js equivalent:
// - getLiveKitNode
// - getLiveKitNodeAll
// - getSVCLiveKit
// - updateHealthCheck
```

#### 2.2.2 User Service

```go
// internal/service/user_service.go

type UserService interface {
    GenerateUser(ctx context.Context, opts GenerateUserOptions) (*GenerateUserResult, error)
    GenerateUserJoinConference(ctx context.Context, room, userName, socketId string) (*UserToken, error)
    GetDomain(ctx context.Context, service int, sender, linkType, linkID string) (string, error)
    AddUser(ctx context.Context, opts AddUserOptions) error
    UpdateUserStatus(ctx context.Context, room, identity, status string) error
    GetUserDetail(ctx context.Context, room, identity, socketId string) (*RoomUser, error)
    ListParticipants(ctx context.Context, room string) ([]livekit.ParticipantInfo, error)
    RemoveParticipant(ctx context.Context, room, identity string) error
    MutePublishedTrack(ctx context.Context, room, identity, trackSid string, muted bool) error
}
```

#### 2.2.3 Room Service

```go
// internal/service/room_service.go

type RoomService interface {
    CreateRoom(ctx context.Context, opts CreateRoomOptions) (*RoomConference, error)
    CloseRoom(ctx context.Context, room string) error
    GetRoomDetail(ctx context.Context, room string) (*RoomConference, error)
    GetRoomConferenceList(ctx context.Context, status string) ([]RoomConference, error)
    UpdateRoomType(ctx context.Context, room string, opts UpdateRoomTypeOptions) error
    UpdateRecordStatus(ctx context.Context, room string, status int) error
    AutoRoomExpiredClose(ctx context.Context) error
    AutoRoomSocketClose(ctx context.Context) error
}
```

#### 2.2.4 Link Service

```go
// internal/service/link_service.go

type LinkService interface {
    CreateLink(ctx context.Context, opts CreateLinkOptions) (*LinkConnect, error)
    GetLinkDetail(ctx context.Context, linkID, room string, userType *string) (*LinkConnect, error)
    GetShareURL(ctx context.Context, room string, userType string) (string, error)
    UpdateLatLng(ctx context.Context, linkID string, lat, lng float64, accuracy int) error
    GetSMSLinkHistory(ctx context.Context, opts HistoryOptions) (*PaginatedResult, error)
    CheckAndUpdateOneTimeLink(ctx context.Context, linkID string) error
    AutoLinkExpiredClose(ctx context.Context) error
}
```

#### 2.2.5 Record Service (LiveKit Egress)

```go
// internal/service/record_service.go

type RecordService interface {
    StartRecord(ctx context.Context, opts StartRecordOptions) (*EgressInfo, error)
    StopRecord(ctx context.Context, recordId string) (*EgressInfo, error)
    ListEgress(ctx context.Context, room *string) ([]livekit.EgressInfo, error)
    StopAllActiveRecord(ctx context.Context) ([]StoppedEgress, error)
    GetFileHistory(ctx context.Context, room *string) ([]RecordMedia, error)
    CheckEgressAvailable(ctx context.Context) (bool, error)
}
```

#### 2.2.6 SMS Service

```go
// internal/service/sms_service.go

type SMSService interface {
    SendSMS(ctx context.Context, phoneNumber, message string) error
    SendCustomMessage(ctx context.Context, phoneNumber, message string) error
}
```

#### 2.2.7 Auth Service

```go
// internal/service/auth_service.go

type AuthService interface {
    CreateToken(ctx context.Context, payload interface{}) (string, error)
    VerifyToken(ctx context.Context, token string) (*Claims, error)
    VerifyUser(ctx context.Context, userName, password string) (*User, error)
}
```

#### 2.2.8 Crontab Service

```go
// internal/service/crontab_service.go

type CrontabService interface {
    InitCronJobs() error
    Cleanup() error
    GetStatus() *CronStatus
    HealthCheck() *HealthStatus
}

// Cron jobs to implement:
// - Room cleanup (*/30 * * * *)
// - Link cleanup (*/30 * * * *)
// - LiveKit health check (*/10 * * * *)
```

---

## 🌐 Phase 3: HTTP Handlers (REST API) (Week 3-4) ✅ COMPLETE

> **สถานะ**: ✅ All 15 handlers implemented
>
> - auth_handler.go
> - room_handler.go
> - user_handler.go
> - link_handler.go
> - system_handler.go
> - chat_handler.go
> - notification_handler.go
> - record_handler.go
> - car_handler.go
> - case_handler.go
> - radio_handler.go
> - stats_handler.go
> - upload_handler.go
> - webhook_handler.go
> - test_handler.go

### 3.1 Route Mapping (Node.js → Golang)

#### Auth Routes (`/auth`)

| Method | Endpoint           | Handler       | Description   |
| ------ | ------------------ | ------------- | ------------- |
| GET    | `/auth/create`     | `CreateToken` | สร้าง token   |
| POST   | `/auth/verifyuser` | `VerifyUser`  | ตรวจสอบผู้ใช้ |

#### Room Routes (`/room`)

| Method | Endpoint             | Handler          | Description          |
| ------ | -------------------- | ---------------- | -------------------- |
| GET    | `/room/detail`       | `GetRoomDetail`  | รายละเอียดห้อง       |
| GET    | `/room/listrooms`    | `ListRooms`      | รายการห้อง           |
| GET    | `/room/checkexpired` | `CheckExpired`   | ตรวจสอบหมดอายุ       |
| GET    | `/room/verifytoken`  | `VerifyToken`    | ตรวจสอบ token        |
| GET    | `/room/picture`      | `GetRoomPicture` | รูปภาพห้อง           |
| POST   | `/room/updateuser`   | `UpdateUser`     | อัพเดทผู้ใช้         |
| POST   | `/room/deleteroom`   | `DeleteRoom`     | ลบห้อง               |
| PUT    | `/room/updatetype`   | `UpdateType`     | อัพเดทประเภท         |
| PUT    | `/room/updatestatus` | `UpdateStatus`   | อัพเดทสถานะ          |
| PUT    | `/room/close`        | `CloseRoom`      | ปิดห้อง              |
| POST   | `/room/verifyuser`   | `VerifyUserRoom` | ตรวจสอบผู้ใช้ (auth) |

#### User Routes (`/user`)

| Method | Endpoint                     | Handler                | Description                |
| ------ | ---------------------------- | ---------------------- | -------------------------- |
| GET    | `/user/getuseralreadyinroom` | `GetUserAlreadyInRoom` | ผู้ใช้ในห้อง               |
| GET    | `/user/getuserdetail`        | `GetUserDetail`        | รายละเอียดผู้ใช้           |
| GET    | `/user/listparticipants`     | `ListParticipants`     | รายการผู้เข้าร่วม          |
| POST   | `/user/generate`             | `GenerateUser`         | สร้างผู้ใช้ (auth)         |
| POST   | `/user/joingenerate`         | `JoinGenerate`         | สร้างผู้ใช้เข้าร่วม (auth) |
| POST   | `/user/generateChatUser`     | `GenerateChatUser`     | สร้างผู้ใช้แชท             |
| POST   | `/user/updateparticipants`   | `UpdateParticipants`   | อัพเดทผู้เข้าร่วม          |
| POST   | `/user/mutepublishedtrack`   | `MutePublishedTrack`   | ปิดเสียง track             |
| POST   | `/user/removeParticipant`    | `RemoveParticipant`    | ลบผู้เข้าร่วม              |
| GET    | `/user/log`                  | `GetUserLog`           | บันทึกผู้ใช้               |
| PUT    | `/user/handle/track`         | `HandleTrack`          | จัดการ track               |

#### Link Routes (`/link`)

| Method | Endpoint                 | Handler          | Description        |
| ------ | ------------------------ | ---------------- | ------------------ |
| GET    | `/link/getdetail`        | `GetLinkDetail`  | รายละเอียดลิ้งค์   |
| GET    | `/link/history`          | `GetLinkHistory` | ประวัติลิ้งค์      |
| POST   | `/link/create`           | `CreateLink`     | สร้างลิ้งค์ (auth) |
| POST   | `/link/create/hls`       | `CreateHLSLink`  | สร้างลิ้งค์ HLS    |
| POST   | `/link/update/latlng`    | `UpdateLatLng`   | อัพเดทพิกัด        |
| POST   | `/link/multilatlng/send` | `MultiLatLng`    | ส่งพิกัดหลายจุด    |
| GET    | `/link/share`            | `GetShareURL`    | URL แชร์           |
| POST   | `/link/cartracking`      | `CarTracking`    | ติดตามรถ           |
| GET    | `/link/get/domain`       | `GetDomain`      | ดึงโดเมน           |
| GET    | `/link/list`             | `GetLinkList`    | รายการลิ้งค์       |

#### Case Routes (`/case`)

| Method | Endpoint        | Handler          | Description  |
| ------ | --------------- | ---------------- | ------------ |
| POST   | `/case/create`  | `CreateCase`     | สร้างเคส     |
| GET    | `/case/get`     | `GetCase`        | ดึงข้อมูลเคส |
| GET    | `/case/history` | `GetCaseHistory` | ประวัติเคส   |
| PUT    | `/case/update`  | `UpdateCase`     | อัพเดทเคส    |

#### Chat Routes (`/chat`)

| Method | Endpoint             | Handler               | Description     |
| ------ | -------------------- | --------------------- | --------------- |
| GET    | `/chat/history`      | `GetChatHistory`      | ประวัติแชท      |
| GET    | `/chat/notification` | `GetChatNotification` | การแจ้งเตือนแชท |

#### Stats Routes (`/stats`)

| Method | Endpoint          | Handler           | Description  |
| ------ | ----------------- | ----------------- | ------------ |
| GET    | `/stats/summary`  | `GetStatsSummary` | สรุปสถิติ    |
| GET    | `/stats/device`   | `GetDeviceStats`  | สถิติอุปกรณ์ |
| GET    | `/stats/type`     | `GetTypeStats`    | สถิติประเภท  |
| GET    | `/stats/gen`      | `GenerateStats`   | สร้างสถิติ   |
| GET    | `/stats/generate` | `GenerateStats2`  | สร้างสถิติ   |
| GET    | `/stats/user`     | `GetUserStats`    | สถิติผู้ใช้  |
| GET    | `/stats/case`     | `GetCaseStats`    | สถิติเคส     |

#### Record Routes (`/record`)

| Method | Endpoint          | Handler          | Description       |
| ------ | ----------------- | ---------------- | ----------------- |
| GET    | `/record/request` | `RequestRecord`  | ขอการบันทึก       |
| GET    | `/record/list`    | `GetRecordList`  | รายการบันทึก      |
| GET    | `/record/stopall` | `StopAllRecords` | หยุดบันทึกทั้งหมด |
| GET    | `/record/file`    | `GetFileHistory` | ประวัติไฟล์       |
| GET    | `/record/check`   | `CheckRecord`    | ตรวจสอบการบันทึก  |

#### Notification Routes (`/notification`)

| Method | Endpoint                               | Handler                  | Description     |
| ------ | -------------------------------------- | ------------------------ | --------------- |
| GET    | `/notification/events`                 | `GetNotificationEvents`  | รายการแจ้งเตือน |
| PUT    | `/notification/update/:notificationId` | `UpdateNotification`     | อัพเดทแจ้งเตือน |
| GET    | `/notification/unread`                 | `GetUnreadNotifications` | ยังไม่อ่าน      |
| GET    | `/notification/:notificationId`        | `GetNotificationById`    | ดึงตาม ID       |
| POST   | `/notification`                        | `CreateNotification`     | สร้างแจ้งเตือน  |

#### Radio Routes (`/radio`)

| Method | Endpoint                   | Handler                | Description       |
| ------ | -------------------------- | ---------------------- | ----------------- |
| GET    | `/radio/device`            | `GetRadioDevices`      | รายการอุปกรณ์     |
| GET    | `/radio/device/:id`        | `GetRadioDeviceById`   | ดึงอุปกรณ์ตาม ID  |
| GET    | `/radio/location`          | `GetRadioLocations`    | รายการตำแหน่ง     |
| GET    | `/radio/location/:radioNo` | `GetRadioLocationByNo` | ตำแหน่งตามหมายเลข |

#### Car Routes (`/car`)

| Method | Endpoint    | Handler          | Description  |
| ------ | ----------- | ---------------- | ------------ |
| GET    | `/car/task` | `GetCarTask`     | ดึงข้อมูลงาน |
| PUT    | `/car/task` | `UpdateCarTask`  | อัพเดทงาน    |
| POST   | `/car/task` | `CreateCarTask`  | สร้างงาน     |
| GET    | `/car/list` | `GetCarTaskList` | รายการงาน    |

#### Service Routes (`/service`)

| Method | Endpoint          | Handler         | Description     |
| ------ | ----------------- | --------------- | --------------- |
| GET    | `/service/get`    | `GetService`    | ดึงข้อมูลบริการ |
| PUT    | `/service/update` | `UpdateService` | อัพเดทบริการ    |

#### Upload Routes (`/upload`)

| Method | Endpoint           | Handler         | Description   |
| ------ | ------------------ | --------------- | ------------- |
| POST   | `/upload/file`     | `UploadFile`    | อัพโหลดไฟล์   |
| POST   | `/upload/video`    | `UploadVideo`   | อัพโหลดวิดีโอ |
| GET    | `/upload/list`     | `GetUploadList` | รายการวิดีโอ  |
| POST   | `/upload/sms/send` | `SendSMS`       | ส่ง SMS       |

#### System Routes

| Method | Endpoint      | Handler          | Description      |
| ------ | ------------- | ---------------- | ---------------- |
| GET    | `/`           | `Root`           | Root             |
| GET    | `/health`     | `HealthCheck`    | สถานะ            |
| GET    | `/status`     | `GetStatus`      | สถานะการใช้งาน   |
| GET    | `/service`    | `GetServiceInfo` | ข้อมูลบริการ     |
| POST   | `/webhook`    | `Webhook`        | LiveKit Webhook  |
| POST   | `/sms/custom` | `SendCustomSMS`  | ส่ง SMS          |
| GET    | `/namespace`  | `GetNamespaces`  | รายการ namespace |
| POST   | `/log`        | `AddLog`         | เพิ่มบันทึก      |

#### Test Routes

| Method | Endpoint                    | Handler               | Description   |
| ------ | --------------------------- | --------------------- | ------------- |
| GET    | `/test`                     | `TestUnMuteAll`       | ทดสอบ         |
| GET    | `/test/get/namespace`       | `GetAllNamespaces`    | ดึง namespace |
| GET    | `/test/redis/connection`    | `TestRedisConnection` | ทดสอบ Redis   |
| GET    | `/test/redis/operations`    | `TestRedisOperations` | ทดสอบ Redis   |
| DELETE | `/test/redis/clear`         | `ClearRedisTestData`  | ล้าง Redis    |
| GET    | `/test/mp4/queue`           | `GetMP4Queue`         | คิว MP4       |
| DELETE | `/test/mp4/queue/:recordId` | `RemoveFromQueue`     | ลบจากคิว      |
| DELETE | `/test/mp4/queue`           | `ClearMP4Queue`       | ล้างคิว       |

---

## � Phase 4: Middleware & Routing (Week 4) ✅ COMPLETE

> **สถานะ**: ✅ All middleware and routing implemented

### 4.1 Middleware Implemented

| Middleware        | File                   | Description                              |
| ----------------- | ---------------------- | ---------------------------------------- |
| Auth Middleware   | `auth_middleware.go`   | JWT authentication (required & optional) |
| CORS Middleware   | `cors_middleware.go`   | Cross-Origin Resource Sharing            |
| Logger Middleware | `logger_middleware.go` | Request logging                          |

### 4.2 Router Implementation

- [x] Created `internal/router/router.go` with all route definitions
- [x] Grouped routes by resource (auth, room, user, link, chat, notification, etc.)
- [x] Applied auth middleware to protected routes
- [x] Static file serving for uploads, videos, images, etc.

### 4.3 Main Entry Point

- [x] Updated `cmd/api/main.go` with full integration:
  - Configuration loading
  - Database, Redis, LiveKit initialization
  - All repositories instantiation
  - All services instantiation
  - All handlers instantiation
  - Crontab service initialization
  - Graceful shutdown handling

---

## �🔌 Phase 5: WebSocket/Socket.IO Implementation (Week 5-6) ⏳ PENDING

### 5.1 WebSocket Architecture Decision

#### Option A: ใช้ Socket.IO Go (go-socket.io)

```go
// ข้อดี: Compatible กับ Socket.IO clients เดิม
// ข้อเสีย: Library อาจไม่ stable เท่า native WebSocket
```

#### Option B: ใช้ Fiber WebSocket + Custom Protocol

```go
// ข้อดี: Performance ดี, ควบคุมได้เต็มที่
// ข้อเสีย: ต้อง implement protocol เอง, frontend ต้องแก้
```

**แนะนำ: Option A** - ใช้ `github.com/googollee/go-socket.io` เพื่อให้ compatible กับ frontend เดิม

### 5.2 Socket Namespaces ที่ต้อง Implement

#### 5.2.1 Room Namespace (`/{roomName}`)

```go
// Events:
// - connection: ตรวจสอบ identity และ join room
// - disconnect: cleanup user session
// - chat-message: ส่ง/รับข้อความ
// - position: อัพเดทตำแหน่ง
// - user-connection: แจ้งเตือนผู้ใช้ connect
// - user-disconnect: แจ้งเตือนผู้ใช้ disconnect
// - room-record: สถานะการบันทึก
// - agentCar: ตำแหน่งรถ
// - conference-status: สถานะ conference
```

#### 5.2.2 Mobile Namespace (`/mobile`)

```go
// Events:
// - connection: validate task id
// - disconnect: cleanup
// - location: รับตำแหน่งจาก mobile app
// - status: อัพเดทสถานะ
// - message: ส่งข้อความกลับ
```

#### 5.2.3 Notification Namespace (`/notification`)

```go
// Events:
// - connection: ไม่ต้อง auth
// - all: ส่งทุก notifications
// - unread: ส่งที่ยังไม่อ่าน
// - new: notification ใหม่
// - update: อัพเดท notification
// - read: mark as read
```

#### 5.2.4 Queue Namespace (`/queue`, `/newqueue`)

```go
// Events:
// - connection: join queue room
// - queue-update: อัพเดท queue
// - newcase: เคสใหม่
```

### 5.3 Redis Adapter for Socket.IO

```go
// ใช้ Redis pub/sub สำหรับ cross-instance communication
// - REDIS_ADAPTER_DB: 1 (Socket.IO adapter)
// - REDIS_STATE_DB: 2 (Socket.IO state)
```

### 5.4 Cross-Instance Event Manager

```go
// internal/socket/event_manager.go

type CrossInstanceEventManager interface {
    Initialize(ctx context.Context) error
    PublishEvent(ctx context.Context, namespace, event string, data interface{}) error
    RegisterHandler(namespace string, handler EventHandler)
    UnregisterHandler(namespace string)
    Cleanup() error
}

// Events:
// - user_disconnect
// - user_connect
// - car_position
// - chat_message
// - room_record
```

---

## 🎯 Phase 6: LiveKit Webhook Handler (Week 6) ⏳ PENDING

> **สถานะ**: ⚠️ Basic webhook handler implemented, needs enhancement

### 6.1 Webhook Events ที่ต้อง Handle

```go
// internal/handler/webhook_handler.go

type WebhookHandler struct {
    roomService     RoomService
    recordService   RecordService
    socketManager   SocketManager
    trackDataCache  sync.Map
    participantTimers sync.Map
}

// Events:
// - room_started
// - room_finished
// - participant_joined
// - participant_left
// - track_published
// - track_subscribed
// - egress_started
// - egress_ended
```

### 6.2 Auto Recording Logic

```go
// เมื่อ participant_joined และ autoRecord=1
// 1. ตรวจสอบ recordStatus
// 2. ถ้า recordType = RoomCompositeVideoAudio หรือ RoomCompositeAudio
//    → เริ่ม startRoomCompositeEgress
// 3. ถ้า recordType = TrackComposite
//    → รอ track_published ทั้ง video และ audio
//    → แล้วเริ่ม startTrackCompositeEgress
```

---

## 🛡️ Phase 7: Testing & Integration (Week 7) ⏳ PENDING

> **หมายเหตุ**: Phase 6 Middleware ได้ implement แล้วใน Phase 4

### 7.1 Unit Tests (Examples from original Phase 6 code)

```go
// internal/middleware/auth_middleware.go

func AuthMiddleware(authService AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "status": "FAIL",
                "message": "Invalid token",
            })
        }

        claims, err := authService.VerifyToken(c.Context(), token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "status": "FAIL",
                "data": err.Error(),
            })
        }

        c.Locals("user", claims)
        return c.Next()
    }
}
```

### 6.2 Join Conference Middleware

```go
// internal/middleware/join_conference_middleware.go

func JoinConferenceMiddleware(linkService LinkService, authService AuthService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        linkID := c.FormValue("linkID")
        if linkID != "" {
            linkDetail, err := linkService.GetLinkDetail(c.Context(), linkID, "", nil)
            if err != nil || linkDetail == nil {
                return c.Status(400).JSON(fiber.Map{"status": "FAIL"})
            }
            c.Locals("room", linkDetail.Room)
            return c.Next()
        }

        // Verify token
        token := c.Get("Authorization")
        claims, err := authService.VerifyToken(c.Context(), token)
        if err != nil {
            c.Locals("decoded", nil)
        } else {
            c.Locals("decoded", claims)
        }
        return c.Next()
    }
}
```

---

## 📊 Phase 8: Static Files & File Upload (Week 7) ✅ COMPLETE

> **สถานะ**: ✅ Implemented in upload_handler.go and router.go

### 8.1 Static File Serving

```go
// Fiber static file serving
app.Static("/logo", "./logo")
app.Static("/videos", "./uploads/videos")
app.Static("/images", "./uploads/images")
app.Static("/thumbnails", "./uploads/thumbnails")
app.Static("/files", "./uploads/files")
app.Static("/record", "./record-file")
```

### 8.2 File Upload Handler

```go
// internal/handler/upload_handler.go

type UploadHandler struct {
    fileService FileService
    maxSize     int64
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
    // Multiple files: myFiles[]
    form, err := c.MultipartForm()
    files := form.File["myFiles[]"]
    // ...
}

func (h *UploadHandler) UploadVideo(c *fiber.Ctx) error {
    // Single video: myVideo
    file, err := c.FormFile("myVideo")
    // ...
}
```

---

## 🔄 Phase 9: Cron Jobs (Week 7) ✅ COMPLETE

> **สถานะ**: ✅ Implemented in crontab_service.go

### 9.1 Cron Job Implementation

```go
// internal/service/crontab_service.go

func (s *CrontabServiceImpl) InitCronJobs() error {
    c := cron.New(cron.WithLocation(time.FixedZone("ICT", 7*60*60)))

    // Room cleanup - every 30 minutes
    c.AddFunc("*/30 * * * *", func() {
        s.roomService.AutoRoomExpiredClose(context.Background())
    })

    // Link cleanup - every 30 minutes
    c.AddFunc("*/30 * * * *", func() {
        s.linkService.AutoLinkExpiredClose(context.Background())
    })

    // LiveKit health check - every 10 minutes
    c.AddFunc("*/10 * * * *", func() {
        s.healthCheckLiveKit(context.Background())
    })

    c.Start()
    return nil
}
```

---

## 🚦 Phase 10: Graceful Shutdown (Week 8) ✅ COMPLETE

> **สถานะ**: ✅ Implemented in cmd/api/main.go

### 10.1 Shutdown Handler

```go
// cmd/api/main.go

func gracefulShutdown(server *fiber.App, services ...Cleanable) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    log.Println("Starting graceful shutdown...")

    // Phase 1: Stop cron jobs
    crontabService.Cleanup()

    // Phase 2: Stop socket connections
    socketManager.Cleanup()

    // Phase 3: Close Redis connections
    redisManager.CloseAll()

    // Phase 4: Close database pool
    db.Close()

    // Phase 5: Shutdown HTTP server
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.ShutdownWithContext(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
}
```

---

## 📝 Phase 11: Testing & Documentation (Week 8-9) ⏳ PENDING

### 11.1 Unit Tests

- [ ] Repository tests
- [ ] Service tests
- [ ] Handler tests
- [ ] Middleware tests

### 11.2 Integration Tests

- [ ] API endpoint tests
- [ ] WebSocket tests
- [ ] Database tests
- [ ] Redis tests

### 11.3 Documentation

- [ ] API documentation (Swagger/OpenAPI)
- [ ] WebSocket protocol documentation
- [ ] Deployment guide
- [ ] Environment variables documentation

---

## 🔧 Environment Variables

```env
# Server
PORT=5500
ENVIRONMENT=development

# Database
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=
MYSQL_DATABASE=conference

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
REDIS_ADAPTER_DB=1
REDIS_STATE_DB=2

# API
API_URL=http://localhost:5500

# SMS
SMS_ENABLE=true
SMS_API_URL=http://portal-api-idems.niems.go.th/inet/sms

# File
RECORD_PATH=./record-file
FILE_SIZE_LIMIT=524288000

# Room
JOIN_ROOM_REPEAT_DELAY=5000
AUTO_CLOSE_ROOM=true
ROOM_DAY_DEFAULT_TIMEOUT=24H

# LiveKit
EGRESS_LIMIT=4

# Radio API
RADIO_LOCATION_API_URL=
RADIO_LOCATION_API_CREDENTIALS_USERNAME=
RADIO_LOCATION_API_CREDENTIALS_PASSWORD=

# Encode API
ENCODE_API=http://encode-api:5600

# Custom
CUSTOM_CHARSET=ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz
```

---

## ⚠️ Critical Points

### 1. Socket.IO Compatibility

- ต้อง test กับ frontend เดิมให้แน่ใจว่า compatible
- อาจต้องใช้ specific version ของ go-socket.io
- Redis adapter configuration ต้องเหมือนกับ Node.js

### 2. LiveKit SDK

- ใช้ `livekit/server-sdk-go/v2`
- Token generation ต้องเหมือนกับ Node.js SDK
- Egress/Recording ต้อง test กับ production environment

### 3. Database Transactions

- Node.js ไม่ได้ใช้ transactions ชัดเจน
- Golang ควรเพิ่ม transaction support สำหรับ operations ที่ซับซ้อน

### 4. Error Handling

- Node.js มี try-catch แต่หลายที่ไม่ส่ง error กลับ
- Golang ต้อง handle errors อย่างถูกต้อง

### 5. Timezone

- ใช้ `Asia/Bangkok` สำหรับ cron jobs
- Database ใช้ `YYYY-MM-DD HH:mm:ss` format

---

## 📅 Timeline Summary

| Phase    | Duration | Description       |
| -------- | -------- | ----------------- |
| Phase 1  | Week 1   | Foundation Setup  |
| Phase 2  | Week 2-3 | Core Services     |
| Phase 3  | Week 3-4 | HTTP Handlers     |
| Phase 4  | Week 4-5 | WebSocket         |
| Phase 5  | Week 5   | Webhook Handler   |
| Phase 6  | Week 5   | Middleware        |
| Phase 7  | Week 6   | Static Files      |
| Phase 8  | Week 6   | Cron Jobs         |
| Phase 9  | Week 7   | Graceful Shutdown |
| Phase 10 | Week 7-8 | Testing           |

**Total Estimated Time: 6-8 Weeks**

---

## ✅ Checklist

### Foundation

- [ ] Config loader
- [ ] Database connection
- [ ] Redis connection
- [ ] Logger setup
- [ ] Models (มีแล้ว ✅)

### Repositories

- [ ] Room repository
- [ ] User repository
- [ ] Link repository
- [ ] Chat repository
- [ ] Notification repository
- [ ] Record repository
- [ ] Car repository
- [ ] Case repository
- [ ] Radio repository
- [ ] Stats repository
- [ ] Usage log repository

### Services

- [ ] Auth service
- [ ] Room service
- [ ] User service
- [ ] Link service
- [ ] Chat service
- [ ] Notification service
- [ ] Record service
- [ ] Car service
- [ ] Case service
- [ ] Radio service
- [ ] Stats service
- [ ] LiveKit service
- [ ] SMS service
- [ ] File service
- [ ] Crontab service

### Handlers

- [ ] Auth handler
- [ ] Room handler
- [ ] User handler
- [ ] Link handler
- [ ] Chat handler
- [ ] Notification handler
- [ ] Record handler
- [ ] Car handler
- [ ] Case handler
- [ ] Radio handler
- [ ] Stats handler
- [ ] Upload handler
- [ ] Webhook handler
- [ ] Service handler
- [ ] Test handler

### Socket

- [ ] Socket.IO integration
- [ ] Room namespace
- [ ] Mobile namespace
- [ ] Notification namespace
- [ ] Queue namespace
- [ ] Redis adapter
- [ ] Cross-instance events

### Middleware

- [ ] Auth middleware
- [ ] Join conference middleware
- [ ] CORS middleware
- [ ] Logger middleware

### Infrastructure

- [ ] Cron jobs
- [ ] Graceful shutdown
- [ ] Static files
- [ ] File upload
- [ ] Health check

### Testing

- [ ] Unit tests
- [ ] Integration tests
- [ ] E2E tests

### Deployment

- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Documentation
