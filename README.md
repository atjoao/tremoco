# Tremoco

# O que é o Tremoco

Tremoco é um projeto de streaming de musicas(simples), não tem como objetivo de fazer reencoding e etcs...

# Como instalar

Ter
- `FFProbe` em `%PATH%` ou parecido?? (ate arranjar maneira de ler exif data)

Transfere o Tremoco [aqui](https://github.com/atjoao/tremoco/releases/)

Remover .example do ficheiro `.env.example`

As variaveis existentes são as unicas no ficheiro `.env` e configurar-lo

Abrir `tremoco.exe`

```
.
├── .env
├── audio
│   ├── Album
│   │   ├── Disc 1
│   │   │   ├── track1.mp3
│   │   │   ├── track2.opus
│   │   │   └── track3.flac
│   │   ├── Disc 2
│   │   │   ├── track1.aac
│   │   │   ├── track2.wav
│   │   │   └── track3.flac
│   │   └── cover.jpg
└── tremoco.exe
```
> Exemplo da extrutura de albuns 

Tremoco vai estar disponivel em se tiver correto > https://localhost:3000

# Endpoints disponiveis
```
/
├── GET / (page)
├── GET /login (page)
├── GET /register (page)
├── GET /logout (redirect/clear)
├── html
│   └── GET /sidebar
├── api (protected)
│   ├── GET /search (params ?q= )
│   ├── GET /video (params ?id= )
│   ├── GET /stream/:audioId
│   ├── GET /cover/:audioId
│   ├── GET /playlists
│   ├── POST /playlist/create (FORMDATA)
│   ├── GET /playlist/get/:audioId
│   ├── POST /playlist/change (FORMDATA)
│   ├── GET /playlist/:playlistId
│   ├── DELETE /playlist/delete/:playlistId
│   ├── POST /playlist/edit/:playlistId (FORMDATA)
│   └── GET /proxy (params ?url= encoded base64 / content-type audio | image)
└── auth
    ├── POST /login (FORMDATA)
    └── POST /register (FORMDATA)

```
