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