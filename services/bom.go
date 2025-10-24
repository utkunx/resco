package services

import (
	"database/sql"
	"fmt"
	"os"
	"resco/db"
	"strings"
	"time"
)

type BOMResult struct {
	BOMRecCode      string  `json:"parent-number"`
	AD              string  `json:"parent-name"`
	ParProSpec      string  `json:"par_pro_spec"`
	BOMRecKaynakCode string `json:"child-number"`
	SubItemName     *string `json:"child-name"`
	SubProSpec      string  `json:"sub_pro_spec"`
	BOMRecKaynak0   float64 `json:"child-quantity"`
	Depth           int     `json:"depth"`
}

type BOMResultCombined struct {
	BOMRecCode      string  `json:"parent-number"`
	AD              string  `json:"parent-name"`
	ADChinese       string  `json:"parent-name-cn"`
	ParProSpec      string  `json:"par_pro_spec"`
	BOMRecKaynakCode string `json:"child-number"`
	SubItemName     *string `json:"child-name"`
	SubItemNameChinese *string `json:"child-name-cn"`
	SubProSpec      string  `json:"sub_pro_spec"`
	BOMRecKaynak0   float64 `json:"child-quantity"`
	Depth           int     `json:"depth"`
}

type BOMTotalResult struct {
	SequenceNumber int    `json:"sequence-number"`
	Code           string `json:"code"`
}

type ProductCheckResult struct {
	SequenceNumber int    `json:"sequence-number"`
	Code           string `json:"code"`
	Status         string `json:"status"` // "OK" or "NOT"
}

