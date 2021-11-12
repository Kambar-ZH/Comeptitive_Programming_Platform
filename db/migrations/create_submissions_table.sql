CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    contest_id INTEGER NOT NULL,
    author_handle VARCHAR(255) NOT NULL,
    submission_date DATE DEFAULT NOW(),
    problem_id INTEGER DEFAULT 0,
    verdict VARCHAR(255) NOT NULL,
    failed_test INTEGER DEFAULT 0,
    CONSTRAINT fk_users
      FOREIGN KEY(author_handle)
	      REFERENCES users(handle)
);