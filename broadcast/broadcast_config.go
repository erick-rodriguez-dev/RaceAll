package broadcast

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type BroadcastConfig struct {
	Root struct {
		UpdListenerPort    int    `json:"updListenerPort"`
		ConnectionPassword string `json:"connectionPassword"`
		CommandPassword    string `json:"commandPassword"`
	}

	lock sync.Mutex
}

var (
	globalLock   sync.Mutex
	globalConfig *BroadcastConfig
)

func init() {
	globalConfig = &BroadcastConfig{}
}

func GetConfiguration() *BroadcastConfigRoot {
	globalLock.Lock()
	defer globalLock.Unlock()

	configPath := GetAccConfigPath() + "broadcasting.json"

	if fileInfo, err := os.Stat(configPath); err == nil && !fileInfo.IsDir() {
		fileStream, err := os.Open(configPath)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		defer fileStream.Close()

		config := getConfigurationFromStream(fileStream)

		if config != nil {
			if config.UpdListenerPort == 0 {
				config.UpdListenerPort = 9000

				data, err := json.MarshalIndent(config, "", "  ")
				if err == nil {
					os.WriteFile(configPath, data, 0644)
					fmt.Printf("Auto-Changed the port number in \"%sbroadcasting.json\" from 0 to 9000.\n", GetAccConfigPath())
				}
			}

			return config
		}
	}

	return nil
}

func getConfigurationFromStream(stream io.Reader) *BroadcastConfigRoot {
	var jsonString string

	data, err := io.ReadAll(stream)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	jsonString = string(data)

	jsonString = removeNullChars(jsonString)

	var config BroadcastConfigRoot

	err = json.Unmarshal([]byte(jsonString), &config)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &config
}

func removeNullChars(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			result = append(result, s[i])
		}
	}
	return string(result)
}

type BroadcastConfigRoot struct {
	UpdListenerPort    int    `json:"updListenerPort"`
	ConnectionPassword string `json:"connectionPassword"`
	CommandPassword    string `json:"commandPassword"`
}
