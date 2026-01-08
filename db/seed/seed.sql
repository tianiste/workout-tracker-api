PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

INSERT OR IGNORE INTO categories (name) VALUES
  ('strength'),
  ('cardio'),
  ('mobility'),
  ('stretching'),
  ('sport');

INSERT OR IGNORE INTO muscle_groups (name) VALUES
  ('chest'),
  ('back'),
  ('legs'),
  ('shoulders'),
  ('arms'),
  ('core'),
  ('glutes'),
  ('full_body');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT
  NULL,
  'Bench Press',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'chest');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Incline Dumbbell Press',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'chest');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Push-Up',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'chest');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Deadlift',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'back');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Barbell Row',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'back');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Pull-Up',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'back');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Lat Pulldown',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'back');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Back Squat',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Front Squat',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Romanian Deadlift',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'glutes');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Lunge',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Leg Press',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Hip Thrust',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'glutes');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Overhead Press',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'shoulders');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Lateral Raise',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'shoulders');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Face Pull',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'shoulders');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Bicep Curl',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'arms');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Triceps Pushdown',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'arms');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Plank',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'core');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Hanging Leg Raise',
  (SELECT id FROM categories WHERE name = 'strength'),
  (SELECT id FROM muscle_groups WHERE name = 'core');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Running',
  (SELECT id FROM categories WHERE name = 'cardio'),
  (SELECT id FROM muscle_groups WHERE name = 'full_body');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Cycling',
  (SELECT id FROM categories WHERE name = 'cardio'),
  (SELECT id FROM muscle_groups WHERE name = 'full_body');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Rowing Machine',
  (SELECT id FROM categories WHERE name = 'cardio'),
  (SELECT id FROM muscle_groups WHERE name = 'full_body');

-- Mobility / Stretching
INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Hip Mobility Flow',
  (SELECT id FROM categories WHERE name = 'mobility'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Shoulder Mobility Flow',
  (SELECT id FROM categories WHERE name = 'mobility'),
  (SELECT id FROM muscle_groups WHERE name = 'shoulders');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Hamstring Stretch',
  (SELECT id FROM categories WHERE name = 'stretching'),
  (SELECT id FROM muscle_groups WHERE name = 'legs');

INSERT OR IGNORE INTO exercises (owner_user_id, name, category_id, muscle_group_id)
SELECT NULL, 'Chest Stretch',
  (SELECT id FROM categories WHERE name = 'stretching'),
  (SELECT id FROM muscle_groups WHERE name = 'chest');


COMMIT;

