# Tremoco

# O que é o Tremoco

Tremoco é um projeto de streaming de musicas(simples), não tem como objetivo de fazer reencoding e etcs...

# Como instalar

Ter uma base de dados `POSTGRES` e definir ligação `POSTGRES_URI` 

Transfere o Tremoco [aqui](https://github.com/atjoao/concordo/releases) e faz run a este comando

Iniciar powershell e executar:
```ps
Get-Content ".env" | ForEach-Object { if ($_ -notmatch '^#') { $name, $value = $_ -split '=', 2; $value = $value.Trim().Trim('"'); [System.Environment]::SetEnvironmentVariable($name.Trim(), $value, [System.EnvironmentVariableTarget]::Process) } }; .\music.exe
```
O comando acima define as variaveis de sistema necessarias para executar.

As variaveis podem ser trocadas no ficheiro `.env`

Tremoco tem uma WebUI > https://localhost:3000

# Endpoints disponiveis
```
/
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