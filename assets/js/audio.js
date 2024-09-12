let todo = null;

const queue = {
    position: 0,
    currentSong: null,

    originalQueue: [], // this queue should not be modified in case of shuffle off
    currentQueue: [], // queue that will be listened for changes
}

const audio = new Audio();

audio.ontimeupdate = (e) => {
    // calculate if 0 - 100% music
    // set width of the player slider watever
};

audio.onended = () => {
    // next track
    // if false clear div content
    // if true start next audio
    // if looping return;
};

audio.onpause = () => {
    // save to localstorage
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
    }

    queue.currentQueue.push(...musicList)
    queue.originalQueue.push(...musicList)

    if (reset) startQueue();

}

function startQueue(){
    // set audio src
    // remove from queue
    // only executes one time

    // first get highest quality (AUDIO_QUALITY_MEDIUM/opus) 
    // if not found select audio/mp4 AUDIO_QUALITY_MEDIUM
    // if not found select for 0 (since the rest is just low tbh)

    stream = todo

    audio.src = queue.originalQueue[queue.position].streams[null].streamUrl
}

function clearQueue(){
    // stops audio
    // clears queue
}

function nextTrack(){
    // set audio src
    // remove from queue
    

    // return false if no more tracks
    // return true if more
}

function shuffle(){
    // randomizes the list?
    // can i make this better instead of a simple random?
}

// event listener that checks for the current playlist and playing song
// so i check for data-playlist and data-music and change the icon / function related to it?
