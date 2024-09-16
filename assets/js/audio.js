const playerToggle = document.getElementById("player-toggle")
playerToggle.addEventListener("click", () => {
    if (audio.paused){
        audio.play()
    } else {
        audio.pause()
    }
})

const playerRepeat = document.getElementById("player-repeat")
playerRepeat.addEventListener("click", () => {
    if (audio.loop){
        audio.loop = false
        playerRepeat.children[0].children[0].style.stroke = "gray"
    } else {
        audio.loop = true
        playerRepeat.children[0].children[0].style.stroke = "white"
    }
})

const playerImage = document.getElementById("player-image")
const playerTitle = document.getElementById("player-title") 
const playerAuthor = document.getElementById("player-author")
const playerVolume = document.getElementById("player-volume")

const playerProgress = document.getElementById("player-progress")
const playerCurrentTime = document.getElementById("player-current")
const playerDuration = document.getElementById("player-duration")

const playerBack = document.getElementById("player-back")
playerBack.addEventListener("click", () => {
    backTrack();
})

const playerNext = document.getElementById("player-next")
playerNext.addEventListener("click", () => {
    nextTrack();
})

const playerShuffle = document.getElementById("player-shuffle")
playerShuffle.addEventListener("click", () => {
    toggleShuffle();
});

function toggleShuffle(force) {
    if (audio.paused && !force) return;

    if (queue.originalQueue.length == 1 ) return;

    queue.shuffle = !queue.shuffle;
    
    if (queue.shuffle) {
        shuffleQueue();
        playerShuffle.children[0].children[0].style.stroke = "white";
    } else {
        resetQueueOrder();
        playerShuffle.children[0].children[0].style.stroke = "gray";
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const volume = localStorage.getItem("volume")
    if (volume) {
        playerVolume.value = volume
        audio.volume = volume / 100
    } else {
        playerVolume.value = 100
        audio.volume = 1
        localStorage.setItem("volume", 100)
    }

    playerProgress.value = 0
})

playerVolume.addEventListener("change", (e)=> {
    audio.volume = e.target.value / 100;
    localStorage.setItem("volume", e.target.value)
})

playerProgress.addEventListener("change", (e)=>{

    if (queue.currentSong == null) {
        e.preventDefault();
        e.target.value = 0;
        return;
    };


    audio.currentTime = e.target.value
})

const queue = {
    position: 0,
    currentSong: null,
    shuffle: false,

    originalQueue: [], // this queue should not be modified in case of shuffle off
    currentQueue: [], // queue that will be listened for changes
    alteradyPlayed: []
}

const audio = new Audio();

const icons = {
    play: `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M16.6582 9.28638C18.098 10.1862 18.8178 10.6361 19.0647 11.2122C19.2803 11.7152 19.2803 12.2847 19.0647 12.7878C18.8178 13.3638 18.098 13.8137 16.6582 14.7136L9.896 18.94C8.29805 19.9387 7.49907 20.4381 6.83973 20.385C6.26501 20.3388 5.73818 20.0469 5.3944 19.584C5 19.053 5 18.1108 5 16.2264V7.77357C5 5.88919 5 4.94701 5.3944 4.41598C5.73818 3.9531 6.26501 3.66111 6.83973 3.6149C7.49907 3.5619 8.29805 4.06126 9.896 5.05998L16.6582 9.28638Z" stroke="#fff" stroke-width="2" stroke-linejoin="round"/>
            </svg>`,
    pause: `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M8 5V19M16 5V19" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>`
}

audio.onpause = () => {
    playerToggle.innerHTML = icons.play
}

