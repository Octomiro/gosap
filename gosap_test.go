package gosap_test

import (
	"encoding/json"
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
	assert.True(t, len(notes.Value) > 0)

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

func TestGetPurchaseOrders(t *testing.T) {
	t.Parallel()

	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	orders, err := session.GetPurchaseOrders(config)
	require.NoError(t, err)

	t.Log(orders)

	gp := filepath.Join("testdata", filepath.FromSlash(t.Name())+".golden")
	if *update {
		err := os.WriteFile(gp, []byte(ToJSON(orders.Value)), 0o600)
		require.NoError(t, err)
	}

	goldenContent, err := os.ReadFile(gp)
	require.NoError(t, err)
	assert.Equal(t, []byte(ToJSON(orders.Value)), goldenContent)
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

	for _, note := range notes[:30] {
		t.Run("Test get specific delivery note", func(t *testing.T) {
			t.Parallel()

			tnote, err := session.GetDeliveryNote(config, strconv.Itoa(note.DocEntry))
			require.NoError(t, err)
			require.NotNil(t, tnote)

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

func TestGetPurchaseOrder(t *testing.T) {
	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	deliveryNotesFile := filepath.Join("testdata", "TestGetPurchaseOrders.golden")

	bytes, err := os.ReadFile(deliveryNotesFile)
	require.NoError(t, err)

	notes := make([]gosap.PurchaseOrder, 0, 10)
	err = FromJSON(string(bytes), &notes)
	require.NoError(t, err)

	t.Log(notes)

	for _, note := range notes[:30] {
		t.Run("purchase_order_with_valid_id_is_retrieved", func(t *testing.T) {
			t.Parallel()

			tnote, err := session.GetPurchaseOrder(config, strconv.Itoa(note.DocEntry))
			require.NoError(t, err)
			require.NotNil(t, tnote)

			t.Log(tnote)

			assert.Equal(t, note, *tnote)
		})
	}

	t.Run("purchase_order_with_invalid_id_returns_error", func(t *testing.T) {
		t.Parallel()

		_, err := session.GetDeliveryNote(config, "3242")
		assert.Error(t, err)
	})
}

func TestGetPurchaseDeliveryNotes(t *testing.T) {
	t.Parallel()

	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	notes, err := session.GetPurchaseDeliveryNotes(config)
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

func TestCreatePurchaseDeliveryNote(t *testing.T) {
	t.Parallel()

	session, err := gosap.Authenticate(config)
	require.NoError(t, err)

	note := gosap.PurchaseDeliveryNote{
		CardCode: "V10000",
		DocumentLines: []gosap.PurchaseDeliveryNoteLine{
			{
				ItemCode: "I00007",
				Quantity: 20,
			},
		},
	}

	payload, err := json.Marshal(note)
	require.NoError(t, err)

	t.Log(string(payload))

	ok, err := session.CreatePurchaseDeliveryNote(config, note)

	assert.True(t, ok)
	assert.NoError(t, err)
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

	for _, item := range items[:20] {
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
