# SubExtract

SubExtract is a Go application that automatically extracts embedded subtitles from video files. It recursively searches through directories, detects video files with embedded subtitle tracks, and extracts them into separate SRT files.

## Features

- Automatic discovery of video files in folders and subfolders
- Multi-language subtitle extraction (preserves original language codes)
- Parallel processing using worker pools for improved performance
- Skips files that already have extracted subtitles
- Supports multiple video formats (MP4, MKV, AVI, MOV, WMV, FLV, WEBM)

## Requirements (locally)

- Go 1.16 or higher
- FFmpeg (must be installed and available in your system PATH)

## Deployment Options

### Using Docker (Recommended)
You can find the Docker image here: [subtitle-extract image](https://hub.docker.com/r/brunopoiano/subtitle-extract)

#### Option 1: Pull from Docker Hub
```bash
docker run -d \
  --name subextract \
  --restart unless-stopped \
  -e TZ=America/Sao_Paulo \
  -v /location/of/videos:/app/subextract/videos:rw \
  docker.io/brunopoiano/subtitle-extract
```

#### Option 2: Using Docker Compose
```bash
git clone https://github.com/BrunoPoiano/subtitle-extract-go
cd subtitle-extract-go
docker compose up -d
```

### Running from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/BrunoPoiano/subtitle-extract-go
   cd subtitle-extract-go
   ```

2. Configure the root directory where your videos are stored by modifying the `rootFolder` variable in `main.go`:
   ```go
   var rootFolder = "/path/to/your/videos"
   ```

3. Test the application:
   ```bash
   go run .
   ```

4. Build the application:
   ```bash
   go build -o subextract .
   ```

5. Run the compiled binary:
   ```bash
   ./subextract
   ```


## Performance Considerations

- The application uses a worker pool pattern to process multiple video files concurrently
- By default, it utilizes all available CPU cores minus one to avoid system overload
- For large libraries, the initial scan may take some time

## Troubleshooting

- Ensure FFmpeg is properly installed and accessible in your PATH
- Check permissions on your video directories
- Review logs for any error messages
- For Docker installations, ensure your volume mounts are correct

## Acknowledgments

- FFmpeg for the underlying subtitle extraction functionality
