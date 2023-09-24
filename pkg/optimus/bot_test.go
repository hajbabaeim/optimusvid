package optimus

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Set the environment variable for the test
	os.Setenv("TELEGRAM_API", "test_api_key")

	optimus := Init()
	if optimus == nil {
		t.Fatal("Expected optimus to be initialized, but got nil")
	}

	if optimus.Bot == nil {
		t.Fatal("Expected bot to be initialized, but got nil")
	}

	if optimus.Format != "mp3" {
		t.Errorf("Expected default format to be mp3, but got %s", optimus.Format)
	}

	if optimus.Quality != "128k" {
		t.Errorf("Expected default quality to be 128k, but got %s", optimus.Quality)
	}
}

func TestExtractAudioFromVideo(t *testing.T) {
	optimus := &Optimus{Format: "mp3", Quality: "128k"}

	// You might need to create a temporary input file for testing
	inputPath := "test_input_path"
	outputPath := "test_output_path"
	audioCodec := "libmp3lame"
	audioBitrate := "128k"

	_, err := optimus.ExtractAudioFromVideo(inputPath, outputPath, audioCodec, audioBitrate)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check if the output file is created and clean up
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	} else {
		os.Remove(outputPath)
	}
}

func TestGetBotDescription(t *testing.T) {
	optimus := &Optimus{}
	description := optimus.GetBotDescription()
	if description == "" {
		t.Error("Expected a non-empty bot description")
	}
}
