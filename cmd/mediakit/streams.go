package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

func (this *app) Streams(ctx context.Context) error {
	count := uint(0)
	files := []*media{}

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := Walk(ctx, paths, &count, this.offset, this.limit, func(path string, info os.FileInfo) error {
		if media, err := this.ProcessStreams(path); err != nil {
			if *this.quiet == false {
				this.Logger.Print(filepath.Base(path), ": ", err)
			}
		} else {
			files = append(files, media)
		}
		return nil
	}); err != nil {
		return err
	}

	// Print out stream information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Stream", "Type"})
	table.SetAutoFormatHeaders(false)
	for _, file := range files {
		table.AppendBulk(FormatStreams(file))
	}
	table.Render()

	return nil
}

func (this *app) ProcessStreams(path string) (*media, error) {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer this.MediaManager.Close(media)

	// Create obj
	m := NewMedia(media)

	// Append streams
	for _, stream := range media.Streams() {
		m.Streams = append(m.Streams, NewStream(stream))
	}

	// Return success
	return m, nil
}
