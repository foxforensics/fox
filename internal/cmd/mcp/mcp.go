package mcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	res "github.com/cuhsat/fox/v4/internal"
	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
)

var Usage = strings.TrimSpace(`
Inits MCP server (blocks).

fox mcp [FLAGS...] [PORT]

Examples:
  $ fox mcp 8080
`)

type Mcp struct {
	Port uint16 `arg:"" optional:"" default:"3000"`
}

func (cmd *Mcp) Run(cli *cli.Globals) error {
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

	cmd.addHunt(cli, srv)

	if cli.Verbose > 0 {
		log.Println("mcp server started")
	}

	addr := fmt.Sprintf(":%d", cmd.Port)

	sse := server.NewStreamableHTTPServer(srv)

	if err := sse.Start(addr); err != nil {
		log.Fatalln(err)
	}

	if cli.Verbose > 0 {
		log.Println("mcp server stopped")
	}

	return nil
}

func (cmd *Mcp) addHunt(cli *cli.Globals, srv *server.MCPServer) {
	srv.AddTool(mcp.NewTool("hunt",
		mcp.WithDescription("Search for suspicious event logs in a file"),
		mcp.WithBoolean("all",
			mcp.Description("Look for all severities"),
		),
		mcp.WithBoolean("sort",
			mcp.Description("Sort events by timestamp"),
		),
		mcp.WithBoolean("uniq",
			mcp.Description("Filter unique events"),
		),
		mcp.WithString("path",
			mcp.Description("Path of the file to search in"),
			mcp.Required(),
		),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		all, err := request.RequireBool("all")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		sort, err := request.RequireBool("sort")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		uniq, err := request.RequireBool("uniq")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		path, err := request.RequireString("path")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if cli.Verbose > 1 {
			log.Println(fmt.Sprintf("rx: %s", path))
		}

		sb := new(strings.Builder)

		cli.Stdout = sb

		mode := hunt.Hunt{
			All:   all,
			Sort:  sort,
			Uniq:  uniq,
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
}
