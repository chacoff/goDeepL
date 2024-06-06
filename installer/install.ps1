<#
    Date: 2024-06-06
    Developer: J.
    Email: jaime.gomez@usach.cl

    Description:
    This PowerShell script is an installer for the goDeepL application.
    It retrieves the current user's profile path and appends a specific 
    subdirectory (.goDeepL) to it. The script can be used to set up or 
    manage the goDeepL application directory within the user's profile directory.
#>

# gets User's profile folder

$userProfilePath = $env:USERPROFILE
$subdirectory = ".goDeepL"
$fullPath = Join-Path -Path $userProfilePath -ChildPath $subdirectory
Write-Output $fullPath

# if the directory does not exist
if (-not (Test-Path $fullPath -PathType Container)) {
    New-Item -Path $fullPath -ItemType Directory -Force
    Write-Output "Created directory: $fullPath"
}

# adds path to windows %PATH%

$newPath = $fullPath
$currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)

if ($currentPath -notlike "*$newPath*") {
    $updatedPath = "$currentPath;$newPath"
    
    [Environment]::SetEnvironmentVariable("PATH", $updatedPath, [EnvironmentVariableTarget]::User)
} else {
    Write-Output "The path is already in the PATH variable."
}

# $env:PATH

# copy goDeepL executable
$sourceFilePath = ".\goDeepL.exe"
Copy-Item -Path $sourceFilePath -Destination $fullPath

# close session
Write-Output "goDeepL installed successfully at: $fullPath"
Exit