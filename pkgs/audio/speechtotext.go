package audio

import (
	"context"
	"encoding/binary"
	"fmt"
	"go_openai_cli/pkgs/config"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	// TODO update to v2
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/eiannone/keyboard"
	"github.com/gordonklaus/portaudio"
	"github.com/jbenet/goprocess"
)

func Transcribe() (string, error) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return "", err
	}
	defer client.Close()

	// Read an audio file.
	data, err := os.ReadFile("test.wav")
	if err != nil {
		fmt.Printf("Failed to read audio file: %v\n", err)
		return "", err
	}

	// Send a request to the API.
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		fmt.Printf("Failed to recognize: %v\n", err)
		return "", err
	}

	// Process the response.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("Transcript: %s\n", alt.Transcript)
			return alt.Transcript, nil
		}
	}
	return "", nil
}

func Record() error {
	fmt.Println("Recording. Press Space or Enter to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	f, err := os.Create(filepath.Join(config.GetDataPath(), "record.wav"))
	if err != nil {
		return err
	}
	defer func() {
		chk(f.Close())
	}()

	if err := portaudio.Initialize(); err != nil {
		return err
	}
	time.Sleep(1)
	defer portaudio.Terminate()
	in := make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return err
	}

	// Write WAV file header
	writeWAVHeader(f, 1, 16000, 16) // 1 channel, 16000 sample rate, 16-bit

	// Create a process to listen for keypress events
	proc := goprocess.WithParent(goprocess.Background())

	stopRecording := false

	proc.Go(func(worker goprocess.Process) {
		err := keyboard.Open()
		if err != nil {
			fmt.Println("Error opening keyboard input:", err)
			return
		}
		defer keyboard.Close()

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println("Error reading key:", err)
				break
			}
			if key == keyboard.KeySpace || key == keyboard.KeyEnter {
				stopRecording = true
				keyboard.Close()
				break
			}
		}
	})

loop:
	for {
		if stopRecording {
			break loop
		}

		if err := stream.Read(); err != nil {
			return err
		}

		for _, sample := range in {
			if err := binary.Write(f, binary.LittleEndian, sample); err != nil {
				return err
			}
		}

		select {
		case <-sig:
			break loop
		default:
		}
	}

	// Stop the keypress listening process
	proc.Close()

	if err := stream.Stop(); err != nil {
		return err
	}

	fmt.Println("Stopped recording")
	return nil
}

func writeWAVHeader(file *os.File, numChannels, sampleRate, bitsPerSample int) {
	chunkSize := 36
	subChunk2Size := chunkSize - 36
	audioFormat := 1

	binary.Write(file, binary.LittleEndian, []byte("RIFF"))
	binary.Write(file, binary.LittleEndian, uint32(4+24+subChunk2Size))
	binary.Write(file, binary.LittleEndian, []byte("WAVE"))
	binary.Write(file, binary.LittleEndian, []byte("fmt "))
	binary.Write(file, binary.LittleEndian, uint32(16))
	binary.Write(file, binary.LittleEndian, uint16(audioFormat))
	binary.Write(file, binary.LittleEndian, uint16(numChannels))
	binary.Write(file, binary.LittleEndian, uint32(sampleRate))
	binary.Write(file, binary.LittleEndian, uint32(sampleRate*numChannels*bitsPerSample/8))
	binary.Write(file, binary.LittleEndian, uint16(numChannels*bitsPerSample/8))
	binary.Write(file, binary.LittleEndian, uint16(bitsPerSample))
	binary.Write(file, binary.LittleEndian, []byte("data"))
	binary.Write(file, binary.LittleEndian, uint32(subChunk2Size))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
