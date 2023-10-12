package tclient

import (
	"fmt"
	"log/slog"

	progress_bar "github.com/Xumeiquer/tmd/internal/progress-bar"
	"github.com/zelenin/go-tdlib/client"
)

func (tc *TClient) Poll(sync chan struct{ Id int32 }, listener *client.Listener, prb *progress_bar.ProgressBar) {
	for info := range listener.Updates {
		var fileId int32

		if info.GetType() == client.TypeUpdateFile {
			updateInfo, ok := info.(*client.UpdateFile)
			if ok {
				fileId = updateInfo.File.Id
				slog.Debug(fmt.Sprintf("file download update for %d: %d/%d", updateInfo.File.Id, updateInfo.File.Local.DownloadedSize, updateInfo.File.Size))

				tc.uploadDownloadTracker(updateInfo.File.Id, updateInfo)

				if prb != nil && prb.IsTracked(updateInfo.File.Id) {
					prb.UpdateTracker(updateInfo.File.Id, updateInfo)
				}
			} else {
				// This should never happen
				slog.Error("unable to cast to *client.UpdateFile", "context", "downloadDocumentEx")
			}
		}

		if tc.isDownloadCompleted() {
			slog.Debug(fmt.Sprintf("file %d download complete. Synching...", fileId))

			sync <- struct{ Id int32 }{fileId}
			if prb != nil {
				prb.MarkAsDone(fileId)
			}
		}
	}
}

func (tc *TClient) Wait(sync chan struct{ Id int32 }) {
	for {
		select {
		case id := <-sync:
			slog.Debug(fmt.Sprintf("%d download complete", id.Id))
			if len(tc.downloadTracker) == 0 {
				slog.Debug("all downloads complete. Shutting down...")
				return
			}
		}
	}
}
