package audio

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

func CreateMp3(message string) string {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	//

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-GB",
			Name:         "en-GB-Wavenet-C",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  1.2,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	filename := "output-" + time.Now().Format("2006-01-02-15-04-05") + ".mp3"

	dir := filepath.Join(".", "audio")
	os.MkdirAll(dir, os.ModePerm)

	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, resp.AudioContent, 0644)

	return filepath
}

func DeleteAudioFolder() error {
	if _, err := os.Stat("audio"); err == nil {
		// audio directory exists, delete it
		err = os.RemoveAll("audio")
		if err != nil {
			return err
		}
	}
	return nil
}
