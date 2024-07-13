package main

import (
	"fmt"
	"os"
	"path/filepath"
	neturl "net/url"

	"github.com/spf13/cobra"
	"github.com/kkdai/youtube/v2"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "video-downloader",
		Short: "A video downloader for YouTube and Instagram",
		Long:  `A video downloader that can download videos from YouTube and Instagram using just the URL.`,
	}

	var downloadCmd = &cobra.Command{
		Use:   "download [URL]",
		Short: "Download a video from the given URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			fmt.Printf("Validating URL: %s\n", url)

			platform, err := validateURL(url)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			downloadDir, err := createDownloadDirectory()
			if err != nil {
				fmt.Println("Error creating download directory:", err)
				return
			}

			fmt.Printf("Downloading %s video to: %s\n", platform, downloadDir)

			switch platform {
			case "youtube":
				err = downloadYouTubeVideo(url, downloadDir)
			case "instagram":
				err = downloadInstagramVideo(url, downloadDir)
			}

			if err != nil {
				fmt.Println("Error downloading video:", err)
			} else {
				fmt.Println("Video downloaded successfully!")
			}
		},
	}

	rootCmd.AddCommand(downloadCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createDownloadDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	downloadDir := filepath.Join(homeDir, "Downloads", "video-downloader")
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return "", err
	}

	return downloadDir, nil
}

func validateURL(url string) (string, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	if host == "www.youtube.com" || host == "youtube.com" || host == "youtu.be" {
		return "youtube", nil
	} else if host == "www.instagram.com" || host == "instagram.com" {
		return "instagram", nil
	}

	return "", fmt.Errorf("unsupported URL: %s", url)
}

func downloadYouTubeVideo(url, downloadDir string) error {
	// TODO: Implement YouTube video downloading
	return fmt.Errorf("YouTube video downloading not implemented yet")
}

func downloadInstagramVideo(url, downloadDir string) error {
	// TODO: Implement Instagram video downloading
	return fmt.Errorf("Instagram video downloading not implemented yet")
}