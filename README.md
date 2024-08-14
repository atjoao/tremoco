# musica

Run in powerhsell 
```ps
Get-Content ".env" | ForEach-Object { if ($_ -notmatch '^#') { $name, $value = $_ -split '=', 2; $value = $value.Trim().Trim('"'); [System.Environment]::SetEnvironmentVariable($name.Trim(), $value, [System.EnvironmentVariableTarget]::Process) } }; go build; .\music.exe
```