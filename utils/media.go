package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/404name/termui-demo/global"
	"github.com/404name/termui-demo/resource"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

func CovertImg(base64Img string) image.Image {
	image, _, err := image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Img)))
	if err != nil {
		log.Fatalf("failed to decode gopher image: %v", err)
		return nil
	}
	return image
}

// 读取不到默认返回LOGO
func LoadImg(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		global.LOG.Errorln(err.Error())
		return CovertImg(resource.LOGO)
	}
	// decode图片
	m, err := png.Decode(f)
	if err != nil {

		global.LOG.Errorln(err.Error())
		return CovertImg(resource.LOGO)

	}
	return m
}

func LoadAudio(path string) oto.Player {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	samplingRate := 16000

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)
	return player
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	//player.Play()

	// We can wait for the sound to finish playing using something like this
	// for player.IsPlaying() {
	// 	print("11")
	// 	//time.Sleep(time.Millisecond)
	// }

	// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
	// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
	// if err != nil{
	//     panic("player.Seek failed: " + err.Error())
	// }
	// println("Player is now at position:", newPos)
	// player.Play()

	// If you don't want the player/sound anymore simply close
	// err = player.Close()
	// if err != nil {
	// 	panic("player.Close failed: " + err.Error())
	// }
}
