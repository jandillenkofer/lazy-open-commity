package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

const topCommitMessageCount = 5

const promptTemplate = `You are a helpful commit message generator.
Generate me an array of the top %d commit message recommendations in json format.
Only one commit is picked from the returned commit message list.
Only answer with the json! Nothing else! No markdown!

# Rules for writing commit messages:
- Use the conventional commit specification.
- Capitalized, short (50 chars or less) summary (title)
- More detailed description, if necessary. Bullet points can be used, too. Typically a hyphen or asterisk is used for the bullet, preceded by a single space, with blank lines in between
- What was the motivation for the change?
- How does it differ from the previous implementation?
- You MUST insert newlines in the description to ensure a maximum of 72 characters per line.
- Write your commit message in the imperative: "Fix bug" and not "Fixed bug" or "Fixes bug." 

# Example output format:
Here is an example of the output format:
[
	{
		"title": "feat: add user authentication flow",
		"description": "Implement login, registration, and JWT-based session management for secure user access."
	}
]

`

type AIModel string

const (
	gpt41 = AIModel("github-copilot/gpt-4.1")
	gpt4o = AIModel("github-copilot/gpt-4o")
)

func runCommandAndGetOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return string(output), err
}

func promptAIAgent(prompt string, model AIModel) (string, error) {
	return runCommandAndGetOutput("opencode", "-m", string(model), "run", prompt)
}

type CommitMessage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	checkForStagedChanges()

	prompt := fmt.Sprintf(promptTemplate, topCommitMessageCount)
	prompt = extendPromptWithCodeChanges(prompt)
	prompt = extendPromptWithCommitHistory(prompt)

	jsonResponse, err := promptAIAgent(prompt, gpt4o)
	if err != nil {
		fmt.Println("Error running 'opencode run':", err)
		os.Exit(1)
	}

	var commitMessages []CommitMessage
	err = json.Unmarshal([]byte(jsonResponse), &commitMessages)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	if len(commitMessages) == 0 {
		fmt.Println("No commit messages generated.")
		os.Exit(1)
	}

	selectedCommitMessage := selectCommitMessage(commitMessages)

	output := gitCommit(selectedCommitMessage)
	fmt.Println(output)
}

func gitCommit(selectedCommitMessage CommitMessage) string {
	output, err := runCommandAndGetOutput("git", "commit", "-m", selectedCommitMessage.Title, "-m", selectedCommitMessage.Description)
	if err != nil {
		fmt.Println("Error running 'git commit':", err)
		os.Exit(1)
	}
	return output
}

func selectCommitMessage(commitMessages []CommitMessage) CommitMessage {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "→ {{ .Title | cyan }}",
		Inactive: "  {{ .Title }}",
		Selected: "✔ {{ .Title | green }}",
		Details:  "{{ .Description }}",
	}

	selectPrompt := promptui.Select{
		Label:     "Select commit message",
		Items:     commitMessages,
		Templates: templates,
		Size:      len(commitMessages),
		Searcher: func(input string, index int) bool {
			commit := commitMessages[index]
			name := commit.Title + " " + commit.Description
			inputLower := strings.ToLower(input)
			nameLower := strings.ToLower(name)
			return strings.Contains(nameLower, inputLower)
		},
		StartInSearchMode: true,
	}

	idx, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		os.Exit(1)
	}
	topCommitMessage := commitMessages[idx]
	return topCommitMessage
}

func checkForStagedChanges() {
	stagedChanges, err := runCommandAndGetOutput("git", "diff", "--cached", "--stat")
	if err != nil {
		fmt.Println("Error running 'git diff --cached --stat':", err)
		os.Exit(1)
	}
	if stagedChanges == "" {
		fmt.Println("No staged changes found. Please stage your changes before running this tool.")
		os.Exit(0)
	}
}

func extendPromptWithCommitHistory(prompt string) string {
	prompt += "# List of the last commits (additional context):\n"
	lastCommits, err := runCommandAndGetOutput("git", "log", "-n", "5", "--pretty=oneline", "--abbrev-commit")
	if err != nil {
		fmt.Println("Error running 'git log -n 5 --pretty=oneline --abbrev-commit':", err)
		os.Exit(1)
	}
	prompt += "```console\n" + "$ git log -n 5 --pretty=oneline --abbrev-commit\n" + lastCommits + "```\n"
	return prompt
}

func extendPromptWithCodeChanges(prompt string) string {
	prompt += "# Diff to analyze:\n"
	diffStats, err := runCommandAndGetOutput("git", "diff", "--cached", "--stat")
	if err != nil {
		fmt.Println("Error running 'git diff --cached --stat':", err)
		os.Exit(1)
	}
	prompt += "```console\n" + "$ git diff --cached --stat\n" + diffStats + "```\n"
	diff, err := runCommandAndGetOutput("git", "diff", "--cached")
	if err != nil {
		fmt.Println("Error running 'git diff --cached':", err)
		os.Exit(1)
	}
	prompt += "```console\n" + "$ git diff --cached\n" + diff + "```\n"
	return prompt
}
