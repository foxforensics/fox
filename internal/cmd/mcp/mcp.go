package mcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	res "github.com/cuhsat/fox/v4/internal"
	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
)

var Usage = strings.TrimSpace(`
Loads MCP server (blocks).

fox mcp [FLAGS...]

Examples:
  $ fox mcp
`)

type Mcp struct{}

func (cmd *Mcp) Run(cli *cli.Globals) error {
	sb := new(strings.Builder)

	// prepare
	cli.Stdout = sb
	cli.NoFile = true
	cli.NoLine = true
	cli.NoColor = true
	cli.NoPretty = true
	cli.NoStrict = true

	srv := server.NewMCPServer(
		"fox",
		res.Version,
		server.WithToolCapabilities(false),
	)

	srv.AddTool(mcp.NewTool("hunt",
		mcp.WithDescription("Search for suspicious event logs in a file"),
		mcp.WithString("path",
			mcp.Description("Path of the file to search in"),
			mcp.Required(),
		),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, err := request.RequireString("path")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if cli.Verbose > 1 {
			log.Println(fmt.Sprintf("rx: %s", path))
		}

		sb.Reset()

		mode := hunt.Hunt{
			Paths: []string{path},
		}

		err = mode.Run(cli)

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if cli.Verbose > 1 {
			log.Println(fmt.Sprintf("tx: %s", sb.String()))
		}

		return mcp.NewToolResultText(sb.String()), nil
	})

	if cli.Verbose > 0 {
		log.Println("mcp server started")
	}

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalln(err)
	}

	if cli.Verbose > 0 {
		log.Println("mcp server stopped")
	}

	return nil
}
