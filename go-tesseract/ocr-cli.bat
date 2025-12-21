@echo off
REM Helper script to run ocr-cli via Docker
REM Usage: .\ocr-cli.bat [IMAGE_PATH] [FLAGS...]

IF "%~1"=="" (
    echo Usage: .\ocr-cli.bat [IMAGE_PATH]
    exit /b 1
)

REM Use full path to docker.exe if needed, or assume it's in PATH
SET DOCKER_CMD="C:\Program Files\Docker\Docker\resources\bin\docker.exe"

REM Run CLI container
%DOCKER_CMD% run --rm -v "%cd%":/app/data -t go-tesseract-ocr-api ./ocr-cli -image /app/data/%1 %2 %3 %4 %5