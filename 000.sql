-- 1️⃣ Check if the temporary table exists, then drop it
IF OBJECT_ID('tempdb..#TempRecursiveResults') IS NOT NULL
    DROP TABLE #TempRecursiveResults;

-- 1️⃣ Check if the temporary table exists, then drop it
IF OBJECT_ID('tempdb..#TempReco') IS NOT NULL
    DROP TABLE #TempReco;


WITH RecursiveSearch AS (
    -- 1️⃣ First Query: Search for the Initial Reference in BOMREC_CODE
    SELECT EVRAKNO, TRNUM, SRNUM, BOMREC_SIRANO, BOMREC_CODE, BOMREC_KAYNAKCODE, BOMREC_KAYNAK0, TLOG_USERNAME, TLOG_LOGTARIH, TLOG_PSTATION, GK_2, 1 AS Depth
    FROM RESCO_2019.dbo.BOMU01T
    WHERE BOMREC_CODE = '360004' AND BOMREC_INPUTTYPE='H' -- Replace with your initial search value

    UNION ALL

    -- 2️⃣ Recursive Query: Search Next Level Using BOMREC_KAYNAKCODE
    SELECT YT.EVRAKNO, YT.TRNUM, YT.SRNUM, YT.BOMREC_SIRANO, YT.BOMREC_CODE, YT.BOMREC_KAYNAKCODE, YT.BOMREC_KAYNAK0, YT.TLOG_USERNAME, YT.TLOG_LOGTARIH, YT.TLOG_PSTATION, YT.GK_2, RS.Depth + 1
    FROM RESCO_2019.dbo.BOMU01T YT
    INNER JOIN RecursiveSearch RS ON YT.BOMREC_CODE = RS.BOMREC_KAYNAKCODE
    WHERE RS.Depth < 10 AND YT.BOMREC_INPUTTYPE='H'  -- Prevent infinite recursion
)
SELECT * INTO #TempRecursiveResults FROM RecursiveSearch
ORDER BY Depth ASC,EVRAKNO ASC, SRNUM ASC;




-- 3️⃣ Sonuçlara ReferenceTable'dan JOIN Yap
SELECT TRR.BOMREC_CODE,RT.AD,
CAST('' AS NVARCHAR(255)) AS ParProSpec, -- Varsayılan değer eklendi
 TRR.BOMREC_KAYNAKCODE,
CAST(NULL AS NVARCHAR(255)) AS SubItemName, -- Boş sütun eklendi
CAST('' AS NVARCHAR(255)) AS SubProSpec, -- Varsayılan değer eklendi
TRR.BOMREC_KAYNAK0,
TRR.Depth,
TRR.EVRAKNO, TRR.SRNUM
INTO #TempReco
FROM #TempRecursiveResults TRR
LEFT JOIN RESCO_2019.dbo.STOK00 RT ON TRR.BOMREC_CODE = RT.KOD 
ORDER BY Depth ASC,EVRAKNO ASC, SRNUM ASC;






-- 4️⃣ BOMREC_KAYNAKCODE'u ReferenceTable'daki ReferenceName ile eşleştir
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

-- 5️⃣ Drop the temporary table at the end to clean up
DROP TABLE #TempRecursiveResults;
DROP TABLE #TempReco;
