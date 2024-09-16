let playlist;

const pageContent = document.getElementById("content");
const searchContainer = document.getElementById("search");
const initalContainer = document.getElementById("initial");
const playlistContainer = document.getElementById("playlist");
const searchForm = document.getElementById("searchForm");

const getSidebar = () => {
    fetch("/html/sidebar")
        .then(r => {
            return r.text();
        })
        .then(html => {
            document.getElementById("sidebar").innerHTML = html;
        })
}

window.executeForm = (form) => {
    fetch(form.action, {
        method: form.method,
        body: new FormData(form)
    }).then(r => {
        if (!r.ok) throw new Error("An error occured while executing this request\n" + r.statusText);
        return r.text();
    }).then(resp => {
        closeModal();
    }).catch(err => {
        alert(err);
    });
}

window.closeModal = () => {
    document.querySelector(".modal")?.remove();
};

window.createModal = (modal) => {
    if (!modal.includes(".html")) return;

    fetch("/assets/modals/" + modal)
        .then((r) => {
            if (!r.ok) throw new Error("Failed to load the modal content");
            return r.text();
        })
        .then((rst) => {
            const parser = new DOMParser();
            const content = parser.parseFromString(rst, "text/html");

            let filterArray = Array.from(content.body.children).filter((src)=>{
                if (src.localName == "script"){
                    return false;
                }
                return true;
            })

            document.body.append(...filterArray);

            content.querySelectorAll("script").forEach(script => {
                console.log("Executing Script: ", script)
                try {
                    new Function(script.innerHTML)(); 
                } catch (err) {
                    console.error("Script execution error:", err);
                }
            });
        })
        .catch((e) => console.error("Error loading modal:", e));
};

function addRecentPlayed(type, info){
    const recentPlayed = JSON.parse(localStorage.getItem("recentPlayed"));
    if (recentPlayed.length >= 10){
        recentPlayed.shift();
    }

    recentPlayed.some((element, index)=>{
        if (element.type == type && element.info.videoid == info.videoid){
            recentPlayed.splice(index, 1);
            return true;
        }
    })

    recentPlayed.push({type, info});
    localStorage.setItem("recentPlayed", JSON.stringify(recentPlayed));
}

function loadRecentlyPlayed() {
    document.getElementById("recentlyPlayed").innerHTML = "";

    if (localStorage.getItem("recentPlayed") === null){
        localStorage.setItem("recentPlayed", JSON.stringify([]));
    } else {
        const recentPlayed = JSON.parse(localStorage.getItem("recentPlayed"));
        if (recentPlayed.length > 0){
            recentPlayed.forEach((element)=>{
                if (element.type == "music"){
                    const html = `
                    <div class="recentPlayedContainer" data-playlistid="%music:id%">
                        <img src="%music:image%"/>
                        <p>%music:name%</p>
                        <p>%music:author%</p>
                        <p>%music:length%</p>
                        <button id="play_%music:id%">Play</button>
                    </div>`

                    document.getElementById("recentlyPlayed").insertAdjacentHTML("beforeend",
                        html.replace("%music:image%", element.info.thumbnails[0].url.includes("https://") ? "/api/proxy?url=" + btoa(element.info.thumbnails[0].url) : element.info.thumbnails[0].url)
                            .replace("%music:name%", element.info.title)
                            .replace("%music:length%", fmtMSS(element.info.duration))
                            .replace("%music:id%", element.info.videoid)
                            .replace("%music:id%", element.info.videoid)
                            .replace("%music:author%", element.info.author)
                        );

                    document.getElementById("play_"+element.info.videoid).addEventListener("click", function (e) {
                        e.preventDefault();
                        e.stopPropagation();
                        fetch("/api/video?id=" + element.info.videoid).then(r => r.json()).then(data => {
                            sendToQueue([data.data], true);
                            addRecentPlayed("music", data.data);
                        });
                    });

                }
                if (element.type == "playlist"){
                    const html = `
                    <div class="recentPlayedContainer" data-playlistid="%playlist:id%">
                        <img src="%playlist:image%"/>
                        <p>%playlist:name%</p>
                        <p>%playlist:length% songs</p>
                        <button onclick='replaceContent(\"playlist\", this.parentElement.dataset.playlistid)'>Open</button>
                        <button id="play_%playlist:id%">Play</button>
                    </div>
                    `

                    document.getElementById("recentlyPlayed").insertAdjacentHTML("beforeend", 
                        html
                            .replace("%playlist:image%", element.info.image)
                            .replace("%playlist:name%", element.info.name)
                            .replace("%playlist:length%", element.info.list.length)
                            .replace("%playlist:id%", element.info.id)
                            .replace("%playlist:id%", element.info.id)
                        );

                    document.getElementById("play_"+element.info.id).addEventListener("click", function (e) {
                        e.preventDefault();
                        e.stopPropagation();
                        fetch("/api/playlist/"+element.info.id).then((r)=> r.json()).then((resp)=>{
                            sendToQueue(resp.playlist.list, true);
                            addRecentPlayed("playlist", resp.playlist);
                        })
                    });

                }
            })
        }
    }
}

