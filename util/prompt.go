package util

import (
	"context"
	"examples/SimpleBankProject/config"
	"examples/SimpleBankProject/prompt"
	"log"
)

func AddPromptStoreToRedis(c context.Context, location string, data prompt.PromptData) error {
	key := PromptStorePrefix + location

	/*jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling prompt data: %v", err)
		return err
	}*/

	err := config.Redis.Client.Set(c, key, data, PromptStoreExp).Err()
	if err != nil {
		log.Printf("Error setting prompt data in Redis: %v", err)
		return err
	}
	return nil
}
