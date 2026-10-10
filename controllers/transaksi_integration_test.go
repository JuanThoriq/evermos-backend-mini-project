package controllers

import (
	"bytes"
	"evermos-mini/middleware"
	"evermos-mini/models"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func transactionFixture(t *testing.T) (*gorm.DB, models.User, models.Alamat, []models.Produk) {
	t.Helper()
	db := integrationDB(t)
	user := models.User{NamaUser: "Buyer", Email: "buyer@test.local", NoTelp: "08111", KataSandi: "hash", Role: "user", TanggalLahir: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	toko := models.Toko{NamaToko: "Seller", UserID: user.ID}
	if err := db.Create(&toko).Error; err != nil {
		t.Fatal(err)
	}
	category := models.Category{NamaCategory: "Test"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}
	alamat := models.Alamat{JudulAlamat: "Rumah", NamaPenerima: "Buyer", NoTelp: "08111", DetailAlamat: "Test", UserID: user.ID}
	if err := db.Create(&alamat).Error; err != nil {
		t.Fatal(err)
	}
	products := []models.Produk{
		{NamaProduk: "A", Slug: "a", HargaKonsumen: 10000, HargaReseller: 8000, Stok: 5, TokoID: toko.ID, CategoryID: category.ID},
		{NamaProduk: "B", Slug: "b", HargaKonsumen: 20000, HargaReseller: 15000, Stok: 1, TokoID: toko.ID, CategoryID: category.ID},
	}
	if err := db.Create(&products).Error; err != nil {
		t.Fatal(err)
	}
	return db, user, alamat, products
}

func transactionApp() *fiber.App {
	app := fiber.New()
	app.Post("/trx", middleware.JWTProtected, CreateTrx)
	return app
}

func TestCreateTransactionCommitsStockAndSnapshot(t *testing.T) {
	db, user, alamat, products := transactionFixture(t)
	body := []byte(`{"method_bayar":"bca","alamat_kirim":` + strconv.FormatUint(uint64(alamat.ID), 10) + `,"detail_trx":[{"product_id":` + strconv.FormatUint(uint64(products[0].ID), 10) + `,"kuantitas":2}]}`)
	req := httptest.NewRequest("POST", "/trx", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, user.ID))
	resp, err := transactionApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var product models.Produk
	if err := db.First(&product, products[0].ID).Error; err != nil {
		t.Fatal(err)
	}
	if product.Stok != 3 {
		t.Fatalf("stock=%d, want 3", product.Stok)
	}
	var transactions, details, logs int64
	db.Model(&models.Transaksi{}).Count(&transactions)
	db.Model(&models.DetailTransaksi{}).Count(&details)
	db.Model(&models.LogProduk{}).Count(&logs)
	if transactions != 1 || details != 1 || logs != 1 {
		t.Fatalf("transactions=%d details=%d logs=%d", transactions, details, logs)
	}
}

func TestCreateTransactionRollsBackWhenAnyItemFails(t *testing.T) {
	db, user, alamat, products := transactionFixture(t)
	body := []byte(`{"method_bayar":"bca","alamat_kirim":` + strconv.FormatUint(uint64(alamat.ID), 10) + `,"detail_trx":[{"product_id":` + strconv.FormatUint(uint64(products[0].ID), 10) + `,"kuantitas":2},{"product_id":` + strconv.FormatUint(uint64(products[1].ID), 10) + `,"kuantitas":99}]}`)
	req := httptest.NewRequest("POST", "/trx", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, user.ID))
	resp, err := transactionApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var product models.Produk
	if err := db.First(&product, products[0].ID).Error; err != nil {
		t.Fatal(err)
	}
	if product.Stok != 5 {
		t.Fatalf("stock changed despite rollback: %d", product.Stok)
	}
	var transactions, details, logs int64
	db.Model(&models.Transaksi{}).Count(&transactions)
	db.Model(&models.DetailTransaksi{}).Count(&details)
	db.Model(&models.LogProduk{}).Count(&logs)
	if transactions != 0 || details != 0 || logs != 0 {
		t.Fatalf("rollback left rows: transactions=%d details=%d logs=%d", transactions, details, logs)
	}
}

func TestCreateTransactionRejectsAnotherUsersAddress(t *testing.T) {
	db, user, _, products := transactionFixture(t)
	other := models.User{NamaUser: "Other", Email: "other@test.local", NoTelp: "08222", KataSandi: "hash", Role: "user", TanggalLahir: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	otherAddress := models.Alamat{JudulAlamat: "Other", NamaPenerima: "Other", NoTelp: "08222", DetailAlamat: "Other", UserID: other.ID}
	if err := db.Create(&otherAddress).Error; err != nil {
		t.Fatal(err)
	}

	body := []byte(`{"method_bayar":"bca","alamat_kirim":` + strconv.FormatUint(uint64(otherAddress.ID), 10) + `,"detail_trx":[{"product_id":` + strconv.FormatUint(uint64(products[0].ID), 10) + `,"kuantitas":1}]}`)
	req := httptest.NewRequest("POST", "/trx", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken(t, user.ID))
	resp, err := transactionApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status=%d", resp.StatusCode)
	}

	var product models.Produk
	db.First(&product, products[0].ID)
	if product.Stok != 5 {
		t.Fatalf("stock changed: %d", product.Stok)
	}
}
