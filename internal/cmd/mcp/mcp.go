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

	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/stat"
	"github.com/cuhsat/fox/v4/internal/cmd/test"
	"github.com/cuhsat/fox/v4/internal/cmd/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Inits MCP server (blocks).

fox mcp [FLAGS...] [PORT]

Examples:
  $ fox mcp 8080
`)

type Mode interface {
	Run(_ *cli.Globals) error
}

type Mcp struct {
	Port uint16 `arg:"" optional:"" default:"3001"`

	// internal
	mcp *server.MCPServer `kong:"-"`
}

func (cmd *Mcp) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	cmd.mcp = server.NewMCPServer("fox", res.Version,
		server.WithToolCapabilities(true),
	)

	return nil
}

func (cmd *Mcp) Run(cli *cli.Globals) error {
	// prepare cli
	cli.NoFile = true
	cli.NoLine = true
	cli.NoColor = true
	cli.NoPretty = true
	cli.NoStrict = true

	// add tools
	cmd.addCat(cli)
	cmd.addHex(cli)
	cmd.addText(cli)
	cmd.addHash(cli)
	cmd.addStat(cli)
	cmd.addTest(cli)
	cmd.addDump(cli)
	cmd.addHunt(cli)

	if cli.Verbose > 0 {
		log.Println(fmt.Sprintf("mcp: started on port %d", cmd.Port))
	}

	return server.NewStreamableHTTPServer(cmd.mcp).Start(fmt.Sprintf(":%d", cmd.Port))
}

func (cmd *Mcp) addCat(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("cat",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("uniq", mcp.Description("")),
		mcp.WithNumber("dist", mcp.Description("")),
		mcp.WithString("regex", mcp.Description("")),
		mcp.WithNumber("context", mcp.Description("")),
		mcp.WithNumber("before", mcp.Description("")),
		mcp.WithNumber("after", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		cli.Regex = request.GetString("regex", "")

		return execute(cli, &cat.Cat{
			Uniq:    request.GetBool("uniq", false),
			Dist:    request.GetFloat("dist", 0),
			Context: uint(request.GetInt("context", 0)),
			Before:  uint(request.GetInt("before", 0)),
			After:   uint(request.GetInt("after", 0)),
			Paths:   paths,
		}), nil
	})
}

func (cmd *Mcp) addHex(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("hex",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("hexdump", mcp.Description("")),
		mcp.WithBoolean("xxd", mcp.Description("")),
		mcp.WithBoolean("raw", mcp.Description("")),
		mcp.WithBoolean("decimal", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &hex.Hex{
			Hexdump: request.GetBool("hexdump", false),
			Xxd:     request.GetBool("xxd", false),
			Raw:     request.GetBool("raw", false),
			Decimal: request.GetBool("decimal", false),
			Paths:   paths,
		}), nil
	})
}

func (cmd *Mcp) addText(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("text",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithNumber("min", mcp.Description("")),
		mcp.WithNumber("max", mcp.Description("")),
		mcp.WithBoolean("ascii", mcp.Description("")),
		mcp.WithBoolean("sort", mcp.Description("")),
		mcp.WithBoolean("decimal", mcp.Description("")),
		mcp.WithNumber("wtf", mcp.Description("")),
		mcp.WithArray("find", mcp.Description(""), mcp.Required()),
		mcp.WithBoolean("first", mcp.Description("")),
		mcp.WithBoolean("list", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &text.Text{
			Min:     uint(request.GetInt("min", 3)),
			Max:     uint(request.GetInt("max", 256)),
			Ascii:   request.GetBool("ascii", false),
			Sort:    request.GetBool("sort", false),
			Wtf:     request.GetInt("wtf", 0),
			Find:    request.GetStringSlice("find", []string{}),
			First:   request.GetBool("first", false),
			List:    request.GetBool("list", false),
			Decimal: request.GetBool("decimal", false),
			Paths:   paths,
		}), nil
	})
}

func (cmd *Mcp) addHash(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("hash",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithString("algo", mcp.Description("")),
		mcp.WithBoolean("all", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &hash.Hash{
			Algo:  request.GetStringSlice("algo", []string{types.SHA256}),
			All:   request.GetBool("all", false),
			Paths: paths,
		}), nil
	})
}

func (cmd *Mcp) addStat(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("stat",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("sort", mcp.Description("")),
		mcp.WithNumber("block", mcp.Description("")),
		mcp.WithNumber("min", mcp.Description("")),
		mcp.WithNumber("max", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &stat.Stat{
			Sort:  request.GetBool("sort", false),
			Block: fmt.Sprintf("%d", request.GetInt("block", 0)),
			Min:   request.GetFloat("min", 0.0),
			Max:   request.GetFloat("max", 1.0),
			Paths: paths,
		}), nil
	})
}

func (cmd *Mcp) addTest(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("test",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("domain", mcp.Description("")),
		mcp.WithBoolean("url", mcp.Description("")),
		mcp.WithBoolean("ip", mcp.Description("")),
		mcp.WithString("key", mcp.Description(""), mcp.Required()),
		mcp.WithArray("paths", mcp.Description(""), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		key, err := request.RequireString("key")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &test.Test{
			Domain: request.GetBool("domain", false),
			Url:    request.GetBool("url", false),
			Ip:     request.GetBool("Ip", false),
			Key:    key,
			Paths:  paths,
		}), nil
	})
}

func (cmd *Mcp) addDump(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("dump",
		mcp.WithDescription(""),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("json", mcp.Description("")),
		mcp.WithBoolean("jsonl", mcp.Description("")),
		mcp.WithBoolean("bootkey", mcp.Description("")),
		mcp.WithBoolean("lm", mcp.Description("")),
		mcp.WithBoolean("nt", mcp.Description("")),
		mcp.WithString("system", mcp.Description(""), mcp.Required()),
		mcp.WithString("ntds", mcp.Description("")),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		system, err := request.RequireString("system")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &dump.Dump{
			Json:    request.GetBool("json", false),
			Jsonl:   request.GetBool("jsonl", false),
			Bootkey: request.GetBool("bootkey", false),
			OnlyLm:  request.GetBool("lm", false),
			OnlyNt:  request.GetBool("nt", false),
			Paths:   []string{system, request.GetString("ntds", "")},
		}), nil
	})
}

func (cmd *Mcp) addHunt(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("hunt",
		mcp.WithDescription("Search for suspicious event logs in files"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("all", mcp.Description("events of all severity")),
		mcp.WithBoolean("uniq", mcp.Description("events without duplicates")),
		mcp.WithBoolean("sort", mcp.Description("events sorted by timestamp")),
		mcp.WithBoolean("json", mcp.Description("")),
		mcp.WithBoolean("jsonl", mcp.Description("")),
		mcp.WithNumber("block", mcp.Description("")),
		mcp.WithString("rule", mcp.Description("")),
		mcp.WithNumber("dist", mcp.Description("")),
		mcp.WithArray("paths", mcp.Description("Paths to search"), mcp.Required()),
	), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(cli, &hunt.Hunt{
			All:   request.GetBool("all", false),
			Uniq:  request.GetBool("uniq", false),
			Sort:  request.GetBool("sort", false),
			Json:  request.GetBool("json", false),
			Jsonl: request.GetBool("jsonl", false),
			Block: uint(request.GetInt("block", 65536)),
			Rule:  request.GetString("rule", ""),
			Dist:  request.GetFloat("dist", 0),
			Paths: paths,
		}), nil
	})
}

func global() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithBoolean("head", mcp.Description("")),
		mcp.WithBoolean("tail", mcp.Description("")),
		mcp.WithNumber("bytes", mcp.Description("")),
		mcp.WithNumber("lines", mcp.Description("")),
		mcp.WithString("regex", mcp.Description("")),
		mcp.WithString("password", mcp.Description("")),
		mcp.WithBoolean("raw", mcp.Description("")),
		mcp.WithBoolean("dry", mcp.Description("")),
		mcp.WithBoolean("deflate", mcp.Description("")),
		mcp.WithBoolean("extract", mcp.Description("")),
		mcp.WithBoolean("convert", mcp.Description("")),
	}
}

func execute(cli *cli.Globals, mode Mode) *mcp.CallToolResult {
	cli.Pipe = bytes.NewBuffer(nil)

	if err := mode.Run(cli); err != nil {
		return mcp.NewToolResultError(err.Error())
	}

	return mcp.NewToolResultText(cli.Pipe.String())
}
