package broadcast

import (
	"os"
	"path/filepath"
)

type FileUtil struct{}

func GetAccPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, "Documents", "Assetto Corsa Competizione") + string(filepath.Separator)
}

func GetAccConfigPath() string {
	return GetAccPath() + "Config" + string(filepath.Separator)
}

func GetCustomsPath() string {
	return GetAccPath() + "Customs" + string(filepath.Separator)
}

func GetCarsPath() string {
	return GetCustomsPath() + "Cars" + string(filepath.Separator)
}

func GetLiveriesPath() string {
	return GetCustomsPath() + "Liveries" + string(filepath.Separator)
}

func GetSetupsPath() string {
	return GetAccPath() + "Setups" + string(filepath.Separator)
}
