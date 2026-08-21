-- migrate:up
ALTER TABLE members
  ADD COLUMN faculty VARCHAR(100),
  ADD COLUMN study_program VARCHAR(150),
  ADD COLUMN cohort_year INT;

-- migrate:down
ALTER TABLE members
  DROP COLUMN faculty,
  DROP COLUMN study_program,
  DROP COLUMN cohort_year;
