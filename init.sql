CREATE TABLE IF NOT EXISTS fictions (
    ID                  SERIAL PRIMARY KEY,
    Contributor_Name    VARCHAR(255) NOT NULL,
    Cover               TEXT,
    Title               VARCHAR(255) NOT NULL,
    Subtitle            VARCHAR(255),
    Owner               VARCHAR(255),
    Author              VARCHAR(255),
    Artist              VARCHAR(255),
    Status              VARCHAR(50) CHECK (Status IN ('Ongoing', 'Completed', 'Hiatus')),
    Synopsis            TEXT,
    Created             DATE DEFAULT CURRENT_DATE
);

