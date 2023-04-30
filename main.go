package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/d3code/github-action-commit-workflow-changes/api"
)

func main() {
	println("Finding some lols...")

	reddit := api.GetReddit()
	var imageUrl string

	posts := reddit.Data.Children
	for _, post := range posts {

		fmt.Println(post.Data.Title)
		fmt.Println(post.Data.Url)

		if strings.HasSuffix(post.Data.Url, ".jpg") ||
			strings.HasSuffix(post.Data.Url, ".gif") ||
			strings.HasSuffix(post.Data.Url, ".png") {

			imageUrl = post.Data.Url
			break
		}
	}

	if imageUrl == "" {
		fmt.Println("No image found")
		return
	}

	fmt.Println("Found image: ", imageUrl)

	workspace := os.Getenv("GITHUB_WORKSPACE")
	os.Chdir(workspace)

	readme, err := os.ReadFile("README.md")
	if err != nil {
		panic(err.Error())
	}

	// Find the line with the image
	lines := strings.Split(string(readme), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "![Reddit]") {
			lines[i] = fmt.Sprintf("![%s](%s)", "Reddit", imageUrl)
			break
		}
	}

	readme = []byte(strings.Join(lines, "\n"))

	err = os.WriteFile("README.md", readme, 0644)
	if err != nil {
		panic(err.Error())
	}

	cmdAdd := exec.Command("git", "add", "README.md")
	cmdAdd.Stdout = os.Stdout
	cmdAdd.Run()

	cmdCommit := exec.Command("git", "commit", "-m", "'Change README.md'")
	cmdCommit.Stdout = os.Stdout
	cmdCommit.Run()

	cmdPush := exec.Command("git", "push")
	cmdPush.Stdout = os.Stdout
	cmdPush.Run()

	println("Done!")
}
