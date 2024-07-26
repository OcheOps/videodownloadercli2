package main

import (
	"fmt"
	"io"
	neturl "net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"encoding/json"
    "net/http"
    
    "regexp"
    "time"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
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
    client := youtube.Client{}

    video, err := client.GetVideo(url)
    if err != nil {
        return fmt.Errorf("error getting video info: %w", err)
    }

    formats := video.Formats.WithAudioChannels()
    sort.Slice(formats, func(i, j int) bool {
        return formats[i].Quality > formats[j].Quality
    })

    var format *youtube.Format
    for _, f := range formats {
        if f.Quality == "hd720" || f.Quality == "medium" {
            format = &f
            break
        }
    }

    if format == nil {
        return fmt.Errorf("no suitable video format found")
    }

    stream, size, err := client.GetStream(video, format)
    if err != nil {
        return fmt.Errorf("error getting video stream: %w", err)
    }
    defer stream.Close()

    fileName := sanitizeFileName(video.Title) + ".mp4"
    filePath := filepath.Join(downloadDir, fileName)

    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("error creating file: %w", err)
    }
    defer file.Close()

    fmt.Printf("Downloading: %s\n", video.Title)
    fmt.Printf("Size: %.2f MB\n", float64(size)/(1024*1024))

    _, err = io.Copy(file, io.TeeReader(stream, &ProgressWriter{Total: size}))
    if err != nil {
        return fmt.Errorf("error downloading video: %w", err)
    }

    fmt.Printf("\nVideo downloaded: %s\n", filePath)
    return nil
}

func sanitizeFileName(fileName string) string {
    fileName = strings.Map(func(r rune) rune {
        if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
            return -1
        }
        return r
    }, fileName)

    fileName = strings.TrimSpace(fileName)

    if len(fileName) > 200 {
        fileName = fileName[:200]
    }

    return fileName
}

type ProgressWriter struct {
	Total int64
	Current int64
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Current += int64(n)
	pw.PrintProgress()
	return n, nil
}

func (pw *ProgressWriter) PrintProgress() {
	if pw.Total == 0 {
		return
	}
	percentage := float64(pw.Current) / float64(pw.Total) * 100
	fmt.Printf("\rProgress: %.2f%%", percentage)
}

func downloadInstagramVideo(url, downloadDir string) error {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error fetching Instagram page: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("error reading Instagram page: %w", err)
    }

    re := regexp.MustCompile(`<script type="text/javascript">window\._sharedData = (.+);</script>`)
    matches := re.FindSubmatch(body)
    if len(matches) < 2 {
        return fmt.Errorf("could not find shared data in Instagram page")
    }

    var data map[string]interface{}
    err = json.Unmarshal(matches[1], &data)
    if err != nil {
        return fmt.Errorf("error parsing JSON data: %w", err)
    }

    entryData, ok := data["entry_data"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("could not find entry_data in JSON")
    }

    postPage, ok := entryData["PostPage"].([]interface{})
    if !ok || len(postPage) == 0 {
        return fmt.Errorf("could not find PostPage in JSON")
    }

    post, ok := postPage[0].(map[string]interface{})
    if !ok {
        return fmt.Errorf("could not find post data in JSON")
    }

    graphql, ok := post["graphql"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("could not find graphql data in JSON")
    }

    shortcodeMedia, ok := graphql["shortcode_media"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("could not find shortcode_media in JSON")
    }

    videoUrl, ok := shortcodeMedia["video_url"].(string)
    if !ok {
        return fmt.Errorf("could not find video_url in JSON")
    }

    videoResp, err := http.Get(videoUrl)
    if err != nil {
        return fmt.Errorf("error downloading video: %w", err)
    }
    defer videoResp.Body.Close()

    fileName := fmt.Sprintf("instagram_video_%d.mp4", time.Now().Unix())
    filePath := filepath.Join(downloadDir, fileName)

    out, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("error creating output file: %w", err)
    }
    defer out.Close()

    fmt.Printf("Downloading Instagram video\n")
    fmt.Printf("Size: %.2f MB\n", float64(videoResp.ContentLength)/(1024*1024))

    _, err = io.Copy(out, io.TeeReader(videoResp.Body, &ProgressWriter{Total: videoResp.ContentLength}))
    if err != nil {
        return fmt.Errorf("error saving video: %w", err)
    }

    fmt.Printf("\nVideo downloaded: %s\n", filePath)
    return nil
}
//TODO Titok video and vinemo video and pixel formats
