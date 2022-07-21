package filesui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	shell "github.com/ipfs/go-ipfs-api"
	files "github.com/ipfs/go-ipfs-files"
)

// getIpfsFiles returns all files in the local IPFS node's root directory.
func getIpfsFiles() []list.Item {
	// Opens a 'shell' to the local IPFS node
	sh := shell.NewShell("localhost:5001")
	ctx := context.TODO()
	filesList, err := sh.FilesLs(ctx, "")

	if err != nil {
		err = fmt.Errorf("error: %s", err)
		fmt.Print(err)
	}

	var items []list.Item

	for i := 0; i < len(filesList); i++ {
		item := item{
			fileName: filesList[i].Name,
		}

		items = append(items, item)
	}

	return items

}

// addFile returns the CID of a file added to IPFS
func addFile(ctx context.Context, path string, fileName string) string {
	// Opens a 'shell' to the local IPFS node
	sh := shell.NewShell("localhost:5001")
	file, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
	}
	fileReader := files.NewReaderFile(file)

	cid, err := sh.Add(fileReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	// When referencing files added to ipfs, you can get them at /ipfs/<cid>
	ipfsPath := fmt.Sprintf("/ipfs/%s", cid)

	// Human-readable name in the / directory of the MFS
	// TODO: allow the user to select which directory this file will go to
	newPath := fmt.Sprintf("/%s", fileName)

	// TODO: Look into utilizing FilesWrite to write directly to IPFS' MFS
	// FilesCp copies the IPFS file we just added to the MFS.
	err = sh.FilesCp(ctx, ipfsPath, newPath)
	if err != nil {
		err = fmt.Errorf("error: %s", err)
		fmt.Print(err)
	}

	return cid
}

type peerInfo []shell.SwarmConnInfo

func getIpfsStats() tea.Msg {
	// Opens a 'shell' to the local IPFS node
	sh := shell.NewShell("localhost:5001")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	swarm, err := sh.SwarmPeers(ctx)
	if err != nil {
		err = fmt.Errorf("error, %s", err)
		fmt.Print(err)
		os.Exit(1)
	}

	return peerInfo(swarm.Peers)
}

func (m Model) getSwarmPeers() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return getIpfsStats()
	})

}
