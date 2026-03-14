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

	ver "github.com/cuhsat/fox/v4/internal"
	cli "github.com/cuhsat/fox/v4/internal/cmd"
	std "github.com/cuhsat/fox/v4/internal/pkg/text"

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
Starts as MCP server.

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
	cmd.mcp = server.NewMCPServer("fox", ver.Version,
		server.WithToolCapabilities(true),
	)

	return nil
}

func (cmd *Mcp) Run(cli *cli.Globals) error {
	// prepare cli
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
		log.Printf("mcp: started on port %d\n", cmd.Port)
	}

	return server.NewStreamableHTTPServer(cmd.mcp).Start(fmt.Sprintf(":%d", cmd.Port))
}

func (cmd *Mcp) addCat(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("cat", addGlobals(
		mcp.WithDescription("Get file contents"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("uniq", mcp.Description("Filter using unique hash")),
		mcp.WithNumber("dist", mcp.Description("Filter lines using Levenshtein distance")),
		mcp.WithString("regex", mcp.Description("Filter lines using regular expression")),
		mcp.WithNumber("context", mcp.Description("Lines surrounding a match")),
		mcp.WithNumber("before", mcp.Description("Lines leading before a match")),
		mcp.WithNumber("after", mcp.Description("Lines trailing after a match")),
		mcp.WithArray("paths", mcp.Description("Process paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		cli.Regex = request.GetString("regex", "")

		return execute(process(cli, &request), &cat.Cat{
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
	cmd.mcp.AddTool(mcp.NewTool("hex", addGlobals(
		mcp.WithDescription("Get file content in hex"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("canonical", mcp.Description("Format output as canonical")),
		mcp.WithBoolean("hexdump", mcp.Description("Format output like hexdump")),
		mcp.WithBoolean("xxd", mcp.Description("Format output like xxd")),
		mcp.WithBoolean("decimal", mcp.Description("Format addresses as decimal")),
		mcp.WithBoolean("no-format", mcp.Description("Don't format output at all")),
		mcp.WithArray("paths", mcp.Description("Process paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &hex.Hex{
			Canonical: request.GetBool("canonical", false),
			Hexdump:   request.GetBool("hexdump", false),
			Xxd:       request.GetBool("xxd", false),
			Decimal:   request.GetBool("decimal", false),
			NoFormat:  request.GetBool("no-format", false),
			Paths:     paths,
		}), nil
	})
}

func (cmd *Mcp) addText(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("text", addGlobals(
		mcp.WithDescription("Get file text contents"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithNumber("min", mcp.Description("Minimum string length")),
		mcp.WithNumber("max", mcp.Description("Maximum string length")),
		mcp.WithBoolean("ascii", mcp.Description("Get only ASCII strings")),
		mcp.WithBoolean("sort", mcp.Description("Sort strings alphabetically")),
		mcp.WithNumber("wtf", mcp.Description("Get string classifications by level")),
		mcp.WithArray("find", mcp.Description("Get only strings the match classes"), mcp.Required()),
		mcp.WithBoolean("first", mcp.Description("Get only the first string class")),
		mcp.WithBoolean("list", mcp.Description("Get only the classification list")),
		mcp.WithBoolean("decimal", mcp.Description("Format addresses as decimal")),
		mcp.WithArray("paths", mcp.Description("Process paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &text.Text{
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
	cmd.mcp.AddTool(mcp.NewTool("hash", addGlobals(
		mcp.WithDescription("Get file hashes and checksums"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithString("algo", mcp.Description("Use algorithms")),
		mcp.WithBoolean("all", mcp.Description("Use all algorithms")),
		mcp.WithArray("paths", mcp.Description("Hash files in paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &hash.Hash{
			Algo:  request.GetStringSlice("algo", []string{types.SHA256}),
			All:   request.GetBool("all", false),
			Paths: paths,
		}), nil
	})
}

func (cmd *Mcp) addStat(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("stat", addGlobals(
		mcp.WithDescription("Get file stats and entropy"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("sort", mcp.Description("Sort files by path")),
		mcp.WithString("block", mcp.Description("Block size for analysis")),
		mcp.WithNumber("min", mcp.Description("Minimum entropy value")),
		mcp.WithNumber("max", mcp.Description("Maximum entropy value")),
		mcp.WithArray("paths", mcp.Description("Analyse files in paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &stat.Stat{
			Sort:  request.GetBool("sort", false),
			Block: request.GetString("block", ""),
			Min:   request.GetFloat("min", 0.0),
			Max:   request.GetFloat("max", 1.0),
			Paths: paths,
		}), nil
	})
}

func (cmd *Mcp) addTest(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("test", addGlobals(
		mcp.WithDescription("Test suspicious files using VirusTotal"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(true),
		mcp.WithBoolean("domain", mcp.Description("Files contain domains")),
		mcp.WithBoolean("url", mcp.Description("Files contain URLs")),
		mcp.WithBoolean("ip", mcp.Description("Files contain IPs")),
		mcp.WithString("key", mcp.Description("Use required VirusTotal API key"), mcp.Required()),
		mcp.WithArray("paths", mcp.Description("Test files in paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		key, err := request.RequireString("key")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &test.Test{
			Domain: request.GetBool("domain", false),
			Url:    request.GetBool("url", false),
			Ip:     request.GetBool("Ip", false),
			Key:    key,
			Paths:  paths,
		}), nil
	})
}

func (cmd *Mcp) addDump(cli *cli.Globals) {
	cmd.mcp.AddTool(mcp.NewTool("dump", addGlobals(
		mcp.WithDescription("Dump sensitive data"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("json", mcp.Description("Dump hashes as JSON objects")),
		mcp.WithBoolean("jsonl", mcp.Description("Dump hashes as JSON lines")),
		mcp.WithBoolean("bootkey", mcp.Description("Extract the bootkey")),
		mcp.WithBoolean("lm", mcp.Description("Extract only LM hashes")),
		mcp.WithBoolean("nt", mcp.Description("Extract only NT hashes")),
		mcp.WithString("system", mcp.Description("System registry hive"), mcp.Required()),
		mcp.WithString("ntds", mcp.Description("Active Directory database")),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		system, err := request.RequireString("system")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &dump.Dump{
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
	cmd.mcp.AddTool(mcp.NewTool("hunt", addGlobals(
		mcp.WithDescription("Search for suspicious events in files"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithBoolean("all", mcp.Description("Get all events regardless of severity")),
		mcp.WithBoolean("uniq", mcp.Description("Get only unique events")),
		mcp.WithBoolean("sort", mcp.Description("Sort events by timestamp")),
		mcp.WithBoolean("json", mcp.Description("Convert events to JSON objects")),
		mcp.WithBoolean("jsonl", mcp.Description("Convert events to JSON lines")),
		mcp.WithString("block", mcp.Description("Block size used for carving")),
		mcp.WithString("rule", mcp.Description("Filter events using Sigma rule file")),
		mcp.WithNumber("dist", mcp.Description("Filter events using Levenshtein distance")),
		mcp.WithArray("paths", mcp.Description("Search in paths"), mcp.Required()),
	)...), func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		paths, err := request.RequireStringSlice("paths")

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return execute(process(cli, &request), &hunt.Hunt{
			All:   request.GetBool("all", false),
			Uniq:  request.GetBool("uniq", false),
			Sort:  request.GetBool("sort", false),
			Json:  request.GetBool("json", false),
			Jsonl: request.GetBool("jsonl", false),
			Block: request.GetString("block", "65536"),
			Rule:  request.GetString("rule", ""),
			Dist:  request.GetFloat("dist", 0),
			Paths: paths,
		}), nil
	})
}

func addGlobals(opts ...mcp.ToolOption) []mcp.ToolOption {
	return append([]mcp.ToolOption{
		mcp.WithBoolean("head", mcp.Description("Limit file head")),
		mcp.WithBoolean("tail", mcp.Description("Limit file tail")),
		mcp.WithString("bytes", mcp.Description("Limit by byte count")),
		mcp.WithString("lines", mcp.Description("Limit by line count")),
		mcp.WithString("regex", mcp.Description("Filter using regular expression")),
		mcp.WithString("password", mcp.Description("Use password for archives")),
		mcp.WithBoolean("dry", mcp.Description("Get only affected files")),
		mcp.WithBoolean("raw", mcp.Description("Don't process file content at all")),
		mcp.WithBoolean("no-deflate", mcp.Description("Don't deflate file content")),
		mcp.WithBoolean("no-extract", mcp.Description("Don't extract file content")),
		mcp.WithBoolean("no-convert", mcp.Description("Don't convert file content")),
	}, opts...)
}

func process(cli *cli.Globals, req *mcp.CallToolRequest) *cli.Globals {
	cli.Head = req.GetBool("head", false)
	cli.Tail = req.GetBool("tail", false)
	cli.Bytes = req.GetString("bytes", "")
	cli.Lines = req.GetString("lines", "")
	cli.Regex = req.GetString("regex", "")
	cli.Password = req.GetString("password", "")
	cli.DryRun = req.GetBool("dry", false)
	cli.Raw = req.GetBool("raw", false)
	cli.NoDeflate = req.GetBool("no-deflate", false)
	cli.NoExtract = req.GetBool("no-extract", false)
	cli.NoConvert = req.GetBool("no-convert", false)

	return cli
}

func execute(cli *cli.Globals, mode Mode) *mcp.CallToolResult {
	pipe := bytes.NewBuffer(nil)

	std.Redirect(pipe)

	if err := mode.Run(cli); err != nil {
		return mcp.NewToolResultError(err.Error())
	}

	return mcp.NewToolResultText(pipe.String())
}
