package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EMSIFA Base URL
const baseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"

var wilayahClient = &http.Client{Timeout: 10 * time.Second}

// GET /provcity/listprovincies
func GetListProvincies(c *fiber.Ctx) error {
	resp, err := wilayahClient.Get(fmt.Sprintf("%s/provinces.json", baseURL))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal mengambil data provinsi"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(fiber.Map{"status": false, "message": "Layanan wilayah bermasalah"})
	}

	body, _ := io.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return c.JSON(fiber.Map{"status": true, "data": data})
}

// GET /provcity/listcities/:prov_id
func GetListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	resp, err := wilayahClient.Get(fmt.Sprintf("%s/regencies/%s.json", baseURL, provID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal mengambil data kota"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(fiber.Map{"status": false, "message": "Layanan wilayah bermasalah"})
	}

	body, _ := io.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return c.JSON(fiber.Map{"status": true, "data": data})
}

// GET /provcity/detailprovince/:prov_id
func GetDetailProvince(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	resp, err := wilayahClient.Get(fmt.Sprintf("%s/province/%s.json", baseURL, provID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal mengambil detail provinsi"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(fiber.Map{"status": false, "message": "Layanan wilayah bermasalah"})
	}

	body, _ := io.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return c.JSON(fiber.Map{"status": true, "data": data})
}

// GET /provcity/detailcity/:city_id
func GetDetailCity(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	resp, err := wilayahClient.Get(fmt.Sprintf("%s/regency/%s.json", baseURL, cityID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal mengambil detail kota"})
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(fiber.Map{"status": false, "message": "Layanan wilayah bermasalah"})
	}

	body, _ := io.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return c.JSON(fiber.Map{"status": true, "data": data})
}
