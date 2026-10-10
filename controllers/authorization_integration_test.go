package controllers

import (
	"evermos-mini/config"
	"evermos-mini/middleware"
	"evermos-mini/models"
	"evermos-mini/utils"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func integrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_DSN to run MySQL integration tests")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("integration database unavailable: %v", err)
	}
	if err := db.Migrator().DropTable(&models.DetailTransaksi{}, &models.Transaksi{}, &models.LogProduk{}, &models.Produk{}, &models.Category{}, &models.Alamat{}, &models.Toko{}, &models.User{}); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Toko{}, &models.Alamat{}, &models.Category{}, &models.Produk{}, &models.LogProduk{}, &models.Transaksi{}, &models.DetailTransaksi{}); err != nil {
		t.Fatal(err)
	}
	config.DB = db
	return db
}

func testToken(t *testing.T, userID uint) string {
	t.Helper()
	t.Setenv("JWT_SECRET", "integration-test-secret")
	token, err := utils.GenerateJWT(userID, "user")
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestAddressOwnership(t *testing.T) {
	db := integrationDB(t)
	users := []models.User{
		{NamaUser: "User A", Email: "a@test.local", NoTelp: "08001", KataSandi: "hash", Role: "user", TanggalLahir: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{NamaUser: "User B", Email: "b@test.local", NoTelp: "08002", KataSandi: "hash", Role: "user", TanggalLahir: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	address := models.Alamat{JudulAlamat: "Rumah B", NamaPenerima: "B", NoTelp: "08002", DetailAlamat: "Alamat B", UserID: users[1].ID}
	if err := db.Create(&address).Error; err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	app.Get("/user/alamat/:id", middleware.JWTProtected, GetAlamatByID)
	req := httptest.NewRequest("GET", "/user/alamat/"+strconv.FormatUint(uint64(address.ID), 10), nil)
	req.Header.Set("Authorization", "Bearer "+testToken(t, users[0].ID))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("user A accessed user B address, status=%d", resp.StatusCode)
	}
}
