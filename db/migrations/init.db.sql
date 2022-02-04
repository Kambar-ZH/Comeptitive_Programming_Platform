CREATE TABLE contests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) DEFAULT '',
    start_date DATE DEFAULT NOW(),
    end_date DATE DEFAULT NOW(),
    description VARCHAR(255) DEFAULT ''
);

CREATE TABLE problems (
    id BIGSERIAL PRIMARY KEY,
    contest_id INTEGER,
    index VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    statement VARCHAR(255),
    input VARCHAR(255),
    output VARCHAR(255),
    difficulty INTEGER,
    CONSTRAINT fk_contests
      FOREIGN KEY(contest_id)
	      REFERENCES contests(id)
);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE problems_tags (
    problem_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    CONSTRAINT pk_problems_tags
        PRIMARY KEY(problem_id, tag_id),
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id),
    CONSTRAINT fk_tags
      FOREIGN KEY(tag_id)
	      REFERENCES tags(id)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    handle VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    country VARCHAR(255) DEFAULT '',
    city VARCHAR(255) DEFAULT '',
    rating INTEGER DEFAULT 0,
    max_rating INTEGER DEFAULT 0,
    avatar VARCHAR(255) DEFAULT '',
    encrypted_password VARCHAR(255) NOT NULL
);

CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    contest_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    submission_date DATE DEFAULT NOW(),
    problem_id INTEGER NOT NULL,
    verdict VARCHAR(255) NOT NULL,
    failed_test INTEGER DEFAULT 0,
    CONSTRAINT fk_users
      FOREIGN KEY(user_id)
	      REFERENCES users(id),
    CONSTRAINT fk_contests
      FOREIGN KEY(contest_id)
	      REFERENCES contests(id),
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);

CREATE TABLE test_cases (
    id BIGSERIAL PRIMARY KEY,
    problem_id INTEGER,
    as_sample BOOLEAN DEFAULT FALSE,
    test_file VARCHAR(255) NOT NULL,
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);

CREATE TABLE validators (
    problem_id INTEGER,
    solution_file VARCHAR(255) NOT NULL,
    CONSTRAINT fk_problems
      FOREIGN KEY(problem_id)
	      REFERENCES problems(id)
);

INSERT INTO contests (
    name
) VALUES ('Div3 755');

INSERT INTO problems (
    contest_id,
    index,
    name,
    statement,
    input,
    output,
    difficulty
) VALUES (1, 'A', 'Sum on the array', 'You given an array find the sum.', 'On the first line given N - the size of the array. On the next line there are N integers seperated be space.', 'Print the answer.', 600),
    (1, 'B', 'Fibonacci', 'Your task is to find the ith fibonacci number.', 'The single line containts integer i.', 'Output the ith fibonacci number.', 800);

INSERT INTO tags (
    name
) VALUES ('dp'), 
    ('math'), 
    ('constructive'), 
    ('binary search'), 
    ('sort');

INSERT INTO problems_tags (
    problem_id,
    tag_id
) VALUES (1, 3), 
    (2, 2);

INSERT INTO test_cases (
    problem_id,
    test_file
) VALUES (1, '../../test/problems/0001/tests/1.txt'),
    (1, '../../test/problems/0001/tests/2.txt'),
    (2, '../../test/problems/0002/tests/1.txt'),
    (2, '../../test/problems/0002/tests/2.txt');

INSERT INTO validators (
    problem_id,
    solution_file
) VALUES (1, '../../test/problems/0001/solution.go'),
    (2, '../../test/problems/0002/solution.go');