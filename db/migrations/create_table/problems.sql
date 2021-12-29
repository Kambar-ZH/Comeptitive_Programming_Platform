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