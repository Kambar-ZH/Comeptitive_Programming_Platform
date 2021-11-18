CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    contest_id INTEGER NOT NULL,
    author_handle VARCHAR(255) NOT NULL,
    submission_date DATE DEFAULT NOW(),
    problem_id INTEGER NOT NULL,
    verdict VARCHAR(255) NOT NULL,
    failed_test INTEGER DEFAULT 0,
    CONSTRAINT fk_users
      FOREIGN KEY(author_handle)
	      REFERENCES users(handle),
    CONSTRAINT fk_contests
      FOREIGN KEY(contest_id)
	      REFERENCES contests(id),
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);