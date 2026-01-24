# Workout tracker API
A rest api for tracking workouts.
Built with go, gin and sqlite, trying to follow a clean repo -> service -> handler architecture, for maintainability, testing...

## Features
- User authentication with JWT access tokens
- Refresh token authentication (rotating refresh tokens stored hashed in DB)
- Create / Update / Delete workouts
- Attach exercises to workouts
- Workout reports and statistics
- Rate limiting (5 requests / second)


## Architecture

- Repo → raw DB access (SQL only)
- Service → business logic (auth, refresh, rotation)
- Handler → HTTP / JSON (Gin)
- Middleware → auth, rate limiting


## Setup

1. Clone the repository
```bash
git clone https://github.com/tianiste/workout-tracker-api
cd workout-tracker-api
```

2. Create a `.env` file with a JWT secret key (256-bit recommended)
```env
JWT_KEY=your-secret-key
```

3. Run the server
```bash
go run main.go
```


## Authentication

### Access Token
Sent on protected routes using:
```http
Authorization: Bearer <access_token>
```

- Short-lived (10 minutes)
- Used for all protected API routes

### Refresh Token
- Stored in an HttpOnly cookie
- Rotated on every `/refresh`
- Stored hashed (SHA-256) in the database
- Revoked on logout

Cookie settings:
- `HttpOnly`
- `SameSite=Lax`
- `Secure=true` (production)
- `Path=/`


## Authentication Flow

1. Login
   - Returns access token
   - Sets refresh token cookie

2. Refresh
   - Validates refresh token from cookie
   - Issues new access token
   - Rotates refresh token

3. Logout
   - Revokes refresh token in DB
   - Clears cookie


## Endpoints

### Health
```
GET /api/ping
```


### Authentication (Public)

#### Register
```
POST /api/register
```
```json
{
  "name": "string",
  "password": "string"
}
```

#### Login
```
POST /api/login
```
```json
{
  "name": "string",
  "password": "string"
}
```

#### Refresh
```
POST /refresh
```
- Uses refresh token cookie
- Returns new access token

#### Logout
```
POST /logout
```
- Revokes refresh token
- Clears cookie


### Workouts (Protected)

#### Create workout
```
POST /api/workouts
```
```json
{
  "performedAt": "RFC3339 timestamp string",
  "durationMinutes": 45,
  "notes": "optional string"
}
```

#### Get all workouts
```
GET /api/workouts
```

#### Get workout by ID
```
GET /api/workouts/:id
```

#### Get workout details
```
GET /api/workouts/:id/details
```

#### Update workout
```
PUT /api/workouts/:id
```
```json
{
  "performedAt": "RFC3339 timestamp string",
  "durationMinutes": 50,
  "notes": "optional string"
}
```

#### Delete workout
```
DELETE /api/workouts/:id
```

#### Workout report
```
GET /api/workouts/:id/report
```

Sample response:
```json
{
  "workoutId": 1,
  "userId": 1,
  "performedAt": "1",
  "notes": "none",
  "createdAt": "2026-01-09T03:13:28Z",
  "totalExercises": 1,
  "totalSets": 1,
  "totalReps": 10,
  "totalVolume": 600,
  "exercises": [
    {
      "exerciseId": 1,
      "exerciseName": "Bench Press",
      "setsCount": 1,
      "totalReps": 10,
      "maxWeight": 60,
      "totalVolume": 600
    }
  ]
}
```


### Workout Exercises (Protected)

```
POST /api/workouts/:id/exercises
```
```json
{
  "exerciseId": 1,
  "exerciseOrder": 1,
  "notes": "optional string"
}
```

```
PUT /api/workout-exercises/:id
```
```json
{
  "exerciseOrder": 2,
  "notes": "optional string"
}
```

```
DELETE /api/workout-exercises/:id
```


### Sets (Protected)

```
POST /api/workout-exercises/:id/sets
```
```json
{
  "setNumber": 1,
  "reps": 10,
  "weight": 60.0
}
```

```
PUT /api/sets/:id
```
```json
{
  "reps": 12,
  "weight": 62.5
}
```

```
DELETE /api/sets/:id
```


### Exercises (Protected)

```
GET /api/exercises
```


## Todo
1. Workout scheduling
2. Custom user-specific exercises
3. Session/device management UI
4. Token reuse detection alerts
5. Production Docker setup


[Project idea](https://roadmap.sh/projects/fitness-workout-tracker)

