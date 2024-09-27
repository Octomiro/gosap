package gosap_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/octomiro/gosap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	config gosap.Config
	update = flag.Bool("update", false, "update .golden files")
)

func TestMain(m *testing.M) {
	var err error
	config, err = gosap.LoadConfig(".")
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(m.Run())
}

func TestGetItems(t *testing.T) {
	t.Parallel()

	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	items, err := session.GetItems(config)
	assert.NoError(t, err)

	gp := filepath.Join("testdata", filepath.FromSlash(t.Name())+".golden")

	if *update {
		err := os.WriteFile(gp, []byte(ItemsToJSON(*items)), 0o600)
		require.NoError(t, err)
	}

	goldenContent, err := os.ReadFile(gp)
	require.NoError(t, err)

	assert.Equal(t, []byte(ItemsToJSON(*items)), goldenContent)

	t.Log(items)
}

func TestGetItem(t *testing.T) {
	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	itemsGp := filepath.Join("testdata", "TestGetItems.golden")

	itemsBytes, err := os.ReadFile(itemsGp)
	require.NoError(t, err)

	items := JSONToItems(string(itemsBytes))

	for _, item := range items {
		t.Run("Test GetItem", func(t *testing.T) {
			t.Parallel()

			tm, err := session.GetItem(config, item.ItemCode)
			assert.NoError(t, err)

			t.Log(tm)

			assert.Equal(t, item.ItemCode, tm.ItemCode)
			assert.Equal(t, item.ItemName, tm.ItemName)
			assert.Equal(t, item.PurchaseUnitWidth, tm.PurchaseUnitWidth)
		})
	}
}
