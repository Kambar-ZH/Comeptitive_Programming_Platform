CREATE TABLE validators (
    problem_id INTEGER,
    solution_file VARCHAR(255) NOT NULL,
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);