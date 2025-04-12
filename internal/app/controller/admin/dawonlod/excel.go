package dawonlod

import (
	"github.com/xuri/excelize/v2"
)

func WriteToEXCEL(filename string, data [][]string) error {
	f := excelize.NewFile()

	sheetName := "Sheet1"
	for i, row := range data {
		for j, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return err
			}
			f.SetCellValue(sheetName, cellName, cell)
		}
	}

	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}