audio.onplay = () => {
    let thumb;

    playerToggle.innerHTML = icons.pause

    if (!queue.currentSong.thumbnails[0].url.includes("https://")){
        thumb = new URL(window.location)
        thumb.pathname = queue.currentSong.thumbnails[0].url
    } else {
        thumb = "/api/proxy?url=" + btoa(queue.currentSong.thumbnails[0].url)
    }
    
    navigator.mediaSession.metadata = new MediaMetadata({
        title: queue.currentSong.title.replace(/\\"/g, '"'),
        artist: queue.currentSong.author,
        artwork: [
            {
                src: thumb.toString(),
                sizes: "96x96",
                type: 'image/png'
            }
        ]
    })

    navigator.mediaSession.setActionHandler('seekbackward', function() {
        seekBackward();
    });
    navigator.mediaSession.setActionHandler('seekforward', function() {
        seekFoward();
    });
    navigator.mediaSession.setActionHandler('previoustrack', function() {
        backTrack();
    });
    navigator.mediaSession.setActionHandler('nexttrack', function() {
        nextTrack();
    });

    playerImage.src = thumb;
    playerTitle.textContent = queue.currentSong.title.replace(/\\"/g, '"');
    playerAuthor.textContent = queue.currentSong.author;

    playerDuration.textContent = fmtMSS(queue.currentSong.duration)

}

audio.ontimeupdate = (e) => {
    playerCurrentTime.textContent = fmtMSS(audio.currentTime)

    playerProgress.max = audio.duration
    playerProgress.ariaValueNow = audio.currentTime
    playerProgress.value = audio.currentTime
};

audio.onended = () => {
    if (audio.loop) return;
    nextTrack();
};

function sendToQueue(musicList, reset){
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

function seekFoward(){
    audio.currentTime += 10
}

function seekBackward(){
    audio.currentTime -= 10
}

function startQueue(){
    let streams = queue.currentQueue[queue.position].streams

    if (streams.length == 1){
        audio.src = queue.currentQueue[queue.position].streams[0].streamUrl
    } else {
        for (let index = 0; index < streams.length; index++) {
            const element = streams[index];
            if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"opus\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"mp4a.40.5\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else {
                audio.src = streams[0].streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(streams[0].streamUrl) : streams[0].streamUrl
            }
        }
    }

    queue.currentSong = queue.currentQueue[queue.position]

    audio.play()

}

function nextTrack(){
    if (!queue.currentQueue[queue.position + 1]) {
        console.log("Next track is undefined, stopping execution.");
        audio.pause();
        audio.src = null;
        navigator.mediaSession.metadata = null;

        playerImage.src = "";
        playerTitle.textContent = "";
        playerAuthor.textContent = "";
        playerDuration.textContent = "00:00:00";
        playerCurrentTime.textContent = "00:00:00";

        playerProgress.value = 0;
        queue.currentSong = null;

        queue.currentQueue = [];
        queue.originalQueue = [];
        queue.alteradyPlayed = [];
        queue.position = 0;

        if (queue.shuffle) {
            toggleShuffle(true);
        }

        return;
    }

    queue.alteradyPlayed.push(queue.currentQueue[queue.position]);

    queue.position = queue.position+1;

    console.log(queue.position)

    let streams = queue.currentQueue[queue.position].streams

    if (streams.length == 1){
        audio.src = queue.currentQueue[queue.position].streams[0].streamUrl
    } else {
        for (let index = 0; index < streams.length; index++) {
            const element = streams[index];
            if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"opus\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"mp4a.40.5\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else {
                audio.src = streams[0].streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(streams[0].streamUrl) : streams[0].streamUrl
            }
        }
    }

    console.log(queue.currentQueue[queue.position])

    queue.currentSong = queue.currentQueue[queue.position];

    audio.play();
}

function backTrack(){
    if (!queue.currentQueue[queue.position - 1]) {
        console.log("Previous track is undefined, stopping execution.");

        return;
    }

    queue.position = queue.position-1;

    console.log(queue.position)

    let streams = queue.currentQueue[queue.position].streams

    if (streams.length == 1){
        audio.src = queue.currentQueue[queue.position].streams[0].streamUrl
    } else {
        for (let index = 0; index < streams.length; index++) {
            const element = streams[index];
            if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"opus\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else if (element.audioQuality.includes("MEDIUM") && element.mimeType.includes("codecs=\"mp4a.40.5\"")){
                audio.src = element.streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(element.streamUrl) : element.streamUrl
            } else {
                audio.src = streams[0].streamUrl.includes("https://") ? "/api/proxy?url=" + btoa(streams[0].streamUrl) : streams[0].streamUrl
            }
        }
    }

    console.log(queue.currentQueue[queue.position])

    queue.currentSong = queue.currentQueue[queue.position];

    audio.play();
}

// bad impl

function shuffleQueue() {
    const unplayed = queue.currentQueue.slice(queue.position + 1);
    
    for (let i = unplayed.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [unplayed[i], unplayed[j]] = [unplayed[j], unplayed[i]];
    }
    
    queue.currentQueue = [
        ...queue.alteradyPlayed, 
        queue.currentQueue[queue.position], 
        ...unplayed
    ];
    
    console.log("Shuffled queue:", queue.currentQueue);
}


function resetQueueOrder() {
    const unplayed = queue.originalQueue.slice(queue.position + 1);
    
    queue.currentQueue = [
        ...queue.alteradyPlayed, 
        queue.currentQueue[queue.position], 
        ...unplayed
    ];

    const originalPosition = queue.originalQueue.findIndex(song => song === queue.currentSong);
    if (originalPosition !== -1) {
        queue.position = originalPosition;
    }
    
    console.log("Reset queue to original order:", queue.currentQueue);
}