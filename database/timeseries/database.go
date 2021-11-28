package timeseries

import "github.com/nakabonne/tstorage"

var Database tstorage.Storage

func Init(directory string) error {
	db, err := tstorage.NewStorage(
		tstorage.WithDataPath(directory),
	)
	if err != nil {
		return err
	}
	Database = db
	return nil
}

func Close() error {
	return Database.Close()
}
