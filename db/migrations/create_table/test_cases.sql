CREATE TABLE test_cases (
    id BIGSERIAL PRIMARY KEY,
    problem_id INTEGER,
    as_sample BOOLEAN DEFAULT FALSE,
    input_file_path VARCHAR(255) NOT NULL,
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);