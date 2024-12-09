CREATE TABLE IF NOT EXISTS books (
     id VARCHAR(36) PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     author VARCHAR(255) NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for sorting and filtering
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);

-- Insert some sample data
INSERT INTO books (id, title, author) VALUES
      ('1', 'The Go Programming Language', 'Alan A. A. Donovan'),
      ('2', 'Clean Code', 'Robert C. Martin'),
      ('3', 'Design Patterns', 'Erich Gamma'),
      ('4', 'Refactoring', 'Martin Fowler'),
      ('5', 'You Don’t Know JS', 'Kyle Simpson'),
      ('6', 'The Pragmatic Programmer', 'Andrew Hunt, David Thomas'),
      ('7', 'Introduction to Algorithms', 'Thomas H. Cormen'),
      ('8', 'JavaScript: The Good Parts', 'Douglas Crockford'),
      ('9', 'Domain-Driven Design', 'Eric Evans'),
      ('10', 'Effective Java', 'Joshua Bloch'),
      ('11', 'Programming Pearls', 'Jon Bentley'),
      ('12', 'Code Complete', 'Steve McConnell'),
      ('13', 'Structure and Interpretation of Computer Programs', 'Harold Abelson, Gerald Jay Sussman'),
      ('14', 'Head First Design Patterns', 'Eric Freeman, Elisabeth Robson'),
      ('15', 'Artificial Intelligence: A Modern Approach', 'Stuart Russell, Peter Norvig'),
      ('16', 'Algorithms to Live By', 'Brian Christian, Tom Griffiths'),
      ('17', 'Deep Learning', 'Ian Goodfellow, Yoshua Bengio, Aaron Courville'),
      ('18', 'Computer Networking: A Top-Down Approach', 'James Kurose, Keith Ross'),
      ('19', 'Python Crash Course', 'Eric Matthes'),
      ('20', 'Fluent Python', 'Luciano Ramalho')
ON CONFLICT (id) DO NOTHING;