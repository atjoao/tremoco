CREATE TABLE IF NOT EXISTS Playlists(
    id INT AUTO_INCREMENT,
    userId INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(userId) REFERENCES Users(id),

    PRIMARY KEY (id),
);

CREATE TABLE IF NOT EXISTS Music(
    id VARCHAR(30) NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    PRIMARY KEY(id),
);

CREATE TABLE IF NOT EXISTS Playlist_Music(
    playlist_id INT,
    music_id VARCHAR(30) NOT NULL,
    FOREIGN KEY (playlist_id) REFERENCES Playlists(id),
    FOREIGN KEY (music_id) REFERENCES Music(id),
    PRIMARY KEY (playlist_id, music_id)
)
-- nao guardar stream links pois espira a cada 6 horas