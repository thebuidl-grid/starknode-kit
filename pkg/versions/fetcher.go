package versions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

type ClientVersions struct {
	Geth       string `json:"geth"`
	Reth       string `json:"reth"`
	Lighthouse string `json:"lighthouse"`
	Prysm      string `json:"prysm"`
	Juno       string `json:"juno"`
}

// FetchLatestVersions fetches the latest versions for all clients from GitHub
func FetchLatestVersions() (*ClientVersions, error) {
	versions := &ClientVersions{}

	// Define GitHub repositories
	repos := map[string]string{
		"geth":       "ethereum/go-ethereum",
		"reth":       "paradigmxyz/reth",
		"lighthouse": "sigp/lighthouse",
		"prysm":      "prysmaticlabs/prysm",
		"juno":       "NethermindEth/juno",
	}

	// Fetch versions concurrently
	type result struct {
		client  string
		version string
		err     error
	}

	results := make(chan result, len(repos))

	for clientName, repo := range repos {
		go func(client, repo string) {
			version, err := fetchGitHubRelease(client, repo)
			results <- result{client, version, err}
		}(clientName, repo)
	}

	// Collect results
	var errors []string
	for i := 0; i < len(repos); i++ {
		res := <-results
		if res.err != nil {
			return nil, res.err
		} else {
			switch res.client {
			case "geth":
				versions.Geth = res.version
			case "reth":
				versions.Reth = res.version
			case "lighthouse":
				versions.Lighthouse = res.version
			case "prysm":
				versions.Prysm = res.version
			case "juno":
				versions.Juno = res.version
			}
		}
	}

	if len(errors) > 0 {
		fmt.Printf("Warning: Some version fetches failed, using fallback versions: %s\n", strings.Join(errors, ", "))
	}

	return versions, nil
}

// fetchGitHubRelease fetches the latest release tag from a GitHub repository
func fetchGitHubRelease(clientName, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Clean version string (remove 'v' prefix if present)
	version := strings.TrimPrefix(release.TagName, "v")

	// Special handling for different clients
	switch clientName {
	case "prysm":
		// Prysm uses different versioning format
		if release.Name != "" {
			return release.Name, nil
		}
	case "juno":
		// Keep 'v' prefix for Juno as it's expected in the codebase
		return release.TagName, nil
	}

	return version, nil
}

// FetchLatestGethVersion fetches the latest Geth version from GitHub
func FetchLatestGethVersion() (string, error) {
	return fetchGitHubRelease("geth", "ethereum/go-ethereum")
}

// FetchLatestRethVersion fetches the latest Reth version from GitHub
func FetchLatestRethVersion() (string, error) {
	return fetchGitHubRelease("reth", "paradigmxyz/reth")
}

// FetchLatestLighthouseVersion fetches the latest Lighthouse version from GitHub
func FetchLatestLighthouseVersion() (string, error) {
	return fetchGitHubRelease("lighthouse", "sigp/lighthouse")
}

// FetchLatestPrysmVersion fetches the latest Prysm version from GitHub
func FetchLatestPrysmVersion() (string, error) {
	return fetchGitHubRelease("prysm", "prysmaticlabs/prysm")
}

// FetchLatestJunoVersion fetches the latest Juno version from GitHub
func FetchLatestJunoVersion() (string, error) {
	return fetchGitHubRelease("juno", "NethermindEth/juno")
}
func FetchLatestStarknetValidatorVersion() (string, error) {
	return fetchGitHubRelease("starknet-staking-v2", "NethermindEth/starknet-staking-v2")
}

// getLatestPrysmVersion returns the hardcoded Prysm version as fallback
// Prysm handles versioning differently through their script
func getLatestPrysmVersion() string {
	return "latest" // Prysm script auto-downloads latest
}
