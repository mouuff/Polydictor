
function Build {
    cmd /c "set GOOS=linux&& set GOARCH=amd64&& go build -o polydictor ./cmd/polydictor"
}

function Invoke-Remote {
    param(
        [string]$Target,
        [string]$Command
    )

    Start-Process ssh -ArgumentList @("-t", $Target, $Command) -NoNewWindow -Wait
}

function BuildAndDeploy {
    param(
        [Parameter(Mandatory=$true, Position=0)]
        [string]$User,

        [Parameter(Mandatory=$true, Position=1)]
        [string]$Hostname
    )

    Build

    $Target = "$User@$Hostname"

    Invoke-Remote $Target "sudo systemctl stop polydictor-serve && sudo systemctl stop polydictor-analyze"

    scp .\polydictor "${Target}:/home/chetyredva/Polydictor/"

    Invoke-Remote $Target "sudo systemctl start polydictor-serve && sudo systemctl start polydictor-analyze"
}