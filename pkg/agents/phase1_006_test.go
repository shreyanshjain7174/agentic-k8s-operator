// Auto-generated tests for phase1-006

```go
package security

import (
	"strings"
	"testing"
	"time"
)

func TestNewRiskCalculator(t *testing.T) {
	rc := NewRiskCalculator()

	if rc.Baseline != 70 {
		t.Errorf("expected baseline 70, got %d", rc.Baseline)
	}

	if rc.Deductions[RiskTypePrivilegeEscalation] != -25 {
		t.Errorf("expected privilege escalation deduction -25, got %d", rc.Deductions[RiskTypePrivilegeEscalation])
	}

	if rc.Bonuses[SecurityEnhancementEncryption] != 10 {
		t.Errorf("expected encryption bonus 10, got %d", rc.Bonuses[SecurityEnhancementEncryption])
	}
}

func TestRiskCalculator_CalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		risks    []RiskAssessment
		bonuses  []SecurityBonus
		expected int
	}{
		{
			name:     "baseline with no risks or bonuses",
			risks:    []RiskAssessment{},
			bonuses:  []SecurityBonus{},
			expected: 70,
		},
		{
			name: "single risk",
			risks: []RiskAssessment{
				{RiskType: RiskTypePrivilegeEscalation, Severity: -25},
			},
			bonuses:  []SecurityBonus{},
			expected: 45,
		},
		{
			name: "multiple risks",
			risks: []RiskAssessment{
				{RiskType: RiskTypePrivilegeEscalation, Severity: -25},
				{RiskType: RiskTypeDestructiveFiles, Severity: -30},
			},
			bonuses:  []SecurityBonus{},
			expected: 15,
		},
		{
			name:  "single bonus",
			risks: []RiskAssessment{},
			bonuses: []SecurityBonus{
				{EnhancementType: SecurityEnhancementEncryption, Points: 10},
			},
			expected: 80,
		},
		{
			name: "risks and bonuses combined",
			risks: []RiskAssessment{
				{RiskType: RiskTypePrivilegeEscalation, Severity: -25},
			},
			bonuses: []SecurityBonus{
				{EnhancementType: SecurityEnhancementEncryption, Points: 10},
			},
			expected: 55,
		},
		{
			name: "score capped at 0",
			risks: []RiskAssessment{
				{RiskType: RiskTypePrivilegeEscalation, Severity: -25},
				{RiskType: RiskTypeDestructiveFiles, Severity: -30},
				{RiskType: RiskTypeNetworkAccess, Severity: -10},
				{RiskType: RiskTypeMissingErrorHandling, Severity: -10},
			},
			bonuses:  []SecurityBonus{},
			expected: 0,
		},
		{
			name:  "score capped at 100",
			risks: []RiskAssessment{},
			bonuses: []SecurityBonus{
				{EnhancementType: SecurityEnhancementEncryption, Points: 10},
				{EnhancementType: SecurityEnhancementAuditLogging, Points: 10},
				{EnhancementType: SecurityEnhancementHardening, Points: 10},
				{EnhancementType: SecurityEnhancementEncryption, Points: 10},
			},
			expected: 100,
		},
	}

	rc := NewRiskCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rc.CalculateScore(tt.risks, tt.bonuses)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestRiskCalculator_DetermineDecision(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected SecurityDecision
	}{
		{"approve at 100", 100, SecurityDecisionApprove},
		{"approve at 80", 80, SecurityDecisionApprove},
		{"conditional at 79", 79, SecurityDecisionConditionalApprove},
		{"conditional at 60", 60, SecurityDecisionConditionalApprove},
		{"reject at 59", 59, SecurityDecisionReject},
		{"reject at 0", 0, SecurityDecisionReject},
	}

	rc := NewRiskCalculator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rc.DetermineDecision(tt.score)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestProposalAnalyzer_hasPrivilegeEscalation(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"sudo command", "run sudo apt-get install", true},
		{"root user", "switch to root user", true},
		{"chmod 777", "chmod 777 file.txt", true},
		{"chown root", "chown root:root /etc/config", true},
		{"etc shadow", "cat /etc/shadow", true},
		{"etc passwd", "cat /etc/passwd", true},
		{"case insensitive sudo", "SUDO make install", true},
		{"no privilege escalation", "run normal command", false},
		{"empty text", "", false},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pa.hasPrivilegeEscalation(tt.text)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProposalAnalyzer_hasDestructiveFileOps(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"rm -rf command", "rm -rf /var/data", true},
		{"rm with wildcard", "rm *.log", true},
		{"dd command", "dd if=/dev/zero of=/dev/sda", true},
		{"mkfs command", "mkfs.ext4 /dev/sdb1", true},
		{"write to block device", "echo test >/dev/sda", true},
		{"format disk", "format disk /dev/sdc", true},
		{"case insensitive", "RM -RF /tmp", true},
		{"safe rm command", "rm specific_file.txt", false},
		{"no destructive ops", "ls -la", false},
		{"empty text", "", false},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pa.hasDestructiveFileOps(tt.text)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProposalAnalyzer_hasUncontrolledNetwork(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"curl without validation", "curl http://example.com", true},
		{"wget without validation", "wget http://example.com/file", true},
		{"fetch without validation", "fetch data from API", true},
		{"curl with validation", "validate url before curl http://example.com", false},
		{"wget with sanitization", "sanitize url then wget http://example.com", false},
		{"curl with whitelist", "whitelist url check curl http://example.com", false},
		{"no network access", "read local file", false},
		{"empty text", "", false},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pa.hasUncontrolledNetwork(tt.text)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for text: %s", tt.expected, result, tt.text)
			}
		})
	}
}

func TestProposalAnalyzer_addsEncryption(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"encrypt keyword", "encrypt sensitive data", true},
		{"cipher keyword", "use cipher for security", true},
		{"AES encryption", "implement AES-256 encryption", true},
		{"RSA encryption", "use RSA public key", true},
		{"TLS encryption", "enable TLS for connections", true},
		{"SSL encryption", "configure SSL certificates", true},
		{"case insensitive", "ENCRYPT the payload", true},
		{"no encryption", "store data in plain text", false},
		{"empty text", "", false},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pa.addsEncryption(tt.text)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProposalAnalyzer_addsAuditLogging(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"audit log", "add audit log for actions", true},
		{"compliance log", "implement compliance logging", true},
		{"security event", "log security events", true},
		{"case insensitive", "AUDIT LOGGING enabled", true},
		{"no audit logging", "regular logging only", false},
		{"empty text", "", false},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pa.addsAuditLogging(tt.text)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProposalAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name          string
		proposalText  string
		expectedRisks int
		expectedBonus int
	}{
		{
			name:          "clean proposal",
			proposalText:  "implement new feature with proper error handling",
			expectedRisks: 0,
			expectedBonus: 0,
		},
		{
			name:          "proposal with privilege escalation",
			proposalText:  "run sudo apt-get install package",
			expectedRisks: 1,
			expectedBonus: 0,
		},
		{
			name:          "proposal with multiple risks",
			proposalText:  "sudo rm -rf /var/data && curl http://untrusted.com",
			expectedRisks: 3,
			expectedBonus: 0,
		},
		{
			name:          "proposal with encryption bonus",
			proposalText:  "implement AES-256 encryption for user data",
			expectedRisks: 0,
			expectedBonus: 1,
		},
		{
			name:          "proposal with audit logging bonus",
			proposalText:  "add audit log for all admin actions",
			expectedRisks: 0,
			expectedBonus: 1,
		},
		{
			name:          "proposal with risks and bonuses",
			proposalText:  "sudo install package, enable TLS encryption, and add compliance logging",
			expectedRisks: 1,
			expectedBonus: 2,
		},
		{
			name:          "empty proposal",
			proposalText:  "",
			expectedRisks: 0,
			expectedBonus: 0,
		},
	}

	pa := NewProposalAnalyzer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risks, bonuses := pa.Analyze(tt.proposalText)
			if len(risks) != tt.expectedRisks {
				t.Errorf("expected %d risks, got %d", tt.expectedRisks, len(risks))
			}
			if len(bonuses) != tt.expectedBonus {
				t.Errorf("expected %d bonuses, got %d", tt.expectedBonus, len(bonuses))
			}
		})
	}
}

func TestNewSecurityAgent(t *testing.T) {
	sa := NewSecurityAgent()

	if sa.analyzer == nil {
		t.Error("expected analyzer to be initialized")
	}
	if sa.calculator == nil {
		t.Error("expected calculator to be initialized")
	}
}

func TestSecurityAgent_Review(t *testing.T) {
	tests := []struct {
		name             string
		proposalID       string
		proposalText     string
		expectedDecision SecurityDecision
		minScore         int
		maxScore         int
	}{
		{
			name:             "clean proposal - approve",
			proposalID:       "PROP-001",
			proposalText:     "implement new feature with tests",
			expectedDecision: SecurityDecisionConditionalApprove,
			minScore:         70,
			maxScore:         70,
		},
		{
			name:             "proposal with encryption - approve",
			proposalID:       "PROP-002",
			proposalText:     "add AES encryption and audit logging",
			expectedDecision: SecurityDecisionApprove,
			minScore:         80,
			maxScore:         90,
		},
		{
			name:             "proposal with single risk - conditional",
			proposalID:       "PROP-003",
			proposalText:     "install package using sudo",
			expectedDecision: SecurityDecisionConditionalApprove,
			minScore:         45,
			maxScore:         45,
		},
		{
			name:             "proposal with multiple risks - reject",
			proposalID:       "PROP-004",
			proposalText:     "sudo rm -rf /var/data && wget http://untrusted.com",
			expectedDecision: SecurityDecisionReject,
			minScore:         0,
			maxScore:         35,
		},
		{
			name:             "empty proposal",
			proposalID:       "PROP-005",
			proposalText:     "",
			expectedDecision: SecurityDecisionConditionalApprove,
			minScore:         70,
			maxScore:         70,
		},
	}

	sa := NewSecurityAgent()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review := sa.Review(tt.proposalID, tt.proposalText)

			if review.ProposalID != tt.proposalID {
				t.Errorf("expected proposal ID %s, got %s", tt.proposalID, review.ProposalID)
			}

			if review.Decision != tt.expectedDecision {
				t.Errorf("expected decision %s, got %s", tt.expectedDecision, review.Decision)
			}

			if review.FinalScore < tt.minScore || review.FinalScore > tt.maxScore {
				t.Errorf("expected score between %d-%d, got %d", tt.minScore, tt.maxScore, review.FinalScore)
			}

			if review.BaseScore != 70 {
				t.Errorf("expected base score 70, got %d", review.BaseScore)
			}

			if review.Feedback == "" {
				t.Error("expected feedback to be generated")
			}

			if review.Timestamp == "" {
				t.Error("expected timestamp to be set")
			}

			// Validate timestamp format
			_, err := time.Parse(time.RFC3339, review.Timestamp)
			if err != nil {
				t.Errorf("invalid timestamp format: %v", err)
			}

			// Validate feedback contains key sections
			if !strings.Contains(review.Feedback, "Security Score:") {
				t.Error("feedback should contain Security Score")
			}
			if !strings.Contains(review.Feedback, "ClawdLinux Security Context") {
				t.Error("feedback should contain ClawdLinux Security Context")
			}
		})
	}
}

func TestSecurityAgent_generateFeedback(t *testing.T) {
	sa := NewSecurityAgent()

	tests := []struct {
		name            string
		risks           []RiskAssessment
		bonuses         []SecurityBonus
		score           int
		decision        SecurityDecision
		expectedContent []string
	}{
		{
			name:    "feedback with risks",
			risks:   []RiskAssessment{{RiskType: RiskTypePrivilegeEscalation, Severity: -25, Description: "sudo detected"}},
			bonuses: []SecurityBonus{},
			score:   45,
			decision: SecurityDecisionConditionalApprove,
			expectedContent: []string{
				"Security Score: 45/100",
				"Decision: CONDITIONAL_APPROVE",
				"## Risks Identified",
				"privilege_escalation",
				"## Conditions for Approval",
			},
		},
		{
			name:  "feedback with bonuses",
			risks: []RiskAssessment{},
			bonuses: []SecurityBonus{
				{EnhancementType: SecurityEnhancementEncryption, Points: 10, Description: "AES added"},
			},
			score:    80,
			decision: SecurityDecisionApprove,
			expectedContent: []string{
				"Security Score: 80/100",
				"Decision: APPROVE",
				"## Security Enhancements",
				"encryption_added",
			},
		},
		{
			name: "feedback with rejection",
			risks: []RiskAssessment{
				{RiskType: RiskTypeDestructiveFiles, Severity: -30, Description: "rm -rf detected"},
			},
			bonuses:  []SecurityBonus{},
			score:    40,
			decision: SecurityDecisionReject,
			expectedContent: []string{
				"Security Score: 40/100",
				"Decision: REJECT",
				"## Rejection Rationale",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedback := sa.generateFeedback(tt.risks, tt.bonuses, tt.score, tt.decision)

			for _, content := range tt.expectedContent {
				if !strings.Contains(feedback, content) {
					t.Errorf("expected feedback to contain %q", content)
				}
			}

			// All feedback should contain baseline context
			if !strings.Contains(feedback, "Firecracker MicroVM isolation") {
				t.Error("feedback should mention Firecracker MicroVM")
			}
		})
	}
}

func TestSecurityReview_Integration(t *testing.T) {
	sa := NewSecurityAgent()

	// Test full integration with realistic proposal
	proposalText := `
		Deploy new microservice with the following:
		- Implement TLS encryption for all API endpoints
		- Add audit logging for compliance
		- Install dependencies using sudo apt-get
		- Download configuration from curl https://config.internal.com
	`

	review := sa.Review("PROP-INTEGRATION-001", proposalText)

	// Should have risks (sudo, unvalidated curl)
	if len(review.Risks) == 0 {
		t.Error("expected risks to be identified")
	}

	// Should have bonuses (TLS, audit logging)
	if len(review.Bonuses) == 0 {
		t.Error("expected security bonuses to be identified")
	}

	// Validate all fields are populated
	if review.ProposalID == "" {
		t.Error("proposal ID should be set")
	}
	if review.BaseScore == 0 {
		t.Error("base score should be set")
	}
	if review.FinalScore == 0 {
		t.Error("final score should be calculated")
	}
	if review.Decision == "" {
		t.Error("decision should be determined")
	}
	if review.Feedback == "" {
		t.Error("feedback should be generated")
	}
	if review.Timestamp == "" {
		t.Error("timestamp should be set")
	}
}
```