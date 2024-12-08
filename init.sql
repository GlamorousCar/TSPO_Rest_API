CREATE TABLE IF NOT EXISTS books (id VARCHAR(36) PRIMARY KEY,title VARCHAR(255) NOT NULL,author VARCHAR(255) NOT NULL);

-- Insert some sample data
INSERT INTO books (id, title, author) VALUES
      ('1', 'The Go Programming Language', 'Alan A. A. Donovan'),
      ('2', 'Clean Code', 'Robert C. Martin'),
      ('3', 'Design Patterns', 'Erich Gamma');