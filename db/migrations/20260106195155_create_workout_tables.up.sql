PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS categories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS muscle_groups (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS exercises (
  id INTEGER PRIMARY KEY AUTOINCREMENT,

  owner_user_id INTEGER, 
  name TEXT NOT NULL,

  category_id INTEGER NOT NULL,
  muscle_group_id INTEGER,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),

  FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  FOREIGN KEY (muscle_group_id) REFERENCES muscle_groups(id),

  UNIQUE (owner_user_id, name)
);

CREATE TABLE IF NOT EXISTS workouts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,

  performed_at TEXT NOT NULL,
  duration_minutes INTEGER,
  notes TEXT,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),

  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workout_exercises (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  workout_id INTEGER NOT NULL,
  exercise_id INTEGER NOT NULL,
  exercise_order INTEGER NOT NULL,
  notes TEXT,

  FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
  FOREIGN KEY (exercise_id) REFERENCES exercises(id)
);

CREATE TABLE IF NOT EXISTS sets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  workout_exercise_id INTEGER NOT NULL,
  set_number INTEGER NOT NULL,

  reps INTEGER,
  weight REAL,

  FOREIGN KEY (workout_exercise_id) REFERENCES workout_exercises(id) ON DELETE CASCADE,
  UNIQUE (workout_exercise_id, set_number)
);

CREATE INDEX IF NOT EXISTS idx_workouts_user_date
  ON workouts(user_id, performed_at);

CREATE INDEX IF NOT EXISTS idx_exercises_owner
  ON exercises(owner_user_id);

CREATE INDEX IF NOT EXISTS idx_exercises_category
  ON exercises(category_id);

CREATE INDEX IF NOT EXISTS idx_workout_exercises_workout
  ON workout_exercises(workout_id);

CREATE INDEX IF NOT EXISTS idx_sets_workout_exercise
  ON sets(workout_exercise_id);

