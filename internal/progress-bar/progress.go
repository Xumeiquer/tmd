package progress_bar

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/zelenin/go-tdlib/client"
)

type ProgressBar struct {
	cfg      *Config
	pw       progress.Writer
	trackers map[int32]*progress.Tracker
	mux      sync.RWMutex
	log      *slog.Logger
}

func New() *ProgressBar {
	cfg, err := Initconfig()
	if err != nil {
		slog.Error("found an error while initializing progress bar", "msg", err.Error())
	}

	prb := &ProgressBar{
		cfg: cfg,
	}

	log := logger.GetLog(cfg.Log.Level, cfg.Log.Type, cfg.Log.To)
	prb.log = log.With("context", "listCmdEx")

	prb.log.Debug("Progress tracker initiated", "pkg", "progress-bar", "context", "init")

	prb.trackers = make(map[int32]*progress.Tracker)

	prb.pw = progress.NewWriter()
	prb.pw.SetOutputWriter(os.Stdout)
	prb.pw.SetAutoStop(true)
	prb.pw.SetTrackerLength(24)
	prb.pw.SetMessageWidth(26)
	// prb.pw.SetNumTrackersExpected(*flagNumTrackers)
	prb.pw.SetSortBy(progress.SortByPercentDsc)
	prb.pw.SetStyle(progress.StyleDefault)
	prb.pw.SetTrackerPosition(progress.PositionRight)
	prb.pw.SetUpdateFrequency(time.Millisecond * 100)
	prb.pw.Style().Colors = progress.StyleColorsExample
	prb.pw.Style().Options.PercentFormat = "%4.1f%%"
	prb.pw.Style().Visibility.ETA = false
	prb.pw.Style().Visibility.ETAOverall = false
	prb.pw.Style().Visibility.Percentage = true
	prb.pw.Style().Visibility.Speed = true
	prb.pw.Style().Visibility.SpeedOverall = false
	prb.pw.Style().Visibility.Time = true
	prb.pw.Style().Visibility.TrackerOverall = true
	prb.pw.Style().Visibility.Value = true
	prb.pw.Style().Visibility.Pinned = false

	// call Render() in async mode; yes we don't have any trackers at the moment
	go prb.pw.Render()

	return prb
}

func (prb *ProgressBar) AddTracker(fileId int32, totalSize int64) {
	prb.mux.Lock()
	defer prb.mux.Unlock()

	prb.log.Debug(fmt.Sprintf("added file id %d to the tracker", fileId), "pkg", "progress-bar", "context", "AddTracker")

	t := &progress.Tracker{
		Message:    fmt.Sprintf("Downloading File    #%d", fileId),
		DeferStart: false,
		Total:      totalSize,
		Units:      progress.UnitsBytes,
	}

	prb.trackers[fileId] = t
	prb.pw.AppendTracker(t)

	t.Start()
}

func (prb *ProgressBar) IsTracked(fileId int32) bool {
	prb.mux.RLock()
	defer prb.mux.RUnlock()

	_, in := prb.trackers[fileId]
	return in
}

func (prb *ProgressBar) UpdateTracker(fileId int32, updateInfo *client.UpdateFile) {
	prb.mux.Lock()
	defer prb.mux.Unlock()

	prb.log.Debug(fmt.Sprintf("Tracker %d updated\n", fileId), "pkg", "progress-bar", "context", "UpdateTracker")

	if tracker, in := prb.trackers[fileId]; in {
		tracker.Increment(updateInfo.File.Local.DownloadedSize - tracker.Value())
		if updateInfo.File.Local.IsDownloadingCompleted {
			tracker.MarkAsDone()
		}
	}
}

func (prb *ProgressBar) MarkAsDone(fileId int32) {
	if tracker, in := prb.trackers[fileId]; in {
		tracker.MarkAsDone()
	}
}
