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

const callModel_addToPlaylist = async (musicId) => {
    const div = document.createElement("div");
    div.classList.add("modal");

    const content = document.createElement("div");
    content.classList.add("modal-content");

    const div_container = document.createElement("div");
    div_container.classList.add("div_container");

    const h1 = document.createElement("h1");
    h1.textContent = "Add to Playlist";

    const close = document.createElement("span");
    close.style.display = "block";
    close.style.width = "min-content";
    close.classList.add("close");
    close.innerHTML = `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" width="32px" height="32px">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M10.9393 12L6.9696 15.9697L8.03026 17.0304L12 13.0607L15.9697 17.0304L17.0304 15.9697L13.0607 12L17.0303 8.03039L15.9696 6.96973L12 10.9393L8.03038 6.96973L6.96972 8.03039L10.9393 12Z" fill="#FFF"></path>
    </svg>`;
    close.onclick = () => {
        div.remove();
    }

    // .map el add event
    

    div_container.append(h1, close);

    content.append(div_container);

    div.appendChild(content);
    document.body.appendChild(div);
}

const callModel_createPlaylist = async () => {
    const div = document.createElement("div");
    div.classList.add("modal");

    const content = document.createElement("div");
    content.classList.add("modal-content");

    const div_container = document.createElement("div");
    div_container.classList.add("div_container");

    const h1 = document.createElement("h1");
    h1.textContent = "Create Playlist";

    const close = document.createElement("span");
    close.style.display = "block";
    close.style.width = "min-content";
    close.classList.add("close");
    close.innerHTML = `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" width="32px" height="32px">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M10.9393 12L6.9696 15.9697L8.03026 17.0304L12 13.0607L15.9697 17.0304L17.0304 15.9697L13.0607 12L17.0303 8.03039L15.9696 6.96973L12 10.9393L8.03038 6.96973L6.96972 8.03039L10.9393 12Z" fill="#FFF"></path>
    </svg>`;
    close.onclick = () => {
        div.remove();
    }

    const form = document.createElement("form");
    form.action = "/api/playlist/create";
    form.method = "POST";

    const label = document.createElement("label");
    label.htmlFor = "name";
    label.textContent = "Playlist Name:";

    const inputName = document.createElement("input");
    inputName.type = "text";
    inputName.id = "name";
    inputName.name = "name";
    inputName.required = true;

    const submitButton = document.createElement("button");
    submitButton.type = "submit";
    submitButton.innerText = "Create";

    /**
     * @param {Event} e 
     */
    form.addEventListener("submit", (e) => {
        e.preventDefault();

        fetch(form.action, {
            method: form.method,
            body: new FormData(form)
        }).then((r)=> {
            if (!r.ok) throw new Error("Ocorreu um erro ao criar a playlist");
            return r.text();
        }).then((resp)=> {
            div.remove();
            getSidebar();
        }).catch((err)=> {
            alert(err);
        })
    });

    div_container.append(h1, close);

    form.append(label, inputName, submitButton);
    content.append( div_container, form);

    div.appendChild(content);
    document.body.appendChild(div);
}

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

                const provider = document.createElement("p");
                provider.textContent = "Provider: "+ music.provider;

                const divButtons = document.createElement("div");
                divButtons.style.display = "flex";
                divButtons.style.justifyContent = "space-between";

                const button = document.createElement("button");
                button.textContent = "Add to Playlist";

                const play = document.createElement("button");
                play.textContent = "Play";

                divButtons.append(button, play);

                div.append(img, p, provider, divButtons);

                searchContainer.appendChild(div);
            })

            break;
        }
        case "playlist":{
            playlistContainer.setAttribute("class", "")
            initalContainer.setAttribute("class", "hidden")
            searchContainer.setAttribute("class", "hidden")
            break;

        }
        default:{
            document.querySelector("title").innerText = `Musica | Home`

            initalContainer.setAttribute("class", "")
            searchContainer.setAttribute("class", "hidden")
            playlistContainer.setAttribute("class", "hidden")

            break;

        }
    }
}