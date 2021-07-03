package config

type StorageConfig struct {
	Folder string
}

func LoadStorageConfig() StorageConfig {
	return StorageConfig{
		Folder: RequireEnvString("STORAGE_FOLDER"),
	}
}
