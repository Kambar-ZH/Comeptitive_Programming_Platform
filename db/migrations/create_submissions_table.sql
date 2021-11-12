CREATE TABLE submissions (
    id BIGSERIAL,
    contest_id INTEGER,
    author_handle VARCHAR(255),
    submission_date DATE,
    problem_id INTEGER,
    verdict VARCHAR(255),
    failed_test INTEGER
);