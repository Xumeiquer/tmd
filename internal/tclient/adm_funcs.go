package tclient

import "github.com/zelenin/go-tdlib/client"

func (tc *TClient) RemoveDocument(fileId int32) error {
	request := &client.DeleteFileRequest{
		FileId: fileId,
	}

	_, err := tc.t.DeleteFile(request)
	return err
}

func (tc *TClient) SearchFileDownloads() (*client.FoundFileDownloads, error) {
	request := &client.SearchFileDownloadsRequest{
		Query:         "",
		OnlyActive:    false,
		OnlyCompleted: false,
		Offset:        "",
		Limit:         1000,
	}

	return tc.t.SearchFileDownloads(request)
}

func (tc *TClient) DeleteFile(fileId int32) (*client.Ok, error) {
	request := &client.DeleteFileRequest{
		FileId: fileId,
	}

	return tc.t.DeleteFile(request)
}
