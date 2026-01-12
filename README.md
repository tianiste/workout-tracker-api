# Workout tracker API
A rest api for tracking workouts.
Built with go, gin and sqlite, trying to follow a clean repo -> service -> handler architecture, for maintainability, testing...

## Features
- User authentication with JWT (Currently no refresh tokens yet)
- Create/Update/Delete workouts
- Attach exercises to workouts

## Architecture

- Repo → raw DB access (SQL only)
- Service → business logic
- Handler → HTTP / JSON
- Middleware → auth, rate limiting, etc. (rate limiting not yet implemented)

## Setup 
1.
```json
 git clone https://github.com/tianiste/workout-tracker-api
```
3. Create a .env with a 256bit JWT key token

4. Run the server
```bash
go run main.go
```

## Authentication 
```bash
 Authorization: Bearer <token>
```

## Rate Limiting
Added a rate limiter that limits requests to 5/second

## Endpoints
### Health
GET /api/ping  
Body: none

### Authentication (Public)

POST /api/register  
Body:
```json
{
  "name": "string",
  "password": "string"
}
```

POST /api/login  
Body:
```json
{
  "name": "string",
  "password": "string"
}
```

### Workouts (Protected)

POST /api/workouts  
Body:
```json
{
  "performedAt": "RFC3339 timestamp string",
  "durationMinutes": 45,
  "notes": "optional string"
}
```

GET /api/workouts  

GET /api/workouts/:id  

GET /api/workouts/:id/details  

PUT /api/workouts/:id  
Body:
```json
{
  "performedAt": "RFC3339 timestamp string",
  "durationMinutes": 50,
  "notes": "optional string"
}
```

DELETE /api/workouts/:id  

GET /api/workouts/:id/report

Sample body:
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

POST /api/workouts/:id/exercises  
Body:
```json
{
  "exerciseId": 1,
  "exerciseOrder": 1,
  "notes": "optional string"
}
```

PUT /api/workout-exercises/:id  
Body:
```json
{
  "exerciseOrder": 2,
  "notes": "optional string"
}
```

DELETE /api/workout-exercises/:id  

### Sets (Protected)

POST /api/workout-exercises/:id/sets  
Body:
```json
{
  "setNumber": 1,
  "reps": 10,
  "weight": 60.0
}
```

PUT /api/sets/:id  
Body:
```json
{
  "reps": 12,
  "weight": 62.5
}
```

DELETE /api/sets/:id  

### Exercises (Protected)

GET /api/exercises  

## Todo
1. Workout scheduling
2. [x] Report generation
3. Maybe asymmetric JWT and refresh tokens
4. [x] Rate limiting
5. Endpoint for adding custom workouts that are exclusive to the user



[Project idea](https://roadmap.sh/projects/fitness-workout-tracker)
