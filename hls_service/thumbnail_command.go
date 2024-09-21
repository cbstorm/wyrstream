package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/logger"
)

const THUMBNAIL_TIME_TICKER = 10 * time.Second

var process_thumbnail_cmd_store *ProcessThumbnailCommandStore
var process_thumbnail_cmd_store_sync sync.Once

func GetProcessThumbnailCommandStore() *ProcessThumbnailCommandStore {
	if process_thumbnail_cmd_store == nil {
		process_thumbnail_cmd_store_sync.Do(func() {
			process_thumbnail_cmd_store = &ProcessThumbnailCommandStore{
				commands: make(map[string]*ProcessThumbnailCommand),
			}
		})
	}
	return process_thumbnail_cmd_store
}

type ProcessThumbnailCommandStore struct {
	commands map[string]*ProcessThumbnailCommand
	mu       sync.Mutex
}

func (s *ProcessThumbnailCommandStore) Add(c *ProcessThumbnailCommand) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commands[c.stream_id] = c
}

func (s *ProcessThumbnailCommandStore) Get(stream_id string) *ProcessThumbnailCommand {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commands[stream_id]
}

func (s *ProcessThumbnailCommandStore) Remove(stream_id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.commands, stream_id)
}

type ProcessThumbnailCommand struct {
	ctx       context.Context
	cancelFn  context.CancelFunc
	name      string
	args      []string
	logger    *logger.Logger
	stream_id string
	input     []string
	output    string
}

func NewProcessThumbnailCommand(stream_id string) *ProcessThumbnailCommand {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProcessThumbnailCommand{
		ctx:       ctx,
		cancelFn:  cancel,
		stream_id: stream_id,
		logger:    logger.NewLogger(fmt.Sprintf("PROCESS_THUMBNAIL_CMD - %s", stream_id)),
		name:      "ffmpeg",
		args: []string{
			"-v", "error",
			"-q:v", "1",
			"-frames:v", "1",
		},
		output: BuildThumbnailFilePath(stream_id),
	}
}
func (c *ProcessThumbnailCommand) setInput(i string) *ProcessThumbnailCommand {
	c.input = []string{"-i", i}
	return c
}

func (c *ProcessThumbnailCommand) Start() {
	ticker := time.NewTicker(THUMBNAIL_TIME_TICKER)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := c.run(); err != nil {
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *ProcessThumbnailCommand) getInput() string {
	files := GetListSegmentFilesByStreamId(c.stream_id)
	if len(*files) == 0 {
		return ""
	}
	s := (*files)[len(*files)-1]
	s = fmt.Sprintf("%s/%s", BuildHLSStreamDir(c.stream_id), s)
	return s
}

func (c *ProcessThumbnailCommand) buildArgs() []string {
	out := []string{}
	out = append(out, c.input...)
	out = append(out, c.args...)
	out = append(out, "-y")
	out = append(out, c.output)
	return out
}

func (c *ProcessThumbnailCommand) run() error {
	c.logger.Info("RUN")
	input := c.getInput()
	if input == "" {
		c.logger.Info("EMPTY_SEGMENT")
		return nil
	}
	c.setInput(input)
	c.Print()
	args := c.buildArgs()
	cmd := exec.CommandContext(c.ctx, c.name, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.Cancel()
		c.logger.Error("Could not pipe stderr with err: %v", err)
		return err
	}
	if err := cmd.Start(); err != nil {
		c.Cancel()
		c.logger.Error("Could not start command with err: %v", err)
		return err
	}
	stderr_output, err := io.ReadAll(stderr)
	if err != nil {
		c.logger.Error("Could not read from stderr with err: %v", err)
		return err
	}
	c.logger.Error(string(stderr_output))
	if err := cmd.Wait(); err != nil {
		c.logger.Error("Could not wait for command complete with err: %v", err)
		return err
	}
	c.Print()
	c.logger.Info("DONE")
	return nil
}

func (c *ProcessThumbnailCommand) Cancel() {
	c.cancelFn()
	c.logger.Info("Stopped.")
}

func (c *ProcessThumbnailCommand) Print() {
	args := c.buildArgs()
	fmt.Printf("%s %s", c.name, args)
}
