package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

const systemPrompt = "You are a helpful sysadmin specialized in Nix. \n" +
	"You are good at understanding the Linux environment a user wants \n" +
	"and knowing what packages from Nix are required."

const userPrompt = "Given the following user requirement for a Nixos configuration, alter the existing Nix files to fill the requirements. \n" +
	"Feel free to think about the packages the user wants and how to get them from Nix. The NixOS configuration file should be returned as a \n" +
	"Markdown code block, formatted like this:\n\n" +
	"```nix\n" +
	"# Your NixOS configuration file here\n" +
	"```\n" +
	"The user request is:\n" +
	"%s\n" +
	"The existing Nix files are:\n" +
	"%s\n"

var modelName = "gpt-4o-mini"

// var modelName = "o1-mini"

var openAIClient *azopenai.Client

func getClient() *azopenai.Client {
	if openAIClient == nil {
		credential := azcore.NewKeyCredential(OPENAI_API_KEY)
		client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", credential, nil)
		if err != nil {
			panic(err)
		}
		openAIClient = client
	}
	return openAIClient
}

func nixFilesToMarkdown(files map[string]string) string {
	var contents string
	for filename, content := range files {
		contents += fmt.Sprintf("### %s\n\n```nix\n%s\n```\n\n", filename, content)
	}
	return contents
}

func runPrompt(input string, files map[string]string) string {
	// Implementation for running the prompt
	client := getClient()

	markdownFiles := nixFilesToMarkdown(files)
	prompt := fmt.Sprintf(userPrompt, input, markdownFiles)

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: azopenai.NewChatRequestSystemMessageContent(systemPrompt)},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(prompt)},
	}
	schemaMap := map[string]interface{}{
		"type": "object",
		"additionalProperties": map[string]interface{}{
			"type":        "string",
			"description": "Content of the file.",
		},
		"description": "A mapping from filename strings to file content strings.",
	}
	schemaBytes, err := json.Marshal(schemaMap)
	schema := azopenai.ChatCompletionsJSONSchemaResponseFormatJSONSchema{
		Name:        to.Ptr("nixos-file-output"),
		Description: to.Ptr("A mapping of .nix file names to their contents, which will all be deployed together as a NixOS configuration."),
		Schema:      schemaBytes,
		Strict:      to.Ptr(true),
	}
	format := azopenai.ChatCompletionsJSONSchemaResponseFormat{
		to.Ptr("json_schema"),
		&schema,
	}

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &modelName,
		ResponseFormat: schema,
	}, nil)

	if err != nil {
		panic(err)
	}

	nixFile := extractNixFile(*resp.ChatCompletions.Choices[0].Message.Content)
	return nixFile
}

/*
Get only the content from a string within
```nix
````
*/
func extractNixFile(str string) string {

	start := "```nix\n"
	end := "\n```"
	startIndex := strings.Index(str, start)
	endIndex := strings.Index(str, end)
	if startIndex == -1 || endIndex == -1 {
		panic("Failed to extract NixOS configuration.")
	}
	nixConfig := str[startIndex+len(start) : endIndex]
	return nixConfig
}
