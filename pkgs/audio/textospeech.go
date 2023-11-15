package audio

import (
	"context"
	"encoding/json"
	"go_openai_cli/pkgs/config"
	"go_openai_cli/pkgs/openai"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bytes"
	"fmt"
	"net/http"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

// Voice represents the available voice options
type Voice string

const (
	Alloy   Voice = "alloy"
	Echo    Voice = "echo"
	Fable   Voice = "fable"
	Onyx    Voice = "onyx"
	Nova    Voice = "nova"
	Shimmer Voice = "shimmer"
)

func CreateMp3(message string) string {
	var bytes []byte
	var err error

	//Defaults to google cloud text to speech. Because it's faster and cheaper. Though not as good.
	if false {
		bytes, err = openaiSynth(message, Voice("fable"))
	} else {
		bytes, err = gooleCloudTextSynthesize(message)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	filename := time.Now().Format("2006-01-02-15-04-05") + ".mp3"
	var audioFolder = filepath.Join(config.GetDataPath(), "audio")
	sysModelAudioDir := filepath.Join(audioFolder, openai.GetSystemModel().String())
	os.MkdirAll(sysModelAudioDir, os.ModePerm)

	filepath := filepath.Join(sysModelAudioDir, filename)
	err = os.WriteFile(filepath, bytes, 0644)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return filepath
}

func DeleteAudioFolder() error {
	var audioFolder = filepath.Join(config.GetDataPath(), "audio")
	if _, err := os.Stat(audioFolder); err == nil {
		// audio directory exists, delete it
		err = os.RemoveAll(audioFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

// Is better than gooleCloudTextSynthesize. But slower and more expensive.
func openaiSynth(input string, voice Voice) ([]byte, error) {
	fmt.Println(input + "\n" + string(voice))
	url := "https://api.openai.com/v1/audio/speech"
	// Replace newline characters with "\n" in the input
	input = strings.ReplaceAll(input, "\n", "\\n")

	requestData := struct {
		Model string `json:"model"`
		Input string `json:"input"`
		Voice string `json:"voice"`
	}{
		Model: "tts-1",
		Input: input,
		Voice: string(voice),
	}

	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestDataBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.Cfg.OpenAIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	println("success!")
	return body, nil
}

// Is faster and cheaper than openaiSynth. But not as good.
func gooleCloudTextSynthesize(message string) ([]byte, error) {

	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var voiceParams *texttospeechpb.VoiceSelectionParams
	var audioConfig *texttospeechpb.AudioConfig

	if openai.GetSystemModel() == openai.Alfred {

		voiceParams = &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-GB",
			Name:         "en-GB-Standard-D",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		}

		audioConfig = &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  1.2,
		}
	} else {
		voiceParams = &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-GB",
			Name:         "en-GB-Wavenet-C",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		}

		audioConfig = &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  1.2,
		}
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},

		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice:       voiceParams,
		AudioConfig: audioConfig,
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.AudioContent, nil
}
