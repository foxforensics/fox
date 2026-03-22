package vt

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
	"github.com/xxtea/xxtea-go/xxtea"
)

const FilesUrl = "https://www.virustotal.com/api/v3/files/%s"

/*
const (
	Clean      = "clean"
	Unknown    = "unknown"
	Unrated    = "unrated"
	Suspicious = "suspicious"
)
*/

// Encrypted reserve keys for emergency use
const (
	ReserveKey1 = "47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8"
	ReserveKey2 = "44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117"
)

/*
var badCategories = []string{
	"malicious",
	"suspicious",
}
*/

type FileReport struct {
	Data Data `json:"data"`
}

// Data represents the data object containing file attributes
type Data struct {
	Attributes Attributes `json:"attributes"`
	// ID         string     `json:"id"`
	// Links      Links      `json:"links"`
	// Type string `json:"type"`
}

/*
// Links represents the links object
type Links struct {
	Self string `json:"self"`
}
*/

// Attributes contains all file scan attributes
type Attributes struct {
	CapabilitiesTags []string `json:"capabilities_tags"`
	CreationDate     int64    `json:"creation_date"`
	// CrowdsourcedIdsResults  []CrowdsourcedIdsResult  `json:"crowdsourced_ids_results"`
	// CrowdsourcedIdsStats    CrowdsourcedIdsStats     `json:"crowdsourced_ids_stats"`
	// CrowdsourcedYaraResults []CrowdsourcedYaraResult `json:"crowdsourced_yara_results"`
	// Downloadable            bool                      `json:"downloadable"`
	FirstSubmissionDate  int64                     `json:"first_submission_date"`
	LastAnalysisDate     int64                     `json:"last_analysis_date"`
	LastAnalysisResults  map[string]AnalysisResult `json:"last_analysis_results"`
	LastAnalysisStats    AnalysisStats             `json:"last_analysis_stats"`
	LastModificationDate int64                     `json:"last_modification_date"`
	LastSubmissionDate   int64                     `json:"last_submission_date"`
	// MD5                         string                      `json:"md5"`
	// MeaningfulName string `json:"meaningful_name"`
	// Names                       []string                    `json:"names"`
	PopularThreatClassification PopularThreatClassification `json:"popular_threat_classification"` // NEW
	Reputation                  int                         `json:"reputation"`
	// SandboxVerdicts             map[string]SandboxVerdict   `json:"sandbox_verdicts"`
	// SHA1                        string                      `json:"sha1"`
	// SHA256 string `json:"sha256"`
	// SigmaAnalysisSummary map[string]SigmaStats `json:"sigma_analysis_summary"`
	// SigmaAnalysisStats   SigmaStats            `json:"sigma_analysis_stats"`
	// SigmaAnalysisResults []SigmaAnalysisResult `json:"sigma_analysis_results"`
	// Size            int64      `json:"size"`
	// Tags            []string   `json:"tags"`
	TimesSubmitted int        `json:"times_submitted"`
	TotalVotes     TotalVotes `json:"total_votes"`
	// TypeDescription string     `json:"type_description"`
	// TypeTag         string     `json:"type_tag"`
	UniqueSources int `json:"unique_sources"`
	// Vhash           string     `json:"vhash"`
}

/*
// CrowdsourcedIdsResult represents IDS detection results
type CrowdsourcedIdsResult struct {
	AlertContext  []AlertContext `json:"alert_context"`
	AlertSeverity string         `json:"alert_severity"`
	RuleCategory  string         `json:"rule_category"`
	RuleID        string         `json:"rule_id"`
	RuleMsg       string         `json:"rule_msg"`
	RuleSource    string         `json:"rule_source"`
}

// AlertContext represents network alert context
type AlertContext struct {
	Proto   string `json:"proto"`
	SrcIP   string `json:"src_ip"`
	SrcPort int    `json:"src_port"`
}

// CrowdsourcedIdsStats contains IDS statistics
type CrowdsourcedIdsStats struct {
	High   int `json:"high"`
	Info   int `json:"info"`
	Low    int `json:"low"`
	Medium int `json:"medium"`
}

// CrowdsourcedYaraResult represents YARA rule matches
type CrowdsourcedYaraResult struct {
	Description    string `json:"description"`
	MatchInSubfile bool   `json:"match_in_subfile"`
	RuleName       string `json:"rule_name"`
	RulesetID      string `json:"ruleset_id"`
	RulesetName    string `json:"ruleset_name"`
	Source         string `json:"source"`
}
*/

