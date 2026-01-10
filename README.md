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
1. `bash
git clone https://github.com/tianiste/workout-tracker-api`
2. Create a .env with a 256bit JWT key token

3. Run the server `bash
go run main.go`

## Authentication 
`bash
Authorization: Bearer <token>`

## Endpoints

### Health
GET /api/ping

### Authentication (Public)
POST /api/register  
POST /api/login  

### Workouts (Protected)
POST /api/workouts  
GET /api/workouts  
GET /api/workouts/:id  
GET /api/workouts/:id/details  
PUT /api/workouts/:id  
DELETE /api/workouts/:id  

### Workout Exercises (Protected)
POST /api/workouts/:id/exercises  
PUT /api/workout-exercises/:id  
DELETE /api/workout-exercises/:id  

### Sets (Protected)
POST /api/workout-exercises/:id/sets  
PUT /api/sets/:id  
DELETE /api/sets/:id  

### Exercises (Protected)
GET /api/exercises


## Todo
1. Workout scheduling
2. Report generation
3. Maybe asymmetric JWT and refresh tokens
4. Rate limiting
5. Endpoint for adding custom workouts that are exclusive to the user
