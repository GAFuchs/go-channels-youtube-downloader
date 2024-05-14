package main

import (
	"io"
	"os"

	"github.com/kkdai/youtube/v2"
)

func YoutubeDownloadClient(videoID string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}
	defer stream.Close()

	file, err := os.Create("videos/" + videoID + ".mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}

	return nil
}
