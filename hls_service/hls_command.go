package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/cbstorm/wyrstream/lib/logger"
)

var process_hls_cmd_store *ProcessHLSCommandStore
var process_hls_cmd_store_sync sync.Once

func GetProcessHLSCommandStore() *ProcessHLSCommandStore {
	if process_hls_cmd_store == nil {
		process_hls_cmd_store_sync.Do(func() {
			process_hls_cmd_store = &ProcessHLSCommandStore{
				commands: make(map[string]*ProcessHLSCommand),
			}
		})
	}
	return process_hls_cmd_store
}

type ProcessHLSCommandStore struct {
	commands map[string]*ProcessHLSCommand
	mu       sync.Mutex
}

func (s *ProcessHLSCommandStore) Add(c *ProcessHLSCommand) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commands[c.stream_id] = c
}

func (s *ProcessHLSCommandStore) Get(stream_id string) *ProcessHLSCommand {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commands[stream_id]
}

func (s *ProcessHLSCommandStore) Remove(stream_id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.commands, stream_id)
}

type ProcessHLSCommand struct {
	ctx       context.Context
	cancelFn  context.CancelFunc
	name      string
	args      []string
	logger    *logger.Logger
	stream_id string
	input     []string
	output    string
}

func NewProcessHLSCommand(stream_id string) *ProcessHLSCommand {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProcessHLSCommand{
		ctx:       ctx,
		cancelFn:  cancel,
		stream_id: stream_id,
		logger:    logger.NewLogger(fmt.Sprintf("PROCESS_HLS_CMD - %s", stream_id)),
		name:      "ffmpeg",
		args: []string{
			"-v", "error",
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:a", "160k",
			"-b:v", "2M",
			"-maxrate:v", "2M",
			"-bufsize", "1M",
			"-crf", "18",
			"-preset", "ultrafast",
			"-f", "hls",
			"-hls_init_time", "2",
			"-hls_time", "5",
			"-hls_list_size", "5",
			"-hls_segment_filename", BuildHLSSegmentFilePath(stream_id),
		},
		output: BuildHLSm3u8FilePath(stream_id),
	}
}

func (c *ProcessHLSCommand) SetInput(i string) *ProcessHLSCommand {
	c.input = []string{"-i", i}
	return c
}

func (c *ProcessHLSCommand) SetStartNumber(s uint) *ProcessHLSCommand {
	c.args = append(c.args, "-start_number", fmt.Sprintf("%d", s))
	return c
}

func (c *ProcessHLSCommand) buildArgs() []string {
	out := []string{}
	out = append(out, c.input...)
	out = append(out, c.args...)
	out = append(out, c.output)
	return out
}
func (c *ProcessHLSCommand) Run() error {
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
	c.logger.Info("Starting process HLS")
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
	return nil
}
func (c *ProcessHLSCommand) Cancel() {
	c.cancelFn()
	c.logger.Info("Stopped.")
}

func (c *ProcessHLSCommand) Print() {
	args := c.buildArgs()
	fmt.Printf("%s %s", c.name, args)
}
