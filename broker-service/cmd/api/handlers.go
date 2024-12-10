package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
	"os"
	"path/filepath"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}

	err := app.readJson(w, r, &requestData)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	excelPath, err := filepath.Abs("ParametreFormu.xlsx")
	if err != nil {
		app.errorJson(w, fmt.Errorf("error getting absolute path: %v", err))
		return
	}
	currentDir, err := os.Getwd()
	if _, err := os.Stat(excelPath); os.IsNotExist(err) {
		app.errorJson(w, fmt.Errorf("excel file not found at path: %s %v", excelPath, currentDir))
		return
	}

	if err != nil {
		app.errorJson(w, fmt.Errorf("error getting current directory: %v %v %v", err, currentDir, excelPath))
		return
	}
	fmt.Printf("Current working directory: %s\n", currentDir)
	fmt.Printf("Trying to open file at: %s\n", excelPath)

	var f *excelize.File
	f, err = excelize.OpenFile(excelPath)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer f.Close() // Don't forget to close the file

	// Rest of your existing code...
	sheetName := "Sheet1"
	rows, _ := f.GetRows(sheetName)
	newRow := len(rows) + 1

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", newRow), requestData["agirlik"])
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", newRow), requestData["sogutmaSuresi"])
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", newRow), requestData["urunAdi"])
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", newRow), requestData["malAlmaZamani"])

	if err := f.SaveAs(excelPath); err != nil {
		app.errorJson(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "Broker verileri başarıyla aldı",
		Data:    requestData,
	}

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err)
	}
}
