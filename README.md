# askllm-cli

`askllm` is a simple CLI interface for asking questions to the OpenAI GPT API from command line. It uses the chat completion API to form responses and feels much faster than the official web interface.
If you need to ask a simple question, especially about something CLI-related, this tool offers a faster and more convenient workflow.

It requires you to have an OpenAI API key, from which you will need a developer subscription with openai.com

Privacy is a sensitive issue when it comes to personal AI requests. That's why this tool is designed as a simple
one-page app, that you can easily audit and compile yourself. It does not send any data to any server, except the OpenAI API.

## Installation with Homebrew

On Mac with homebrew, just run the following:

```bash
brew tap procommerz/askllm
brew install askllm
```

## Installation from the Repository

First `cd` to the cloned repositoty.

Use one of the file from dist/[your platform]/askllm and link it to /usr/local/bin/askllm:

```bash
ln -s ./dist/macos_arch64/askllm /usr/local/bin/askllm
```

...or compile it yourself with

```bash
go build -o askllm main.go
ln -s ./askllm /usr/local/bin/askllm
```

Run it once, to create the config file

```bash
askllm hello
```

Edit the config file in ~/.askllm and add your OpenAI API key

## Use It

```
askllm How to restart nginx?
```

## Command Line Usage

```
Usage: askllm [-s] [-f filename] question. OpenAI API key is must be set in the ~/.askllm file
-s, --system			Prepend the question with system info from ~/.askllm
-f, --file			Add the file contents to the prompt for analysis
-m, --more			Triples the max_tokens value, for longer answers
-h, --help			Show this help
```