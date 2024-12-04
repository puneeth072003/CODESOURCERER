package generators

import (
	"fmt"
	"log"
	"os/exec"
)

// Function to commit new files to a given branch
func CommitNewFilesToBranch(newBranch string) error {
	// Check out to the new branch
	cmd := exec.Command("git", "checkout", newBranch)
	cmd.Dir = "/path/to/your/repo" // Set the directory to your local repository path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error checking out branch %s: %w", newBranch, err)
	}

	// Add the generated test files to the staging area
	cmd = exec.Command("git", "add", ".") // Add all the new files
	cmd.Dir = "/path/to/your/repo"        // Set the directory to your local repository path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adding files to staging area: %w", err)
	}

	// Commit the changes with a commit message
	commitMessage := "Add generated test files"
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	cmd.Dir = "/path/to/your/repo" // Set the directory to your local repository path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error committing files: %w", err)
	}

	// Optionally, push the changes to the remote repository
	// If you want to push the branch after committing
	cmd = exec.Command("git", "push", "origin", newBranch)
	cmd.Dir = "/path/to/your/repo" // Set the directory to your local repository path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pushing changes to remote repository: %w", err)
	}

	log.Printf("Successfully committed new test files to branch: %s", newBranch)
	return nil
}
