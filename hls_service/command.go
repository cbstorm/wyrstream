package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/cbstorm/wyrstream/lib/logger"
)

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
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:a", "160k",
			"-b:v", "2M",
			"-maxrate:v", "2M",
			"-bufsize", "1M",
			"-crf", "18",
			"-preset", "ultrafast",
			"-f", "hls",
			"-hls_time", "6",
			"-hls_list_size", "6",
			"-hls_segment_filename", "public/" + stream_id + "/seg-%05d.ts",
		},
	}
}

func (c *ProcessHLSCommand) SetInput(i string) *ProcessHLSCommand {
	c.input = []string{"-i", i}
	return c
}

func (c *ProcessHLSCommand) SetOutput(m3u8_file string) *ProcessHLSCommand {
	c.output = fmt.Sprintf("public/%s/%s", c.stream_id, m3u8_file)
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
}

func (c *ProcessHLSCommand) Print() {
	args_str := strings.Join(c.args, " ")
	fmt.Printf("%s %s", c.name, args_str)
}
