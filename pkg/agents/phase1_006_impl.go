// Auto-generated for phase1-006
// Task: Implement Security agent with risk-based scoring (70% baseline)

```go
package security

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// SecurityDecision represents the outcome of a security review
type SecurityDecision string

const (
	SecurityDecisionApprove            SecurityDecision = "APPROVE"
	SecurityDecisionConditionalApprove SecurityDecision = "CONDITIONAL_APPROVE"
	SecurityDecisionReject             SecurityDecision = "REJECT"
)

// RiskType represents categories of security risks
type RiskType string

const (
	RiskTypePrivilegeEscalation  RiskType = "privilege_escalation"
	RiskTypeDestructiveFiles     RiskType = "destructive_file_ops"
	RiskTypeNetworkAccess        RiskType = "uncontrolled_network"
	RiskTypeMissingErrorHandling RiskType = "missing_error_handling"
	RiskTypeMissingAuditLog      RiskType = "missing_audit_logging"
)

// SecurityEnhancement represents security improvements
type SecurityEnhancement string

const (
	SecurityEnhancementEncryption   SecurityEnhancement = "encryption_added"
	SecurityEnhancementAuditLogging SecurityEnhancement = "audit_logging_added"
	SecurityEnhancementHardening    SecurityEnhancement = "security_hardening"
)

// RiskAssessment represents an individual security risk
type RiskAssessment struct {
	RiskType    RiskType `json:"risk_type"`
	Severity    int      `json:"severity"`
	Description string   `json:"description"`
	Mitigation  string   `json:"mitigation,omitempty"`
}

// SecurityBonus represents a security enhancement
type SecurityBonus struct {
	EnhancementType SecurityEnhancement `json:"enhancement_type"`
	Points          int                 `json:"points"`
	Description     string              `json:"description"`
}

// SecurityReview represents a complete security review result
type SecurityReview struct {
	ProposalID string             `json:"proposal_id"`
	BaseScore  int                `json:"base_score"`
	Risks      []RiskAssessment   `json:"risks"`
	Bonuses    []SecurityBonus    `json:"bonuses"`
	FinalScore int                `json:"final_score"`
	Decision   SecurityDecision   `json:"decision"`
	Feedback   string             `json:"feedback"`
	Timestamp  string             `json:"timestamp"`
}

// RiskCalculator calculates security scores based on risks and bonuses
type RiskCalculator struct {
	Baseline   int
	Deductions map[RiskType]int
	Bonuses    map[SecurityEnhancement]int
}

// NewRiskCalculator creates a new RiskCalculator with default values
func NewRiskCalculator() *RiskCalculator {
	return &RiskCalculator{
		Baseline: 70,
		Deductions: map[RiskType]int{
			RiskTypePrivilegeEscalation:  -25,
			RiskTypeDestructiveFiles:     -30,
			RiskTypeNetworkAccess:        -10,
			RiskTypeMissingErrorHandling: -10,
			RiskTypeMissingAuditLog:      -10,
		},
		Bonuses: map[SecurityEnhancement]int{
			SecurityEnhancementEncryption:   10,
			SecurityEnhancementAuditLogging: 10,
			SecurityEnhancementHardening:    10,
		},
	}
}

// CalculateScore computes the final security score
func (rc *RiskCalculator) CalculateScore(risks []RiskAssessment, bonuses []SecurityBonus) int {
	score := rc.Baseline

	for _, risk := range risks {
		score += risk.Severity
	}

	for _, bonus := range bonuses {
		score += bonus.Points
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// DetermineDecision determines the security decision based on score
func (rc *RiskCalculator) DetermineDecision(score int) SecurityDecision {
	if score >= 80 {
		return SecurityDecisionApprove
	} else if score >= 60 {
		return SecurityDecisionConditionalApprove
	}
	return SecurityDecisionReject
}

// ProposalAnalyzer parses proposal content and identifies security patterns
type ProposalAnalyzer struct {
	privilegePatterns   []*regexp.Regexp
	destructivePatterns []*regexp.Regexp
	networkPattern      *regexp.Regexp
	validationPattern   *regexp.Regexp
	encryptionPattern   *regexp.Regexp
	auditLogPattern     *regexp.Regexp
}

// NewProposalAnalyzer creates a new ProposalAnalyzer with compiled regex patterns
func NewProposalAnalyzer() *ProposalAnalyzer {
	return &ProposalAnalyzer{
		privilegePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\bsudo\b`),
			regexp.MustCompile(`(?i)\broot\b`),
			regexp.MustCompile(`(?i)chmod\s+777`),
			regexp.MustCompile(`(?i)chown.*root`),
			regexp.MustCompile(`(?i)/etc/shadow`),
			regexp.MustCompile(`(?i)/etc/passwd`),
		},
		destructivePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)rm\s+-rf`),
			regexp.MustCompile(`(?i)\brm\s+.*\*`),
			regexp.MustCompile(`(?i)\bdd\b.*of=`),
			regexp.MustCompile(`(?i)\bmkfs\b`),
			regexp.MustCompile(`(?i)>/dev/sd`),
			regexp.MustCompile(`(?i)format.*disk`),
		},
		networkPattern:    regexp.MustCompile(`(?i)\b(curl|wget|fetch)\b`),
		validationPattern: regexp.MustCompile(`(?i)(validate|sanitize|whitelist).*url`),
		encryptionPattern: regexp.MustCompile(`(?i)\b(encrypt|cipher|AES|RSA|TLS|SSL)\b`),
		auditLogPattern:   regexp.MustCompile(`(?i)\b(audit.*log|compliance.*log|security.*event)\b`),
	}
}

// Analyze analyzes a proposal and returns risks and bonuses
func (pa *ProposalAnalyzer) Analyze(proposalText string) ([]RiskAssessment, []SecurityBonus) {
	risks := []RiskAssessment{}
	bonuses := []SecurityBonus{}

	if pa.hasPrivilegeEscalation(proposalText) {
		risks = append(risks, RiskAssessment{
			RiskType:    RiskTypePrivilegeEscalation,
			Severity:    -25,
			Description: "Proposal includes sudo/root operations without validation",
			Mitigation:  "Add user confirmation or restrict to non-privileged operations",
		})
	}

	if pa.hasDestructiveFileOps(proposalText) {
		risks = append(risks, RiskAssessment{
			RiskType:    RiskTypeDestructiveFiles,
			Severity:    -30,
			Description: "Proposal includes rm -rf or similar destructive operations",
			Mitigation:  "Use trash/safe-rm or add confirmation prompts",
		})
	}

	if pa.hasUncontrolledNetwork(proposalText) {
		risks = append(risks, RiskAssessment{
			RiskType:    RiskTypeNetworkAccess,
			Severity:    -10,
			Description: "Proposal makes network calls without proper validation",
			Mitigation:  "Add URL validation and timeout limits",
		})
	}

	if pa.addsEncryption(proposalText) {
		bonuses = append(bonuses, SecurityBonus{
			EnhancementType: SecurityEnhancementEncryption,
			Points:          10,
			Description:     "Proposal adds encryption for sensitive data",
		})
	}

	if pa.addsAuditLogging(proposalText) {
		bonuses = append(bonuses, SecurityBonus{
			EnhancementType: SecurityEnhancementAuditLogging,
			Points:          10,
			Description:     "Proposal adds audit logging for compliance",
		})
	}

	return risks, bonuses
}

func (pa *ProposalAnalyzer) hasPrivilegeEscalation(text string) bool {
	for _, pattern := range pa.privilegePatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func (pa *ProposalAnalyzer) hasDestructiveFileOps(text string) bool {
	for _, pattern := range pa.destructivePatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func (pa *ProposalAnalyzer) hasUncontrolledNetwork(text string) bool {
	hasNetwork := pa.networkPattern.MatchString(text)
	hasValidation := pa.validationPattern.MatchString(text)
	return hasNetwork && !hasValidation
}

func (pa *ProposalAnalyzer) addsEncryption(text string) bool {
	return pa.encryptionPattern.MatchString(text)
}

func (pa *ProposalAnalyzer) addsAuditLogging(text string) bool {
	return pa.auditLogPattern.MatchString(text)
}

// SecurityAgent is the main security agent with risk-based scoring
type SecurityAgent struct {
	analyzer   *ProposalAnalyzer
	calculator *RiskCalculator
}

// NewSecurityAgent creates a new SecurityAgent
func NewSecurityAgent() *SecurityAgent {
	return &SecurityAgent{
		analyzer:   NewProposalAnalyzer(),
		calculator: NewRiskCalculator(),
	}
}

// Review performs a security review of a proposal
func (sa *SecurityAgent) Review(proposalID, proposalText string) *SecurityReview {
	risks, bonuses := sa.analyzer.Analyze(proposalText)
	finalScore := sa.calculator.CalculateScore(risks, bonuses)
	decision := sa.calculator.DetermineDecision(finalScore)
	feedback := sa.generateFeedback(risks, bonuses, finalScore, decision)

	return &SecurityReview{
		ProposalID: proposalID,
		BaseScore:  sa.calculator.Baseline,
		Risks:      risks,
		Bonuses:    bonuses,
		FinalScore: finalScore,
		Decision:   decision,
		Feedback:   feedback,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

func (sa *SecurityAgent) generateFeedback(risks []RiskAssessment, bonuses []SecurityBonus, score int, decision SecurityDecision) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Security Score: %d/100", score))
	parts = append(parts, fmt.Sprintf("Decision: %s", decision))
	parts = append(parts, "")
	parts = append(parts, "## ClawdLinux Security Context")
	parts = append(parts, "Baseline assumes existing infrastructure security:")
	parts = append(parts, "- Firecracker MicroVM isolation")
	parts = append(parts, "- Workspace filesystem sandboxing")
	parts = append(parts, "- Network egress firewall (RFC1918 blocked)")
	parts = append(parts, "- Process resource limits (AVE)")
	parts = append(parts, "")

	if len(risks) > 0 {
		parts = append(parts, "## Risks Identified")
		for _, risk := range risks {
			parts = append(parts, fmt.Sprintf("- **%s** (%d pts): %s", risk.RiskType, risk.Severity, risk.Description))
			if risk.Mitigation != "" {
				parts = append(parts, fmt.Sprintf("  → Mitigation: %s", risk.Mitigation))
			}
		}
		parts = append(parts, "")
	}

	if len(bonuses) > 0 {
		parts = append(parts, "## Security Enhancements")
		for _, bonus := range bonuses {
			parts = append(parts, fmt.Sprintf("- **%s** (+%d pts): %s", bonus.EnhancementType, bonus.Points, bonus.Description))
		}
		parts = append(parts, "")
	}

	if decision == SecurityDecisionConditionalApprove {
		parts = append(parts, "## Conditions for Approval")
		parts = append(parts, "- Address identified risks before execution")
		parts = append(parts, "- Add mitigation controls as suggested")
	} else if decision == SecurityDecisionReject {
		parts = append(parts, "## Rejection Rationale")
		parts = append(parts, "Critical security risks must be addressed before resubmission")
	}

	return strings.Join(parts, "\n")
}
```