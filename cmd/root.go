package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "video-downloader",
	Short: "A CLI program to download Instagram, TikTok, and YouTube videos.",
	Long: `This program allows you to download videos from Instagram, TikTok, and YouTube. 
To use the program, simply run the following command:

video-downloader <platform> <video_url>

Supported platforms: instagram, tiktok, youtube`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		platform := args[0]
		videoURL := args[1]

		switch platform {
		case "instagram":
			downloadInstagramVideo(videoURL)
		case "tiktok":
			downloadTikTokVideo(videoURL)
		case "youtube":
			// Implement YouTube video download logic here
			fmt.Println("Downloading YouTube video is not yet implemented.")
		default:
			fmt.Println("Unsupported platform. Please use 'instagram', 'tiktok', or 'youtube'.")
		}
	},
}

func downloadInstagramVideo(videoURL string) {
	// Get the Instagram video ID.
	videoID := strings.Split(videoURL, "/")[4]

	// Download the Instagram video using the `youtube-dl` command-line tool.
	cmd := exec.Command("youtube-dl", "-f", "mp4", videoID)

	// Create a new folder named `ochevideos` if it does not already exist.
	err := os.MkdirAll("ochevideos", 0755)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Set the current working directory to the `ochevideos` folder.
	err = os.Chdir("ochevideos")
	if err != nil {
		fmt.Println("Error changing working directory:", err)
		return
	}

	// Execute the `youtube-dl` command.
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Video downloaded successfully!")
}

func downloadTikTokVideo(videoURL string) {
	// Get the TikTok video ID.
	videoID := strings.Split(videoURL, "/")[3]

	// Download the TikTok video using the `tiktok-video-downloader` command-line tool.
	cmd := exec.Command("tiktok-video-downloader", videoID)

	// Create a new folder named `ochevideos` if it does not already exist.
	err := os.MkdirAll("ochevideos", 0755)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Set the current working directory to the `ochevideos` folder.
	err = os.Chdir("ochevideos")
	if err != nil {
		fmt.Println("Error changing working directory:", err)
		return
	}

	// Execute the `tiktok-video-downloader` command.
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error downloading video:", err)
		return
	}

	fmt.Println("Video downloaded successfully!")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