// GetBOMByCode executes the recursive BOM query for a given item code
func GetBOMByCode(itemCode string) ([]BOMResult, error) {
	// Read the SQL file
	sqlContent, err := os.ReadFile("000.sql")
	if err != nil {
		return nil, fmt.Errorf("error reading SQL file: %v", err)
	}

	// Replace the hardcoded item code with the provided one
	sqlQuery := strings.Replace(string(sqlContent), "'360004'", fmt.Sprintf("'%s'", itemCode), 1)

	// Execute the query
	rows, err := db.DB.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []BOMResult
	for rows.Next() {
		var result BOMResult
		err := rows.Scan(
			&result.BOMRecCode,
			&result.AD,
			&result.ParProSpec,
			&result.BOMRecKaynakCode,
			&result.SubItemName,
			&result.SubProSpec,
			&result.BOMRecKaynak0,
			&result.Depth,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return results, nil
}

// GetBOMByCodeParameterized executes the recursive BOM query using parameterized query
func GetBOMByCodeParameterized(itemCode string) ([]BOMResult, error) {
	// Split the SQL into executable batches
	sqlBatch1 := `
	IF OBJECT_ID('tempdb..#TempRecursiveResults') IS NOT NULL
		DROP TABLE #TempRecursiveResults;

	IF OBJECT_ID('tempdb..#TempReco') IS NOT NULL
		DROP TABLE #TempReco;

	WITH RecursiveSearch AS (
		SELECT EVRAKNO, TRNUM, SRNUM, BOMREC_SIRANO, BOMREC_CODE, BOMREC_KAYNAKCODE, BOMREC_KAYNAK0, TLOG_USERNAME, TLOG_LOGTARIH, TLOG_PSTATION, GK_2, 1 AS Depth
		FROM RESCO_2019.dbo.BOMU01T
		WHERE BOMREC_CODE = @p1 AND BOMREC_INPUTTYPE='H'

		UNION ALL

		SELECT YT.EVRAKNO, YT.TRNUM, YT.SRNUM, YT.BOMREC_SIRANO, YT.BOMREC_CODE, YT.BOMREC_KAYNAKCODE, YT.BOMREC_KAYNAK0, YT.TLOG_USERNAME, YT.TLOG_LOGTARIH, YT.TLOG_PSTATION, YT.GK_2, RS.Depth + 1
		FROM RESCO_2019.dbo.BOMU01T YT
		INNER JOIN RecursiveSearch RS ON YT.BOMREC_CODE = RS.BOMREC_KAYNAKCODE
		WHERE RS.Depth < 10 AND YT.BOMREC_INPUTTYPE='H'
	)
	SELECT * INTO #TempRecursiveResults FROM RecursiveSearch
	ORDER BY Depth ASC,EVRAKNO ASC, SRNUM ASC;

	SELECT TRR.BOMREC_CODE,RT.AD,
	CAST('' AS NVARCHAR(255)) AS ParProSpec,
	TRR.BOMREC_KAYNAKCODE,
	CAST(NULL AS NVARCHAR(255)) AS SubItemName,
	CAST('' AS NVARCHAR(255)) AS SubProSpec,
	TRR.BOMREC_KAYNAK0,
	TRR.Depth,
	TRR.EVRAKNO, TRR.SRNUM
	INTO #TempReco
	FROM #TempRecursiveResults TRR
	LEFT JOIN RESCO_2019.dbo.STOK00 RT ON TRR.BOMREC_CODE = RT.KOD
	ORDER BY Depth ASC,EVRAKNO ASC, SRNUM ASC;

	UPDATE T
	SET T.SubItemName = R.AD
	FROM #TempReco T
	LEFT JOIN RESCO_2019.dbo.STOK00 R ON T.BOMREC_KAYNAKCODE = R.KOD;

	ALTER TABLE #TempReco
	ALTER COLUMN BOMREC_CODE VARCHAR(24);

	ALTER TABLE #TempReco
	ALTER COLUMN BOMREC_KAYNAKCODE VARCHAR(24);

	UPDATE #TempReco
	SET BOMREC_CODE = TRIM(BOMREC_CODE),
		AD = TRIM(AD),
		SubItemName = LTRIM(RTRIM(SubItemName)),
		BOMREC_KAYNAKCODE = TRIM(BOMREC_KAYNAKCODE)
	WHERE BOMREC_CODE IS NOT NULL OR AD IS NOT NULL OR SubItemName IS NOT NULL OR BOMREC_KAYNAKCODE IS NOT NULL;

	SELECT BOMREC_CODE, AD, ParProSpec,BOMREC_KAYNAKCODE, SubItemName, SubProSpec,BOMREC_KAYNAK0,Depth FROM #TempReco
	ORDER BY Depth ASC,EVRAKNO ASC, SRNUM ASC;

	DROP TABLE #TempRecursiveResults;
	DROP TABLE #TempReco;
	`

	rows, err := db.DB.Query(sqlBatch1, sql.Named("p1", itemCode))
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []BOMResult
	for rows.Next() {
		var result BOMResult
		err := rows.Scan(
			&result.BOMRecCode,
			&result.AD,
			&result.ParProSpec,
			&result.BOMRecKaynakCode,
			&result.SubItemName,
			&result.SubProSpec,
			&result.BOMRecKaynak0,
			&result.Depth,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return results, nil
}

// GetBOMByCodeWithTranslation executes the recursive BOM query and applies Chinese translations
func GetBOMByCodeWithTranslation(itemCode string) ([]BOMResult, error) {
	// Get the BOM data
	results, err := GetBOMByCodeParameterized(itemCode)
	if err != nil {
		return nil, err
	}

	// Load translations if not already loaded
	err = LoadTranslations()
	if err != nil {
		return nil, fmt.Errorf("error loading translations: %v", err)
	}

	// Load fallback translations if not already loaded
	err = LoadFallbackTranslations()
	if err != nil {
		return nil, fmt.Errorf("error loading fallback translations: %v", err)
	}

	// Apply translations (with fallback support)
	translatedResults := ApplyTranslationsToBOM(results)

	return translatedResults, nil
}

// GetBOMByCodeWithTranslationTracking executes the recursive BOM query and applies Chinese translations
// Returns translated results and a slice of item codes that failed to translate
func GetBOMByCodeWithTranslationTracking(itemCode string) ([]BOMResult, []string, error) {
	// Get the BOM data
	results, err := GetBOMByCodeParameterized(itemCode)
	if err != nil {
		return nil, nil, err
	}

	// Load translations if not already loaded
	err = LoadTranslations()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading translations: %v", err)
	}

	// Load fallback translations if not already loaded
	err = LoadFallbackTranslations()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading fallback translations: %v", err)
	}

	// Apply translations and track failures
	translatedResults, untranslatedCodes := ApplyTranslationsToBOMWithTracking(results)

	return translatedResults, untranslatedCodes, nil
}

// GetBOMByCodeCombined executes the recursive BOM query and returns both Turkish and Chinese
func GetBOMByCodeCombined(itemCode string) ([]BOMResultCombined, error) {
	// Get the BOM data
	results, err := GetBOMByCodeParameterized(itemCode)
	if err != nil {
		return nil, err
	}

	// Load translations if not already loaded
	err = LoadTranslations()
	if err != nil {
		return nil, fmt.Errorf("error loading translations: %v", err)
	}

	// Load fallback translations if not already loaded
	err = LoadFallbackTranslations()
	if err != nil {
		return nil, fmt.Errorf("error loading fallback translations: %v", err)
	}

	// Create combined results with both Turkish and Chinese
	combinedResults := make([]BOMResultCombined, len(results))

	for i, result := range results {
		combinedResults[i] = BOMResultCombined{
			BOMRecCode:      result.BOMRecCode,
			AD:              result.AD,
			ADChinese:       TranslateWithFallback(result.AD, result.BOMRecCode),
			ParProSpec:      result.ParProSpec,
			BOMRecKaynakCode: result.BOMRecKaynakCode,
			SubItemName:     result.SubItemName,
			SubProSpec:      result.SubProSpec,
			BOMRecKaynak0:   result.BOMRecKaynak0,
			Depth:           result.Depth,
		}

		// Translate child name if it exists
		if result.SubItemName != nil && *result.SubItemName != "" {
			translated := TranslateWithFallback(*result.SubItemName, result.BOMRecKaynakCode)
			combinedResults[i].SubItemNameChinese = &translated
		}
	}

	return combinedResults, nil
}

// GetBOMByCodeCombinedWithTracking executes the recursive BOM query and returns both Turkish and Chinese
// Returns combined results and a slice of item codes that failed to translate
func GetBOMByCodeCombinedWithTracking(itemCode string) ([]BOMResultCombined, []string, error) {
	// Get the BOM data
	results, err := GetBOMByCodeParameterized(itemCode)
	if err != nil {
		return nil, nil, err
	}

	// Load translations if not already loaded
	err = LoadTranslations()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading translations: %v", err)
	}

	// Load fallback translations if not already loaded
	err = LoadFallbackTranslations()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading fallback translations: %v", err)
	}

	// Track untranslated codes
	untranslatedCodesMap := make(map[string]bool)
	var untranslatedCodes []string

	// Create combined results with both Turkish and Chinese
	combinedResults := make([]BOMResultCombined, len(results))

	for i, result := range results {
		// Translate parent name
		parentTranslated, parentSuccess := TranslateWithFallbackTracking(result.AD, result.BOMRecCode)
		if !parentSuccess && result.BOMRecCode != "" {
			if !untranslatedCodesMap[result.BOMRecCode] {
				untranslatedCodesMap[result.BOMRecCode] = true
				untranslatedCodes = append(untranslatedCodes, result.BOMRecCode)
			}
		}

		combinedResults[i] = BOMResultCombined{
			BOMRecCode:      result.BOMRecCode,
			AD:              result.AD,
			ADChinese:       parentTranslated,
			ParProSpec:      result.ParProSpec,
			BOMRecKaynakCode: result.BOMRecKaynakCode,
			SubItemName:     result.SubItemName,
			SubProSpec:      result.SubProSpec,
			BOMRecKaynak0:   result.BOMRecKaynak0,
			Depth:           result.Depth,
		}

		// Translate child name if it exists
		if result.SubItemName != nil && *result.SubItemName != "" {
			childTranslated, childSuccess := TranslateWithFallbackTracking(*result.SubItemName, result.BOMRecKaynakCode)
			combinedResults[i].SubItemNameChinese = &childTranslated

			if !childSuccess && result.BOMRecKaynakCode != "" {
				if !untranslatedCodesMap[result.BOMRecKaynakCode] {
					untranslatedCodesMap[result.BOMRecKaynakCode] = true
					untranslatedCodes = append(untranslatedCodes, result.BOMRecKaynakCode)
				}
			}
		}
	}

	return combinedResults, untranslatedCodes, nil
}

