package openai

import (
	"fmt"
	"regexp"
	"testing"
)

func TestHelloName(t *testing.T) {

	expected_prompt := "3D render of a cute tropical fish in an aquarium on a dark blue background, digital art"
	str := `\"IMAGE_PROMPT: [3D render of a cute tropical fish in an aquarium on a dark blue background, digital art]\""`
	re := regexp.MustCompile(`IMAGE_PROMPT:\s*\[([^]]*)\]`)
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		if match[1] == expected_prompt {

			t.Log(match[1])
			fmt.Println(match[1])
		}

	} else {
		t.Fatalf(`Failed to find expected prompt`)
	}
}
