/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package tg

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/zelenin/go-tdlib/client"
)

type DownloadManager interface {
	AddFile(*client.File, string)
	IsDownloadCompleted() bool
	Wait()
}

type DM struct {
	listener      *client.Listener
	mux           *sync.RWMutex
	downloadQueue []*item
	pollRunning   bool
}

type item struct {
	file      *client.File
	storePath string
}

func NewDM(listener *client.Listener) *DM {
	return &DM{
		listener:      listener,
		mux:           new(sync.RWMutex),
		pollRunning:   false,
		downloadQueue: make([]*item, 0),
	}
}

func (dm *DM) AddFile(file *client.File, storePath string) {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	elem := &item{
		file:      file,
		storePath: storePath,
	}
	dm.downloadQueue = append(dm.downloadQueue, elem)
}

func (dm *DM) IsDownloadCompleted() bool {
	dm.mux.RLock()
	defer dm.mux.RUnlock()

	for _, item := range dm.downloadQueue {
		if item.file.Local.IsDownloadingActive {
			return false
		}
	}
	return true
}

func (dm *DM) enqueued(update *client.UpdateFile) bool {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	for idx, e := range dm.downloadQueue {
		if e.file.Remote.Id == update.File.Remote.Id {
			dm.downloadQueue[idx].file = update.File
			return true
		}
	}
	return false
}

func (dm *DM) Wait() {
	if !dm.listener.IsActive() {
		slog.Info("listener is not active")
		return
	}

	for info := range dm.listener.Updates {
		switch info.GetType() {
		case client.TypeUpdateFile:
			updateInfo, ok := info.(*client.UpdateFile)
			if ok {
				if dm.enqueued(updateInfo) {
					dm.printProcess(updateInfo)
					if updateInfo.File.Local.IsDownloadingCompleted {
						fmt.Print("\n")
						dm.downloadCompleted(updateInfo.File.Remote.Id)
					}
				}
			}
		}
		if dm.IsDownloadCompleted() {
			break
		}
	}
}

func (dm *DM) downloadCompleted(id string) {
	dm.mux.Lock()
	defer dm.mux.Unlock()

	for idx, e := range dm.downloadQueue {
		if e.file.Remote.Id == id {
			go dm.moveFile(e.file.Local.Path, e.storePath)
			dm.downloadQueue = append(dm.downloadQueue[:idx], dm.downloadQueue[idx+1:]...)
			return
		}
	}
}

func (dm *DM) moveFile(src, dst string) error {
	destinatinDir, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	fileName := filepath.Base(src)
	sourcePath := src
	destPath := filepath.Join(destinatinDir, fileName)

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DM) printProcess(update *client.UpdateFile) {
	id := update.File.Id
	size := update.File.ExpectedSize
	dsize := update.File.Local.DownloadedSize
	percent := (float32(dsize) / float32(size)) * 100
	downlaoded := update.File.Local.IsDownloadingCompleted

	fmt.Printf("\033[1A\033[K")
	fmt.Printf("Downloading %d...\t%d/%d (%.2f%%)\t%t\n", id, dsize, size, percent, downlaoded)
}