// GetBOMTotal executes the recursive BOM query and returns unique codes with sequential numbers
func GetBOMTotal(itemCode string) ([]BOMTotalResult, error) {
	// Get the BOM data
	results, err := GetBOMByCodeParameterized(itemCode)
	if err != nil {
		return nil, err
	}

	// Use a map to track unique codes
	uniqueCodes := make(map[string]bool)
	var orderedCodes []string

	// Collect all parent and child codes
	for _, result := range results {
		// Add parent code if not already present
		if !uniqueCodes[result.BOMRecCode] {
			uniqueCodes[result.BOMRecCode] = true
			orderedCodes = append(orderedCodes, result.BOMRecCode)
		}

		// Add child code if not already present
		if !uniqueCodes[result.BOMRecKaynakCode] {
			uniqueCodes[result.BOMRecKaynakCode] = true
			orderedCodes = append(orderedCodes, result.BOMRecKaynakCode)
		}
	}

	// Create results with sequential numbers
	totalResults := make([]BOMTotalResult, len(orderedCodes))
	for i, code := range orderedCodes {
		totalResults[i] = BOMTotalResult{
			SequenceNumber: i + 1,
			Code:           code,
		}
	}

	return totalResults, nil
}

// CheckProducts checks all BOM products against Heihu API with rate limiting
func CheckProducts(itemCode string) ([]ProductCheckResult, error) {
	// Step 1: Get all unique codes from BOM
	bomTotal, err := GetBOMTotal(itemCode)
	if err != nil {
		return nil, fmt.Errorf("error getting BOM total: %v", err)
	}

	// Step 2: Check each code against Heihu API with rate limiting
	results := make([]ProductCheckResult, len(bomTotal))

	for i, item := range bomTotal {
		// Query Heihu API
		_, err := QueryHeihu(item.Code)

		// Determine status based on error
		status := "OK"
		if err != nil {
			// Check if error message contains "no product found with exact code"
			if strings.Contains(err.Error(), "no product found with exact code") {
				status = "NOT"
			} else {
				// For other errors, we still mark as NOT but could log the actual error
				status = "NOT"
			}
		}

		results[i] = ProductCheckResult{
			SequenceNumber: item.SequenceNumber,
			Code:           item.Code,
			Status:         status,
		}

		// Rate limiting: Wait 100ms between requests to stay under 20 QPS limit
		// 100ms delay = 10 requests/second (well under the 20 QPS limit)
		if i < len(bomTotal)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return results, nil
}