package tclient

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

func (tc *TClient) addDownload(file *client.File) {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	slog.Debug(fmt.Sprintf("adding file %d to the file manager", file.Id))
	tc.downloadTracker[file.Id] = file
}

func (tc *TClient) removeDownload(fileId int32) {
	tc.mux.Lock()
	defer tc.mux.Unlock()

	slog.Debug(fmt.Sprintf("deleting file %d from the file manager", fileId))
	delete(tc.downloadTracker, fileId)
}

func (tc *TClient) uploadDownloadTracker(fileId int32, updateInfo *client.UpdateFile) {
	checkExists := func() bool {
		tc.mux.RLock()
		defer tc.mux.RUnlock()

		_, in := tc.downloadTracker[fileId]

		return in
	}

	updateTracker := func() {
		tc.mux.Lock()
		defer tc.mux.Unlock()

		tc.downloadTracker[fileId] = updateInfo.File
	}

	if !checkExists() {
		return
	}

	if updateInfo.File.Local.IsDownloadingCompleted {
		slog.Debug(fmt.Sprintf("download complete for file %d. Requesting deleteion from file manager", fileId))
		tc.removeDownload(fileId)
	} else {
		slog.Debug(fmt.Sprintf("updating file %d with the new information", fileId))
		updateTracker()
	}
}

func (tc *TClient) isDownloadCompleted() bool {
	tc.mux.RLock()
	defer tc.mux.RUnlock()

	slog.Debug(fmt.Sprintf("are all files fully download? %t", len(tc.downloadTracker) == 0))
	return len(tc.downloadTracker) == 0
}

func (tc *TClient) MoveDownloaded(storePath string) error {
	moveFunc := func(source, destination string) (err error) {
		src, err := os.Open(source)
		if err != nil {
			return err
		}
		defer src.Close()
		fi, err := src.Stat()
		if err != nil {
			return err
		}
		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		perm := fi.Mode() & os.ModePerm
		dst, err := os.OpenFile(destination, flag, perm)
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			dst.Close()
			os.Remove(destination)
			return err
		}
		err = dst.Close()
		if err != nil {
			return err
		}
		err = src.Close()
		if err != nil {
			return err
		}
		err = os.Remove(source)
		if err != nil {
			return err
		}
		return nil
	}

	slog.Debug(fmt.Sprintf("reading folder %s/documents", tc.cfg.Cache.File))
	var err error
	files, err := os.ReadDir(filepath.Join(tc.cfg.Cache.File, "documents"))
	if err != nil {
		return err
	}

	os.MkdirAll(storePath, os.ModePerm)

	for _, file := range files {
		if file.Type().IsRegular() {
			err = moveFunc(filepath.Join(tc.cfg.Cache.File, "documents", file.Name()), filepath.Join(storePath, file.Name()))
			slog.Debug(fmt.Sprintf("moving file %s to %s", file.Name(), filepath.Join(storePath, file.Name())))
			if err != nil {
				slog.Error("unable to move file.", "msg", err.Error())
				continue
			}
		}
	}

	return nil
}
