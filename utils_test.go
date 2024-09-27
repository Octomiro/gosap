package gosap_test

import (
	"encoding/json"

	"github.com/octomiro/gosap"
)

func ItemsToJSON(items gosap.Items) string {
	json, err := json.MarshalIndent(items.Value, "", "  ")
	if err != nil {
		return ""
	}

	return string(json)
}

func JSONToItems(items string) []gosap.Item {
	var Items []gosap.Item

	err := json.Unmarshal([]byte(items), &Items)
	if err != nil {
		return nil
	}

	return Items
}
