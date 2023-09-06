package s3

type UploadError struct {
	Err error
}

func (ue *UploadError) Error() string {
	return "an error occurred while uploading to storage: " + ue.Err.Error()
}

type DownloadError struct {
	Err error
}

func (de *DownloadError) Error() string {
	return "an error occurred while downloading from storage: " + de.Err.Error()
}
