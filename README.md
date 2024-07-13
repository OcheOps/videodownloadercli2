# Video Downloader CLI

Video Downloader CLI is a command-line tool written in Go that allows you to download videos from YouTube and Instagram. without Pesky ads from your comand line 

## Features

- Download YouTube videos
- Download Instagram videos (public posts only)
- Progress bar for download status
- Automatic creation of download directory

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go 1.16 or higher installed on your system
- Internet connection

## Installation

1. Clone the repository:

git clone https://github.com/OcheOps/videodownloadercli2.git

2. Navigate to the project directory:

cd videodownloadercli2

3. Build the project:

go build

## Usage

After building the project, you can use the Video Downloader CLI as follows:

1. For YouTube videos:

./videodownloadercli2 download "https://www.youtube.com/watch?v=VIDEO_ID"

Replace `VIDEO_ID` with the actual YouTube video ID.

2. For Instagram videos:

./videodownloadercli2 download "https://www.instagram.com/p/POST_ID/"

Replace `POST_ID` with the actual Instagram post ID.

Alternatively, you can run the program without building it:


go run main.go download "https://www.youtube.com/watch?v=VIDEO_ID"

or

go run main.go download "https://www.instagram.com/p/POST_ID/"

The downloaded videos will be saved in a "video-downloader" folder in your user's Downloads directory.

## Configuration

The download directory is set to `~/Downloads/video-downloader` by default. If you want to change this, modify the `createDownloadDirectory` function in `main.go`.

## Limitations

- Instagram downloader only works for public posts
- The tool may break if YouTube or Instagram change their page structure
- Downloading content may violate the terms of service of YouTube and Instagram
- No other social media 

## Contributing

Contributions to the Video Downloader CLI are welcome. Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE