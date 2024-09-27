package gosap_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

func TestGetDeliveryNotes(t *testing.T) {
	t.Parallel()

	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	notes, err := session.GetDeliveryNotes(config)
	require.NoError(t, err)

	t.Log(notes)

	gp := filepath.Join("testdata", filepath.FromSlash(t.Name())+".golden")
	if *update {
		err := os.WriteFile(gp, []byte(ToJSON(notes.Value)), 0o600)
		require.NoError(t, err)
	}

	goldenContent, err := os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, []byte(ToJSON(notes.Value)), goldenContent)
}

func TestGetDeliveryNote(t *testing.T) {
	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	deliveryNotesFile := filepath.Join("testdata", "TestGetDeliveryNotes.golden")

	bytes, err := os.ReadFile(deliveryNotesFile)
	require.NoError(t, err)

	notes := make([]gosap.DeliveryNote, 0, 10)
	err = FromJSON(string(bytes), &notes)
	require.NoError(t, err)

	t.Log(notes)

	for _, note := range notes {
		t.Run("Test get specific delivery note", func(t *testing.T) {
			t.Parallel()

			tnote, err := session.GetDeliveryNote(config, strconv.Itoa(note.DocEntry))
			assert.NoError(t, err)

			t.Log(tnote)

			assert.Equal(t, note, *tnote)
		})
	}

	t.Run("Test getting a delivery note with wrong id", func(t *testing.T) {
		t.Parallel()

		_, err := session.GetDeliveryNote(config, "3242")
		assert.Error(t, err)
	})
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
