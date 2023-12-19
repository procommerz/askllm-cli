package main

import (
	"context"
	"errors"
	"fmt"
	color "github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"os"
	"strconv"
	"strings"
)

var openAIKey = ""  // Will be filled with the API key from ~/.askllm (openai_key)
var systemInfo = "" // Will be filled with the system info from ~/.askllm (system_info)
var maxTokens = 400

const defaultModel = "gpt-4-1106-preview"

var prependSystemInfo bool = false
var fileAnalysis bool = false
var fileAnalysisFilename string = ""

var cli_red = color.New(color.FgRed)

func main() {
	// Try to read the settings file in the home directory:
	homeDirname, _ := os.UserHomeDir()
	settingsContents, err := os.ReadFile(homeDirname + "/.askllm")

	// Assist in creating the settings file
	if err != nil {
		cli_red.Println("You have to have an ~/.askllm file with your OpenAI API key in it, like 'openai_key=xxxxxxxxxxxxx' (", err, ")")
		// Try to create the file in the home directory:
		err = os.WriteFile(homeDirname+"/.askllm", []byte("openai_key=REPLACE_WITH_YOUR_ACTUAL_KEY\nmax_tokens=300\nsystem_info="), 0644)
		if err != nil {
			c2 := color.New(color.FgHiYellow)
			c2.Println("Tried to create a template ~/.askllm file, but failed due to error:", err)
		}

		return
	}

	// Read the settings file, as key/value pairs:
	settings := strings.Split(string(settingsContents), "\n")

	// Loop through the settings and find the key:
	for _, setting := range settings {
		// Split the setting into key/value:
		keyValue := strings.Split(setting, "=")

		// Check if the key is "openai_key":
		if keyValue[0] == "openai_key" {
			// Set the global variable:
			openAIKey = strings.TrimSpace(keyValue[1])
		} else if keyValue[0] == "system_info" {
			// Set the global variable:
			systemInfo = strings.TrimSpace(keyValue[1])
		} else if keyValue[0] == "max_tokens" {
			// Set the global variable, parsing to int
			conv_err := error(nil)
			maxTokens, conv_err = strconv.Atoi(strings.TrimSpace(keyValue[1]))
			if conv_err != nil {
				cli_red.Println("Error converting max_tokens to int:", conv_err)
				cli_red.Println("Using default value of 400, but please fix the value in the ~/.askllm file")
				maxTokens = 400
			}
		}
	}

	// Warn and exit if the API key is still set to stub value
	if openAIKey == "" || openAIKey == "REPLACE_WITH_YOUR_ACTUAL_KEY" {
		c := color.New(color.FgHiRed)
		c.Println("Replace the key in the ~/.askllm file with your actual OpenAI API key, that you can get here: https://platform.openai.com/api-keys")
		return
	}

	// Process command line arguments.
	// Discard the first argument (the program name):
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	// Check for special arguments,
	// if the first token is "-s" or "--system-info",
	// set prependSystemInfo to true:
	if args[0] == "-s" || args[0] == "--system" {
		prependSystemInfo = true
		args = args[1:]
	}

	// -f argument allows asking a question about a specific file,
	// it's contents will be fed to the OpenAI API as part of the prompt:
	if args[0] == "-f" || args[0] == "--file" {
		fileAnalysis = true
		fileAnalysisFilename = args[1]
		args = args[2:]
	}

	// -f argument allows asking a question about a specific file,
	// it's contents will be fed to the OpenAI API as part of the prompt:
	if args[0] == "-m" || args[0] == "--more" {
		maxTokens = maxTokens * 3
		args = args[1:]
	}

	if args[0] == "-h" || args[0] == "--help" {
		fmt.Println("Usage: askllm [-s] [-f filename] question. OpenAI API key is must be set in the ~/.askllm file")
		fmt.Println("  -s, --system			Prepend the question with system info from ~/.askllm")
		fmt.Println("  -f, --file			Add the file contents to the prompt for analysis")
		fmt.Println("  -m, --more			Triples the max_tokens value, for longer answers")
		fmt.Println("  -h, --help			Show this help")
		return
	}

	// The rest of the arguments are the prompt:
	prompt := strings.Join(args, " ")

	// Strip whitespace from the prompt:
	prompt = strings.TrimSpace(prompt)

	if prompt == "" || len(prompt) < 3 {
		fmt.Println("No question means no answer. Formulate your request as a question or a problem, don't enter gibberish.")
		c := color.New(color.FgMagenta)
		c.Println("Just think about all the wasted electricity that went into this request 😩")
		return
	}

	// send a streaming request:
	sendStreamingChatRequest(prompt)
}

func sendStreamingChatRequest(prompt string) {
	client := openai.NewClient(openAIKey)
	ctx := context.Background()

	// Prepare an empty message list:
	messages := []openai.ChatCompletionMessage{}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Answer the following question in a brief and informative way. Format the output for a bash-like console output.",
	})

	if prependSystemInfo == true {
		// Get the operating system name and version:
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemInfo,
		})
	}

	if fileAnalysis == true {
		// Read the fileAnalysisFilename file into a string:
		fileContents, err := os.ReadFile(fileAnalysisFilename)

		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// Get the filename part of the fileAnalysisFilename:
		filenameParts := strings.Split(fileAnalysisFilename, "/")
		onlyFilename := filenameParts[len(filenameParts)-1]

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "My question is about the file " + onlyFilename + ". Here's it's contents:\n",
		})

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: string(fileContents),
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	req := openai.ChatCompletionRequest{
		Model:     defaultModel,
		MaxTokens: maxTokens,
		Messages:  messages,
		Stream:    true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	c := color.New(color.FgYellow)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			//fmt.Println("Stream finished")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		c.Print(response.Choices[0].Delta.Content)
	}
}
