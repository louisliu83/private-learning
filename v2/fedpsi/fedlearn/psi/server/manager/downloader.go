package manager

import (
	"context"

	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
)

type downloadTask struct {
	ctx context.Context
	ds  *model.Dataset
}

type HttpDownloader struct {
	downloadTaskPool chan *downloadTask
}

var (
	downloader = &HttpDownloader{
		downloadTaskPool: make(chan *downloadTask, 64),
	}
)

func GetHttpDownloader() *HttpDownloader {
	return downloader
}

func (s *HttpDownloader) DownloadDataset(ctx context.Context, ds *model.Dataset) error {
	s.downloadTaskPool <- &downloadTask{
		ctx: ctx,
		ds:  ds,
	}
	return nil
}

func (s *HttpDownloader) Run(ctx context.Context) {
	log.Infoln(ctx, "HTTP Downloader starts to run")
	paralism := config.GetConfig().DataSet.Downloaders
	if paralism == 0 {
		paralism = 1
	}

	for i := 1; i <= paralism; i++ {
		log.Infof(ctx, "HTTP Downloader %d/%d is running...", i, paralism)
		go func() {
			for {
				t := <-s.downloadTaskPool
				if err := download(t.ctx, t.ds); err != nil {
					log.Errorf(ctx, "Download failed.")
				}
			}
		}()
	}
}
