# screen record

Use webrtc to record the screen and crop the file with ffmpeg.

## Install

require:

- ffmpeg
- go

```bash
go install github.com/fzdwx/screenrd@main
```

## Usage

```bash
screenrd -p 8080 # default port is 8080

# open browser and visit http://localhost:8080
```

## Screenshot

![screenshot](.github/show.gif)
