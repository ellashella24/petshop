package service

import (
	"fmt"
	"petshop/entity"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportExcel(
	transactionData []entity.Transaction, transactionDetailData []entity.TransactionDetail,
	productData []entity.Product, email string,
) error {
	headers := map[string]string{
		"A1": "Invoice ID",
		"B1": "Nama Barang",
		"C1": "Quantity",
		"D1": "Paid At",
	}

	file := excelize.NewFile()

	activeSheet := file.NewSheet("Sheet1")

	styleHeader, _ := file.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Size:   10,
				Family: "Arial",
				Color:  "000000",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
		},
	)

	for i, v := range headers {
		file.SetCellValue("Sheet1", i, v)
	}

	file.SetCellStyle("Sheet1", "A1", "D1", styleHeader)

	for i := 0; i < len(transactionData); i++ {
		appendRow(file, i, transactionData, transactionDetailData, productData)
	}

	file.SetActiveSheet(activeSheet)

	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()

	filename := fmt.Sprint(
		"./hasil-export/transaction-report-store", productData[0].StoreID, "-", year, month, day, "-", hour, minute,
		second, ".xlsx",
	)

	filename = strings.ReplaceAll(filename, " ", "")
	fmt.Println(filename)
	err := file.SaveAs(filename)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = SendMail(filename, email)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func appendRow(
	file *excelize.File, index int, transactionData []entity.Transaction,
	transactionDetailData []entity.TransactionDetail, productData []entity.Product,
) (fileWriter *excelize.File) {
	rowCount := index + 2

	styleContent, _ := file.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Size:   10,
				Family: "Arial",
				Color:  "000000",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
		},
	)

	file.SetCellValue("Sheet1", fmt.Sprint("A", rowCount), fmt.Sprint(transactionData[index].InvoiceID))
	file.SetCellValue("Sheet1", fmt.Sprint("B", rowCount), fmt.Sprint(productData[index].Name))
	file.SetCellValue("Sheet1", fmt.Sprint("C", rowCount), fmt.Sprint(transactionDetailData[index].Quantity))
	file.SetCellValue("Sheet1", fmt.Sprint("D", rowCount), fmt.Sprint(transactionData[index].PaidAt))

	file.SetCellStyle("Sheet1", fmt.Sprint("A", rowCount), fmt.Sprint("D", rowCount), styleContent)

	return file
}
