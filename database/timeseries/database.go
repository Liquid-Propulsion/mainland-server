package timeseries

import "github.com/nakabonne/tstorage"

var Database tstorage.Storage

func Init(inMemory bool, directory string) error {
	if !inMemory {
		db, err := tstorage.NewStorage(
			tstorage.WithDataPath(directory),
		)
		if err != nil {
			return err
		}
		Database = db
	} else {
		db, err := tstorage.NewStorage()
		if err != nil {
			return err
		}
		Database = db
	}
	return nil
}

func Close() error {
	return Database.Close()
}
