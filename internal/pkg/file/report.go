package file

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/xxtea/xxtea-go/xxtea"

	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
)

const api = "https://www.virustotal.com/api/v3/files/%s"

var keys = [2]string{
	"47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8",
	"44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117",
}

type report struct {
	Data struct {
		Report Report `json:"attributes"`
	} `json:"data"`
}

type Report struct {
	PopularThreatClassification struct {
		SuggestedThreatLabel string `json:"suggested_threat_label"`
	} `json:"popular_threat_classification"`
	FirstSubmissionDate int64 `json:"first_submission_date"`
	LastAnalysisDate    int64 `json:"last_analysis_date"`
	LastAnalysisStats   struct {
		ConfirmedTimeout int `json:"confirmed-timeout"`
		Failure          int `json:"failure"`
		Harmless         int `json:"harmless"`
		Malicious        int `json:"malicious"`
		Suspicious       int `json:"suspicious"`
		Timeout          int `json:"timeout"`
		TypeUnsupported  int `json:"type-unsupported"`
		Undetected       int `json:"undetected"`
	} `json:"last_analysis_stats"`
	LastAnalysisResults map[string]struct {
		Category      string  `json:"category"`
		EngineName    string  `json:"engine_name"`
		EngineUpdate  string  `json:"engine_update"`
		EngineVersion string  `json:"engine_version"`
		Method        string  `json:"method"`
		Result        *string `json:"result"`
	} `json:"last_analysis_results"`
}

func (rep *Report) String() string {
	return rep.PopularThreatClassification.SuggestedThreatLabel
}

func (rep *Report) ToJSON() string {
	b, _ := json.MarshalIndent(rep, "", "  ")
	return string(b)
}

func (rep *Report) ToJSONL() string {
	b, _ := json.Marshal(rep)
	return string(b)
}

func GetReport(sum, key string) *Report {
	req, _ := http.NewRequest("GET", fmt.Sprintf(api, sum), nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", key)

	res, err := client.Default().Do(req)

	if err != nil {
		log.Println(err)
		return nil
	}

	rep := new(report)

	defer func() {
		_ = res.Body.Close()
	}()

	if err = json.NewDecoder(res.Body).Decode(&rep); err != nil {
		log.Println(err)
		return nil
	}

	return &rep.Data.Report
}

func ReserveKey(n int, key string) string {
	b, _ := hex.DecodeString(keys[n-1])
	b = xxtea.Decrypt(b, []byte(key))

	return string(b)
}
