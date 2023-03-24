-- CREATE TABLE tasks (
--     id INT AUTO_INCREMENT NOT NULL,
--     task_name VARCHAR(50) NOT NULL,
--     completed BOOLEAN,
--     owner_id INT,
--     PRIMARY KEY (id),
--     FOREIGN KEY (owner_id) REFERENCES owners(owner_id) ON DELETE CASCADE
-- );

-- CREATE TABLE owners (
--     owner_id INT AUTO_INCREMENT NOT NULL,
--     owner_name VARCHAR(50),
--     PRIMARY KEY (owner_id)
-- );

-- INSERT INTO owners (owner_name)
-- VALUES ("Sean The Programmer");

INSERT INTO tasks (task_name, completed, owner_id)
VALUES 
("complete this project", false, 1);