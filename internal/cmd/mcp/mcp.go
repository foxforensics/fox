package mcp

import (
	"context"
	"log"
	"strings"

	res "github.com/cuhsat/fox/v4/internal"
	cli "github.com/cuhsat/fox/v4/internal/cmd"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Usage = strings.TrimSpace(`
Loads MCP server (blocking).

fox mcp [FLAGS...]

Examples:
  $ fox mcp
`)

type Mcp struct {
}

func (cmd *Mcp) Run(cli *cli.Globals) error {
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

		sb := new(strings.Builder)

		cli.Stdout = sb
		cli.NoFile = true
		cli.NoLine = true
		cli.NoColor = true
		cli.NoPretty = true
		cli.NoStrict = true

		mode := hunt.Hunt{
			Paths: []string{path},
		}

		err = mode.Run(cli)

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
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
