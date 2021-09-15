package util

import (
	"time"

	"github.com/cheggaaa/pb/v3"
)

const refreshRate = 500 * time.Millisecond

func newProgressBar(progressText string) *pb.ProgressBar {
	var limit int64 = 1024 * 1024 * 500

	tmpl := `{{(cycle . "[ ]" "[/]" "[x]")}} {{(string . "progressText")}}`

	bar := pb.ProgressBarTemplate(tmpl).Start64(limit)
	bar.Set("progressText", progressText)
	bar.SetRefreshRate(refreshRate)

	return bar
}

func showProgress(bar *pb.ProgressBar) {
	// moves to next cycle state
	time.Sleep(refreshRate)
	bar.Increment()
}

func endProgress(bar *pb.ProgressBar) {
	time.Sleep(refreshRate)
	bar.Finish()
}

type ProgressFunc func() error

func StartProgress(progressName string, progress ProgressFunc) error {
	pb := newProgressBar(progressName)

	err := progress()
	if err != nil {
		return err
	}

	showProgress(pb)
	endProgress(pb)

	return nil
}
