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
        if (!r.ok) throw new Error("An error occured while executing this request");
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


document.addEventListener("DOMContentLoaded", () => {
    getSidebar();
});

document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
        // if modal is up, remove it
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
                img.src = music.thumbnail;
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

                divButtons.append(button, play);

                div.append(img, p,author, provider, divButtons);

                searchContainer.appendChild(div);
            })



            break;
        }
        case "playlist":{
            let playlist;
            playlistContainer.setAttribute("class", "")
            initalContainer.setAttribute("class", "hidden")
            searchContainer.setAttribute("class", "hidden")

            playlistContainer.innerHTML = ""
            searchForm.value = null;
            searchForm.innerText = null;

            fetch("/api/playlist/"+content).then((r)=> r.json())
            .then((resp)=>{
                playlist = resp.playlist
                console.log(playlist)

                const e = document.createElement('div');
    
                const e0 = document.createElement('div');
                e0.setAttribute('class', 'head');

                const e1 = document.createElement('div');
                
                const e2 = document.createElement("img");
                e2.src = playlist.image == "" ? "assets/images/default_album.png" : playlist.image;
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
                e8.innerHTML = `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M16.6582 9.28638C18.098 10.1862 18.8178 10.6361 19.0647 11.2122C19.2803 11.7152 19.2803 12.2847 19.0647 12.7878C18.8178 13.3638 18.098 13.8137 16.6582 14.7136L9.896 18.94C8.29805 19.9387 7.49907 20.4381 6.83973 20.385C6.26501 20.3388 5.73818 20.0469 5.3944 19.584C5 19.053 5 18.1108 5 16.2264V7.77357C5 5.88919 5 4.94701 5.3944 4.41598C5.73818 3.9531 6.26501 3.66111 6.83973 3.6149C7.49907 3.5619 8.29805 4.06126 9.896 5.05998L16.6582 9.28638Z" stroke="#000000" stroke-width="2" stroke-linejoin="round"/>
                                </svg>`;
                e7.appendChild(e8);
                
                const e9 = document.createElement('span');
                e9.innerHTML = `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M12.5535 16.5061C12.4114 16.6615 12.2106 16.75 12 16.75C11.7894 16.75 11.5886 16.6615 11.4465 16.5061L7.44648 12.1311C7.16698 11.8254 7.18822 11.351 7.49392 11.0715C7.79963 10.792 8.27402 10.8132 8.55352 11.1189L11.25 14.0682V3C11.25 2.58579 11.5858 2.25 12 2.25C12.4142 2.25 12.75 2.58579 12.75 3V14.0682L15.4465 11.1189C15.726 10.8132 16.2004 10.792 16.5061 11.0715C16.8118 11.351 16.833 11.8254 16.5535 12.1311L12.5535 16.5061Z" fill="#000"/>
                                    <path d="M3.75 15C3.75 14.5858 3.41422 14.25 3 14.25C2.58579 14.25 2.25 14.5858 2.25 15V15.0549C2.24998 16.4225 2.24996 17.5248 2.36652 18.3918C2.48754 19.2919 2.74643 20.0497 3.34835 20.6516C3.95027 21.2536 4.70814 21.5125 5.60825 21.6335C6.47522 21.75 7.57754 21.75 8.94513 21.75H15.0549C16.4225 21.75 17.5248 21.75 18.3918 21.6335C19.2919 21.5125 20.0497 21.2536 20.6517 20.6516C21.2536 20.0497 21.5125 19.2919 21.6335 18.3918C21.75 17.5248 21.75 16.4225 21.75 15.0549V15C21.75 14.5858 21.4142 14.25 21 14.25C20.5858 14.25 20.25 14.5858 20.25 15C20.25 16.4354 20.2484 17.4365 20.1469 18.1919C20.0482 18.9257 19.8678 19.3142 19.591 19.591C19.3142 19.8678 18.9257 20.0482 18.1919 20.1469C17.4365 20.2484 16.4354 20.25 15 20.25H9C7.56459 20.25 6.56347 20.2484 5.80812 20.1469C5.07435 20.0482 4.68577 19.8678 4.40901 19.591C4.13225 19.3142 3.9518 18.9257 3.85315 18.1919C3.75159 17.4365 3.75 16.4354 3.75 15Z" fill="#000"/>
                                </svg>`;
                e7.appendChild(e9);
                
                const e10 = document.createElement('span');
                e10.innerHTML = `<svg width="1em" height="1em" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M18 6L17.1991 18.0129C17.129 19.065 17.0939 19.5911 16.8667 19.99C16.6666 20.3412 16.3648 20.6235 16.0011 20.7998C15.588 21 15.0607 21 14.0062 21H9.99377C8.93927 21 8.41202 21 7.99889 20.7998C7.63517 20.6235 7.33339 20.3412 7.13332 19.99C6.90607 19.5911 6.871 19.065 6.80086 18.0129L6 6M4 6H20M16 6L15.7294 5.18807C15.4671 4.40125 15.3359 4.00784 15.0927 3.71698C14.8779 3.46013 14.6021 3.26132 14.2905 3.13878C13.9376 3 13.523 3 12.6936 3H11.3064C10.477 3 10.0624 3 9.70951 3.13878C9.39792 3.26132 9.12208 3.46013 8.90729 3.71698C8.66405 4.00784 8.53292 4.40125 8.27064 5.18807L8 6" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                </svg>`;
                e7.appendChild(e10);
                
                /* const e11 = document.createElement('span');
                e11.textContent = `b`;
                e7.appendChild(e11); */
                
                e6.appendChild(e7);

                const e12 = document.createElement('div');
                const e13 = document.createElement("input");
                e13.type = "text";
                e13.setAttribute('name', '""');
                e13.setAttribute('id', '""');
                
                e12.appendChild(e13);
                e6.appendChild(e12);
                e0.appendChild(e6);

                playlistContainer.appendChild(e0)

            })

            const divMusic = `
                <div data-musicid=%%>
                    <img src=%image%>
                    <div>
                        <p>%song:name%</p>
                        <p>%song:author%</p>
                    </div>
                    <p>%duration%</p>
                </div>
            `

            break;

        }
        default:{
            document.querySelector("title").innerText = `Musica | Home`

            initalContainer.setAttribute("class", "")
            searchContainer.setAttribute("class", "hidden")
            playlistContainer.setAttribute("class", "hidden")

            searchForm.value = null;
            searchForm.innerText = null;

            break;

        }
    }
}