document.addEventListener("DOMContentLoaded", () => {
    getSidebar();
    loadRecentlyPlayed();
});

document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
        document.querySelector(".modal")?.remove();
    }
});


function debounce(callback, delay) {
    let timeout;
    return (...args) => {
        if (timeout) clearTimeout(timeout);
        timeout = setTimeout(() => callback(...args), delay);
    };
}

async function performSearch(query) {
    try {
        const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
        if (!res.ok) throw new Error(`Error: ${res.statusText}`);
        const data = await res.json();
        console.log(data.videos);

        replaceContent("search", data.videos)

    } catch (error) {
        console.error('Failed to fetch search results:', error);
    }
}

const debouncedSearch = debounce(() => {
    const query = searchForm.value.trim();
    if (query.length > 0) {
        performSearch(query);
    }
    if (query.length === 0) {
        replaceContent("reset", null);
    }
}, 300);

searchForm.addEventListener("input", debouncedSearch);

function clearText(value){
    const element = document.createElement("div")
    element.innerText = value
    return element.innerHTML
}

function fmtMSS(seconds) {
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    
    return [
        hrs.toString().padStart(2, '0'),
        mins.toString().padStart(2, '0'),
        secs.toString().padStart(2, '0')
    ].join(':');
}


function replaceContent(type, content){
    
    switch(type){
        case "search":{
            document.querySelector("title").innerText = `Musica | Search ${clearText(searchForm.value)}`
            searchContainer.setAttribute("class", "")
            playlistContainer.setAttribute("class", "hidden")
            initalContainer.setAttribute("class", "hidden")

            searchContainer.innerHTML = `<h1 style="height: fit-content;display: block;width: 100%;">Search Results for ${clearText(searchForm.value)}</h1>`;

            content.map(music => {
                const div = document.createElement("div");
                div.classList.add("music");
                div.dataset.musicid = music.id;
                
                const img = document.createElement("img");
                img.src = music.thumbnail.includes("https://") ? "/api/proxy?url=" + btoa(music.thumbnail) : music.thumbnail;
                img.alt = "Thumbnail";

                const p = document.createElement("p");
                p.classList.add("title");
                p.textContent = music.title;

                const author = document.createElement("p");
                author.textContent = "Author: "+ music.author;

                const provider = document.createElement("p");
                provider.textContent = "Provider: "+ music.provider;

                const divButtons = document.createElement("div");
                divButtons.style.display = "flex";
                divButtons.style.justifyContent = "space-between";

                const button = document.createElement("button");
                button.textContent = "Add to Playlist";
                button.onclick = () => {
                    globalThis.addToPlaylistMusicId = music.id;
                    createModal("addToPlaylist.html");
                }

                const play = document.createElement("button");
                play.textContent = "Play";
                play.onclick = () => {
                    fetch("/api/video?id=" + music.id).then(r => r.json()).then(data => {
                        sendToQueue([data.data], true);
                        addRecentPlayed("music", data.data);
                    });

                }

                divButtons.append(button, play);

                div.append(img, p,author, provider, divButtons);

                searchContainer.appendChild(div);
            })

            break;
        }
        case "playlist":{
            playlistContainer.setAttribute("class", "")
            initalContainer.setAttribute("class", "hidden")
            searchContainer.setAttribute("class", "hidden")

            playlistContainer.innerHTML = ""
            searchForm.value = null;
            searchForm.innerText = null;

            const divMusic = `
                <div data-musicid=%song:id% class="music_column">
                    <div>
                        <div class="relative">
                            <div id="PlayOverlay">
                                <svg id="playbtn" width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M16.6582 9.28638C18.098 10.1862 18.8178 10.6361 19.0647 11.2122C19.2803 11.7152 19.2803 12.2847 19.0647 12.7878C18.8178 13.3638 18.098 13.8137 16.6582 14.7136L9.896 18.94C8.29805 19.9387 7.49907 20.4381 6.83973 20.385C6.26501 20.3388 5.73818 20.0469 5.3944 19.584C5 19.053 5 18.1108 5 16.2264V7.77357C5 5.88919 5 4.94701 5.3944 4.41598C5.73818 3.9531 6.26501 3.66111 6.83973 3.6149C7.49907 3.5619 8.29805 4.06126 9.896 5.05998L16.6582 9.28638Z" stroke="#FFF" fill="#FFF" stroke-width="2" stroke-linejoin="round"/>
                                </svg>
                            </div>
                            <img src=%image%>
                        </div>
                        <div>
                            <p id="title">%song:name%</p>
                            <p id="author">%song:author%</p>
                        </div>
                    </div>
                    <div class="musicInfo">
                        <svg class="playlistRemove" width="1.5em" height="1.5em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M10 11V17" stroke="#FFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M14 11V17" stroke="#FFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M4 7H20" stroke="#FFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M6 7H12H18V18C18 19.6569 16.6569 21 15 21H9C7.34315 21 6 19.6569 6 18V7Z" stroke="#FFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            <path d="M9 5C9 3.89543 9.89543 3 11 3H13C14.1046 3 15 3.89543 15 5V7H9V5Z" stroke="#FFF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        <p>%duration%</p>
                    </div>
                </div>
            `

            fetch("/api/playlist/"+content).then((r)=> r.json())
            .then((resp)=>{
                globalThis.currentPlaylist = content;
                const musicContainer = document.createElement("div")
                playlist = resp.playlist
                console.log(playlist)

                const e = document.createElement('div');
    
                const e0 = document.createElement('div');
                e0.setAttribute('class', 'head');

                const e1 = document.createElement('div');
                e1.setAttribute('style', 'display: flex;align-items: center;gap: 20px;');
                
                const e2 = document.createElement("img");
                e2.src = playlist.image == "" ? "assets/images/default_album.png" : playlist.image;
                e2.width = 128;
                e2.height = 128;
                e2.style.borderRadius = "10px";
                e2.style.objectFit = "cover";

                e1.appendChild(e2);

                const e3 = document.createElement('div');
                e3.setAttribute('class', 'info');
                const e4 = document.createElement('h1');
                e4.innerText = playlist.name

                e3.appendChild(e4);
                const e5 = document.createElement('p');
                e5.innerText = playlist.list?.length ? `${playlist.list.length} songs` : "0 songs";
                e3.appendChild(e5);
                e1.appendChild(e3);
                
                e0.appendChild(e1);
                const e6 = document.createElement('div');

                const e7 = document.createElement('div');
                e7.setAttribute('class', 'buttons');

                const e8 = document.createElement('span');
                e8.innerHTML = `<svg data-playlist=${content} id="PlaylistPlay" width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M16.6582 9.28638C18.098 10.1862 18.8178 10.6361 19.0647 11.2122C19.2803 11.7152 19.2803 12.2847 19.0647 12.7878C18.8178 13.3638 18.098 13.8137 16.6582 14.7136L9.896 18.94C8.29805 19.9387 7.49907 20.4381 6.83973 20.385C6.26501 20.3388 5.73818 20.0469 5.3944 19.584C5 19.053 5 18.1108 5 16.2264V7.77357C5 5.88919 5 4.94701 5.3944 4.41598C5.73818 3.9531 6.26501 3.66111 6.83973 3.6149C7.49907 3.5619 8.29805 4.06126 9.896 5.05998L16.6582 9.28638Z" stroke="#000000" stroke-width="2" stroke-linejoin="round"/>
                                </svg>`;
                e8.addEventListener("click", function (e) {
                    if(playlist.list){
                        sendToQueue(playlist.list, true)
                        addRecentPlayed("playlist", playlist)
                    }
                })
                e7.appendChild(e8);

                const e10 = document.createElement('span');
                e10.innerHTML = `<svg width="1em" height="1em" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path stroke-width="2" stroke-linejoin="round" stroke="#FFF" stroke-linecap="round" fill="black" style="fill: black !important;" d="M18 6L17.1991 18.0129C17.129 19.065 17.0939 19.5911 16.8667 19.99C16.6666 20.3412 16.3648 20.6235 16.0011 20.7998C15.588 21 15.0607 21 14.0062 21H9.99377C8.93927 21 8.41202 21 7.99889 20.7998C7.63517 20.6235 7.33339 20.3412 7.13332 19.99C6.90607 19.5911 6.871 19.065 6.80086 18.0129L6 6M4 6H20M16 6L15.7294 5.18807C15.4671 4.40125 15.3359 4.00784 15.0927 3.71698C14.8779 3.46013 14.6021 3.26132 14.2905 3.13878C13.9376 3 13.523 3 12.6936 3H11.3064C10.477 3 10.0624 3 9.70951 3.13878C9.39792 3.26132 9.12208 3.46013 8.90729 3.71698C8.66405 4.00784 8.53292 4.40125 8.27064 5.18807L8 6"></path>
                                </svg>`;
                e10.addEventListener("click", function (e) {
                    createModal("deletePlaylist.html");
                })
                e7.appendChild(e10);

                const e47 = document.createElement('span');
                e47.innerHTML = `<svg fill="#000000" width="1em" height="1em" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                    <path fill-rule="evenodd" d="M12.3023235,7.94519388 L4.69610276,15.549589 C4.29095108,15.9079238 4.04030835,16.4092335 4,16.8678295 L4,20.0029438 L7.06398288,20.004826 C7.5982069,19.9670062 8.09548693,19.7183782 8.49479322,19.2616227 L16.0567001,11.6997158 L12.3023235,7.94519388 Z M13.7167068,6.53115006 L17.4709137,10.2855022 L19.8647941,7.89162181 C19.9513987,7.80501747 20.0000526,7.68755666 20.0000526,7.56507948 C20.0000526,7.4426023 19.9513987,7.32514149 19.8647932,7.23853626 L16.7611243,4.13485646 C16.6754884,4.04854589 16.5589355,4 16.43735,4 C16.3157645,4 16.1992116,4.04854589 16.1135757,4.13485646 L13.7167068,6.53115006 Z M16.43735,2 C17.0920882,2 17.7197259,2.26141978 18.1781068,2.7234227 L21.2790059,5.82432181 C21.7406843,6.28599904 22.0000526,6.91216845 22.0000526,7.56507948 C22.0000526,8.21799052 21.7406843,8.84415992 21.2790068,9.30583626 L9.95750718,20.6237545 C9.25902448,21.4294925 8.26890003,21.9245308 7.1346,22.0023295 L2,22.0023295 L2,21.0023295 L2.00324765,16.7873015 C2.08843822,15.7328366 2.57866679,14.7523321 3.32649633,14.0934196 L14.6953877,2.72462818 C15.1563921,2.2608295 15.7833514,2 16.43735,2 Z"/>
                                </svg>`

                e47.addEventListener("click", function (e) {
                    createModal("editPlaylist.html");
                })

                e7.appendChild(e47);

                const e9 = document.createElement('span');
                e9.innerHTML = `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M12.5535 16.5061C12.4114 16.6615 12.2106 16.75 12 16.75C11.7894 16.75 11.5886 16.6615 11.4465 16.5061L7.44648 12.1311C7.16698 11.8254 7.18822 11.351 7.49392 11.0715C7.79963 10.792 8.27402 10.8132 8.55352 11.1189L11.25 14.0682V3C11.25 2.58579 11.5858 2.25 12 2.25C12.4142 2.25 12.75 2.58579 12.75 3V14.0682L15.4465 11.1189C15.726 10.8132 16.2004 10.792 16.5061 11.0715C16.8118 11.351 16.833 11.8254 16.5535 12.1311L12.5535 16.5061Z" fill="#000"/>
                                    <path d="M3.75 15C3.75 14.5858 3.41422 14.25 3 14.25C2.58579 14.25 2.25 14.5858 2.25 15V15.0549C2.24998 16.4225 2.24996 17.5248 2.36652 18.3918C2.48754 19.2919 2.74643 20.0497 3.34835 20.6516C3.95027 21.2536 4.70814 21.5125 5.60825 21.6335C6.47522 21.75 7.57754 21.75 8.94513 21.75H15.0549C16.4225 21.75 17.5248 21.75 18.3918 21.6335C19.2919 21.5125 20.0497 21.2536 20.6517 20.6516C21.2536 20.0497 21.5125 19.2919 21.6335 18.3918C21.75 17.5248 21.75 16.4225 21.75 15.0549V15C21.75 14.5858 21.4142 14.25 21 14.25C20.5858 14.25 20.25 14.5858 20.25 15C20.25 16.4354 20.2484 17.4365 20.1469 18.1919C20.0482 18.9257 19.8678 19.3142 19.591 19.591C19.3142 19.8678 18.9257 20.0482 18.1919 20.1469C17.4365 20.2484 16.4354 20.25 15 20.25H9C7.56459 20.25 6.56347 20.2484 5.80812 20.1469C5.07435 20.0482 4.68577 19.8678 4.40901 19.591C4.13225 19.3142 3.9518 18.9257 3.85315 18.1919C3.75159 17.4365 3.75 16.4354 3.75 15Z" fill="#000"/>
                                </svg>`;
                e9.addEventListener("click", function (e) {
                    // download playlist
                    if (!playlist.list) return;
                    alert("todo...")
                });

                e7.appendChild(e9);
                
                e6.appendChild(e7);
                e6.setAttribute("style", "justify-content: space-between;display: flex;padding: 10px;")
                e6.style.justifyContent = "space-between"

                const e12 = document.createElement('div');
                const e13 = document.createElement("input");
                e13.type = "text";
                e13.setAttribute('name', "search");
                e13.setAttribute('id', "playlistSearch");
                e13.placeholder = "Search inside the playlist"
                e13.addEventListener("keyup", function (e) {
                    if (e.key.includes("Arrow")) return;
                    
                    const filterValue = e13.value.toLowerCase();

                    const musicChilds = Array.from(musicContainer.children);
                    
                    musicChilds.forEach(musicItem => {
                        const musicName = musicItem.querySelector("p#title").textContent.toLowerCase();
                        const authorName = musicItem.querySelector("p#author").textContent.toLowerCase();
                        if (musicName.includes(filterValue) || authorName.includes(filterValue)) {
                            musicItem.style.display = "flex";
                        } else {
                            musicItem.style.display = "none";
                        }
                    });
                })

                
                e12.appendChild(e13);
                e6.appendChild(e12);
                e0.appendChild(e6);

                playlistContainer.appendChild(e0)

                musicContainer.setAttribute("class", "music_list")
                if (playlist.list != null){
                    for (let index = 0; index < playlist.list.length; index++) {
                        const element = playlist.list[index];
                        musicContainer.insertAdjacentHTML("beforeend", 
                        divMusic
                            .replace("%image%", element.thumbnails[0].url.includes("https://") ? "/api/proxy?url=" + btoa(element.thumbnails[0].url) : element.thumbnails[0].url)
                            .replace("%song:id%", element.videoid)
                            .replace("%song:name%", element.title)
                            .replace("%song:author%", element.author)
                            .replace("%duration%", fmtMSS(element.duration))
                        )
                    }    
                }

                playlistContainer.appendChild(musicContainer)

                document.querySelector("title").innerText = `Musica | Playlist ${playlist.name}`
                document.querySelectorAll("#PlayOverlay").forEach((element)=>{
                    element.addEventListener("click", function (e) {
                        if(playlist.list){
                            const playlistFiltered = playlist.list.filter((song)=>{
                                return song.videoid == element.parentElement.parentElement.parentElement.dataset.musicid
                            })
                            sendToQueue(playlistFiltered, true)
                            addRecentPlayed("music", playlistFiltered[0])
                        }
                    })
                })

                document.querySelectorAll(".playlistRemove").forEach((element)=>{
                    element.addEventListener("click", function (e) {
                        const musicId = element.parentElement.parentElement.dataset.musicid
                        const formdt = new FormData()
                        formdt.append("playlistId", content)
                        formdt.append("audioId", musicId)

                        fetch("/api/playlist/change", {
                            method: "POST",
                            body: formdt
                        })
                        .then((r)=> r.json())
                        .then((resp)=>{
                            element.parentElement.parentElement.remove()
                            fetch("/api/playlist/"+content).then((r)=> r.json()).then((resp)=>{
                                playlist = resp.playlist
                            })
                        })
                    })
                })
            })

            break;
        }
        default:{
            document.querySelector("title").innerText = `Musica | Home`

            initalContainer.setAttribute("class", "")
            searchContainer.setAttribute("class", "hidden")
            playlistContainer.setAttribute("class", "hidden")

            searchForm.value = null;
            searchForm.innerText = null;

            loadRecentlyPlayed();

            break;

        }
    }
}