Start powershell {
    Set-Location -Path .\api
    &go run main.go
}

Start powershell {
    Set-Location -Path ".\web-ui"
    &npm run start
}