// AnalysisResult represents a single antivirus engine result
type AnalysisResult struct {
	Category      string  `json:"category"`
	EngineName    string  `json:"engine_name"`
	EngineUpdate  string  `json:"engine_update"`
	EngineVersion string  `json:"engine_version"`
	Method        string  `json:"method"`
	Result        *string `json:"result"` // Pointer to handle null values
}

// AnalysisStats contains analysis statistics
type AnalysisStats struct {
	ConfirmedTimeout int `json:"confirmed-timeout"`
	Failure          int `json:"failure"`
	Harmless         int `json:"harmless"`
	Malicious        int `json:"malicious"`
	Suspicious       int `json:"suspicious"`
	Timeout          int `json:"timeout"`
	TypeUnsupported  int `json:"type-unsupported"`
	Undetected       int `json:"undetected"`
}

/*
// SandboxVerdict represents sandbox analysis results
type SandboxVerdict struct {
	Category              string   `json:"category"`
	Confidence            int      `json:"confidence"`
	MalwareClassification []string `json:"malware_classification"`
	MalwareNames          []string `json:"malware_names"`
	SandboxName           string   `json:"sandbox_name"`
}

// SigmaStats contains sigma rule statistics
type SigmaStats struct {
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Critical int `json:"critical"`
	Low      int `json:"low"`
}

// SigmaAnalysisResult represents sigma rule matches
type SigmaAnalysisResult struct {
	RuleTitle       string         `json:"rule_title"`
	RuleSource      string         `json:"rule_source"`
	MatchContext    []MatchContext `json:"match_context"`
	RuleLevel       string         `json:"rule_level"`
	RuleDescription string         `json:"rule_description"`
	RuleAuthor      string         `json:"rule_author"`
	RuleID          string         `json:"rule_id"`
}

// MatchContext contains sigma match context
type MatchContext struct {
	Values map[string]string `json:"values"`
}
*/

// TotalVotes represents community votes
type TotalVotes struct {
	Harmless  int `json:"harmless"`
	Malicious int `json:"malicious"`
}

// PopularThreatClassification contains crowd-sourced threat classification
type PopularThreatClassification struct {
	SuggestedThreatLabel string `json:"suggested_threat_label"`
	// PopularThreatCategory []ThreatCount `json:"popular_threat_category"`
	// PopularThreatName     []ThreatCount `json:"popular_threat_name"`
}

/*
// ThreatCount represents a threat category/name with occurrence count
type ThreatCount struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}
*/

func (fr *FileReport) String() string {
	return fr.Data.Attributes.PopularThreatClassification.SuggestedThreatLabel
}

func (fr *FileReport) ToJSON() string {
	b, _ := json.MarshalIndent(fr, "", "  ")
	return string(b)
}

func (fr *FileReport) ToJSONL() string {
	b, _ := json.Marshal(fr)
	return string(b)
}

func CheckHash(sha, key string) *FileReport {
	req, _ := http.NewRequest("GET", fmt.Sprintf(FilesUrl, sha), nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", key)
	res, err := client.Default().Do(req)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer func() {
		_ = res.Body.Close()
	}()

	rep := new(FileReport)

	if err = json.NewDecoder(res.Body).Decode(&rep); err != nil {
		log.Println(err)
		return nil
	}

	/*
		res.Stats.Bad = countStats(obj, badCategories)
		res.Stats.All = countStats(obj, []string{
			"malicious",
			"suspicious",
			"undetected",
			"harmless",
			"timeout",
			"confirmed-timeout",
			"failure",
			"type-unsupported",
		})

		res.Verdict, _ = obj.GetString("popular_threat_classification.suggested_threat_label")

		if len(res.Verdict) == 0 {
			switch {
			case res.Stats.Bad > 0:
				res.Verdict = api.Suspicious
			case res.Stats.All > 0:
				res.Verdict = api.Clean
			default:
				res.Verdict = api.Unrated
			}
		}
	*/

	return rep
}

func Decrypt(s, k string) string {
	v, _ := hex.DecodeString(s)
	return string(xxtea.Decrypt(v, []byte(k)))
}
