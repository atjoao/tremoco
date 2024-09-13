const playerImage = document.getElementById("player-image")
const playerTitle = document.getElementById("player-title") 
const playerAuthor = document.getElementById("player-author") 

let todo = null;

const queue = {
    position: 0,
    currentSong: null,

    originalQueue: [], // this queue should not be modified in case of shuffle off
    currentQueue: [], // queue that will be listened for changes
    alteradyPlayed: []
}

const audio = new Audio();

audio.onplay = () => {
    let thumb;

    if (!queue.currentSong.thumbnails[0].url.includes("https://")){
        thumb = new URL(window.location)
        thumb.pathname = queue.currentSong.thumbnails[0].url
    } else {
        thumb = queue.currentSong.thumbnails[0].url
    }
    
    navigator.mediaSession.metadata = new MediaMetadata({
        title: queue.currentSong.title,
        artist: queue.currentSong.author,
        artwork: [
            {
                src: thumb.toString(),
                sizes: "96x96",
                type: 'image/png'
            }
        ]
    })

    playerImage.src = thumb;
    playerTitle.textContent = queue.currentSong.title;
    playerAuthor.textContent = queue.currentSong.author;

}

audio.ontimeupdate = (e) => {
    // calculate if 0 - 100% music
    // set width of the player slider watever
};

audio.onended = () => {
    nextTrack();
    // next track
    // if false clear div content
    // if true start next audio
    // if looping return;
};

function sendToQueue(musicList, reset){
    // add playlist to the set
    // if array check array correct and add
    // if not add to array and start bnluh bluh
    console.log(musicList, reset)
    if (reset){
        audio.pause()
        audio.src = null;

        queue.currentQueue = []
        queue.originalQueue = []
        queue.position = 0;
    }

    queue.currentQueue.push(...musicList)
    queue.originalQueue.push(...musicList)

    if (reset) startQueue();

}

function startQueue(){
    let streams = queue.currentQueue[queue.position].streams

    if (streams.length == 1){
        audio.src = queue.currentQueue[queue.position].streams[0].streamUrl
    } else {
        for (let index = 0; index < streams.length; index++) {
            const element = streams[index];
            if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"opus\"")){
                audio.src = element.streamUrl
            } else if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"mp4a.40.5\"")){
                audio.src = element.streamUrl
            } else {
                audio.src = streams[0]
            }
        }
    }

    queue.currentSong = queue.currentQueue[queue.position]

    audio.play()

}

function clearQueue(){
    // stops audio
    // clears queue
}

function nextTrack(){
    if (!queue.currentQueue[queue.position + 1]) {
        console.log("Next track is undefined, stopping execution.");

        navigator.mediaSession.metadata = null;

        return;
    }

    queue.alteradyPlayed.push(queue.currentQueue[queue.position]);

    queue.position = queue.position+1;

    console.log(queue.position) // 1

    let streams = queue.currentQueue[queue.position].streams

    if (streams.length == 1){
        audio.src = queue.currentQueue[queue.position].streams[0].streamUrl
    } else {
        for (let index = 0; index < streams.length; index++) {
            const element = streams[index];
            if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"opus\"")){
                audio.src = element.streamUrl
            } else if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("audio/mp4;")){
                audio.src = element.streamUrl
            } else {
                audio.src = streams[0]
            }
        }
    }

    console.log(queue.currentQueue[queue.position])

    queue.currentSong = queue.currentQueue[queue.position];

    audio.play();
}

function shuffle(){
    // randomizes the list?
    // can i make this better instead of a simple random?
}

// event listener that checks for the current playlist and playing song
// so i check for data-playlist and data-music and change the icon / function related to it?
