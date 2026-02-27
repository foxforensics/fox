package mcp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
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
	Port uint16 `arg:"" optional:"" default:"3001"`

	// internal
	addr string `kong:"-"`
}

func (cmd *Mcp) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	cmd.addr = fmt.Sprintf(":%d", cmd.Port)

	return nil
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
		server.WithToolCapabilities(true),
	)

	cmd.addHunt(cli, srv)

	if cli.Verbose > 0 {
		log.Println(fmt.Sprintf("mcp server started on %d", cmd.Port))
	}

	sse := server.NewStreamableHTTPServer(srv)

	if err := sse.Start(cmd.addr); err != nil {
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
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
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

		mode := hunt.Hunt{
			All:   true,
			Uniq:  true,
			Paths: []string{path},
		}

		pipe := bytes.NewBuffer(nil)

		cli.Stdout = pipe

		err = mode.Run(cli)

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if cli.Verbose > 1 {
			log.Println(fmt.Sprintf("tx: %s", pipe.String()))
		}

		return mcp.NewToolResultText(pipe.String()), nil
	})
}
