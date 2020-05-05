package rest

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
	"strings"
)

type GetCodeRunTestInput struct {
	OrgId   int32 `webcore:"orgId"`
	RunId   int64 `webcore:"runId"`
	Summary bool  `webcore:"summary"`
}

func getCodeRunTest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeRunTestInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if inputs.Summary {
		summary, err := database.GetCodeRunTestSummary(inputs.RunId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get test summary: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonWriter.Encode(summary)
	} else {
		core.Warning("Non summary not supported yet.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type ExportTestInput struct {
	OrgId int32 `webcore:"orgId"`
	RunId int64 `webcore:"runId"`
}

func exportTests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	xlsx := excelize.NewFile()

	inputs := ExportTestInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Failed to find org: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	headerStyle := excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	}

	headerStyleIdx, err := xlsx.NewStyle(&headerStyle)
	if err != nil {
		core.Warning("Failed to create style: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create one worksheet for each data source and then one worksheet for the tests.
	// The data that's stored in the database is encrypted so we need to decrypt it before
	// sending it to the user. Since there can be a lot of sources and data, we do a batch
	// decrypt to minimize the back and forth w/ Vault.
	sources, err := database.GetTrackedSourcesForRunId(inputs.RunId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to retrieve tracked sources: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	totalData := 0
	sourceData := map[int64][]*core.TrackedData{}
	for _, s := range sources {
		sourceData[s.Id], err = database.GetTrackedDataForSource(s.Id, role)
		if err != nil {
			core.Warning("Failed to retrieve tracked data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		totalData = totalData + len(sourceData[s.Id])
	}

	srcEncrypted := make([][]byte, len(sources))
	dataEncrypted := make([][]byte, totalData)

	dataIdx := 0
	for idx, s := range sources {
		srcEncrypted[idx] = []byte(s.Src)

		for _, d := range sourceData[s.Id] {
			dataEncrypted[dataIdx] = []byte(d.Data)
			dataIdx = dataIdx + 1
		}
	}

	srcRaw, err := vault.BatchTransitDecrypt("scripts", srcEncrypted)
	if err != nil {
		core.Warning("Failed to decrypt src: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataRaw, err := vault.BatchTransitDecrypt("scripts", dataEncrypted)
	if err != nil {
		core.Warning("Failed to decrypt data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sourceToSheetName := map[int64]string{}
	dataIdToReference := map[int64]string{}
	dataIdFieldToColumn := map[int64]map[string]string{}

	dataIdx = 0
	for idx, s := range sources {
		clientData, err := database.GetClientDataFromId(s.DataId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get client data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sheetName := fmt.Sprintf("Data - %s", clientData.Data.Name)
		sourceToSheetName[s.Id] = sheetName
		xlsx.NewSheet(sheetName)

		xlsx.SetCellStr(sheetName, "A1", "Source Method: ")
		xlsx.SetCellStr(sheetName, "B1", string(srcRaw[idx]))

		xlsx.SetCellStr(sheetName, "A2", "GRCHive Data: ")
		xlsx.SetCellStr(sheetName, "B2", webcore.MustGetRouteUrlAbsolute(
			webcore.SingleClientDataRouteName,
			core.DashboardOrgOrgQueryId,
			org.OktaGroupName,
			core.DashboardOrgClientDataQueryId,
			strconv.FormatInt(s.DataId, 10),
		))

		headerRow := 4
		currentDataRow := 5
		fieldToColumn := map[string]string{}

		for didx, d := range sourceData[s.Id] {
			dataValue := map[string]interface{}{}
			err = json.Unmarshal(dataRaw[dataIdx], &dataValue)
			if err != nil {
				core.Warning("Failed to parse data: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Special case for index = 0 where we need to populate the
			// column headers. Assume that all the other data entries will
			// have the same format as the 0th index.
			if didx == 0 {
				colIdx := 1
				for k, _ := range dataValue {
					fieldToColumn[k], err = excelize.ColumnNumberToName(colIdx)
					if err != nil {
						core.Warning("Invalid column number: " + err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					colIdx = colIdx + 1

					xlsx.SetCellStr(
						sheetName,
						fmt.Sprintf("%s%d", fieldToColumn[k], headerRow),
						k)
				}

				startCol, _ := excelize.ColumnNumberToName(1)
				endCol, _ := excelize.ColumnNumberToName(colIdx - 1)
				err = xlsx.SetCellStyle(sheetName,
					fmt.Sprintf("%s%d", startCol, headerRow),
					fmt.Sprintf("%s%d", endCol, headerRow),
					headerStyleIdx,
				)

				if err != nil {
					core.Warning("Failed to set header style: " + err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			for k, v := range dataValue {
				cell := fmt.Sprintf("%s%d", fieldToColumn[k], currentDataRow)

				err = xlsx.SetCellValue(sheetName, cell, v)
				if err != nil {
					// Assume that an error indicates that it isn't a supported format
					// so fallback to JSON formatting.
					jsonV, err := json.Marshal(v)
					if err != nil {
						core.Warning("Failed to marshal value: " + err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					xlsx.SetCellStr(
						sheetName,
						cell,
						string(jsonV),
					)
				}
			}

			dataIdToReference[d.Id] = fmt.Sprintf("'%s'!%%s%d", sheetName, currentDataRow)
			dataIdFieldToColumn[d.Id] = fieldToColumn
			currentDataRow = currentDataRow + 1
			dataIdx = dataIdx + 1
		}
	}

	tests, err := database.GetTrackedTestsForRun(inputs.RunId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to retrieve tracked tests: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	testSheetName := "Tests"
	xlsx.NewSheet(testSheetName)
	xlsx.SetCellStr(
		testSheetName,
		"A1",
		"Test Value",
	)

	xlsx.SetCellStr(
		testSheetName,
		"B1",
		"Test Field",
	)

	xlsx.SetCellStr(
		testSheetName,
		"C1",
		"Reference Value",
	)

	xlsx.SetCellStr(
		testSheetName,
		"D1",
		"Reference Field",
	)

	xlsx.SetCellStr(
		testSheetName,
		"E1",
		"Action",
	)

	xlsx.SetCellStr(
		testSheetName,
		"F1",
		"Success",
	)

	xlsx.SetCellStyle(testSheetName,
		"A1",
		"F1",
		headerStyleIdx,
	)

	testRow := 2
	for _, t := range tests {
		splitFields := strings.Split(t.Field, ",")

		if t.DataAId.NullInt64.Valid {
			id := t.DataAId.NullInt64.Int64
			xlsx.SetCellFormula(
				testSheetName,
				fmt.Sprintf("A%d", testRow),
				fmt.Sprintf(dataIdToReference[id], dataIdFieldToColumn[id][splitFields[0]]),
			)
		}

		if t.DataBId.NullInt64.Valid {
			id := t.DataBId.NullInt64.Int64
			xlsx.SetCellFormula(
				testSheetName,
				fmt.Sprintf("C%d", testRow),
				fmt.Sprintf(dataIdToReference[id], dataIdFieldToColumn[id][splitFields[1]]),
			)
		}

		xlsx.SetCellValue(
			testSheetName,
			fmt.Sprintf("B%d", testRow),
			splitFields[0],
		)

		xlsx.SetCellValue(
			testSheetName,
			fmt.Sprintf("D%d", testRow),
			splitFields[1],
		)

		xlsx.SetCellValue(
			testSheetName,
			fmt.Sprintf("E%d", testRow),
			t.Action,
		)

		xlsx.SetCellValue(
			testSheetName,
			fmt.Sprintf("F%d", testRow),
			t.Ok,
		)

		testRow = testRow + 1
	}

	xlsx.DeleteSheet("Sheet1")
	buf, err := xlsx.WriteToBuffer()
	if err != nil {
		core.Warning("Failed to save XLSX to buffer: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}
