CREATE TABLE IF NOT EXISTS Playlists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INT NOT NULL,
    name TEXT NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    image TEXT DEFAULT NULL,
    FOREIGN KEY(userId) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Music(
    id TEXT NOT NULL, -- local_<id>
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    duration INT NOT NULL,
    genre TEXT NOT NULL,
    location TEXT NOT NULL,
    PRIMARY KEY(id),
    UNIQUE(id, location)
);
-- album detected by folder name
CREATE TABLE IF NOT EXISTS Album(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cover TEXT,
    UNIQUE(name, cover)
);

CREATE TABLE IF NOT EXISTS Album_Music (
    album_id INT,
    music_id TEXT,
    FOREIGN KEY (album_id) REFERENCES Album(id),
    FOREIGN KEY (music_id) REFERENCES Music(id),
    PRIMARY KEY (album_id, music_id),
    UNIQUE(music_id)
);

CREATE TABLE IF NOT EXISTS Playlist_Music(
    playlist_id INT,
    music_id TEXT NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES Playlists(id),
    PRIMARY KEY (playlist_id, music_id)
);
