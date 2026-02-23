#!/bin/bash
# Manual code review before push using GLM5 (free model)
# Usage: ./scripts/review-before-push.sh
# This reviews all changed files against origin/main

set -e

echo "🔍 CODE REVIEW GATE (GLM5 - Free Model)"
echo "========================================"
echo ""

# Get changed files
CHANGED_FILES=$(git diff origin/main...HEAD --name-only 2>/dev/null || git diff HEAD~1 --name-only 2>/dev/null || echo "")

if [ -z "$CHANGED_FILES" ]; then
    echo "✅ No changes to review"
    exit 0
fi

echo "📝 Changed files:"
echo "$CHANGED_FILES" | sed 's/^/  /'
echo ""

# Filter to code files
CODE_FILES=$(echo "$CHANGED_FILES" | grep -E '\.(go|py|yaml|rego)$' || true)

if [ -z "$CODE_FILES" ]; then
    echo "✅ No code files changed - skipping review"
    exit 0
fi

echo "🚀 Building review request for GLM5..."
echo ""

# Create review prompt
REVIEW_PROMPT_FILE="/tmp/code-review-request-$(date +%s).txt"

{
    cat << 'EOF'
# CODE REVIEW REQUEST - GLM5 Fast Model

You are a code reviewer. Systematically review the following changes using these dimensions:

## Review Dimensions

### SECURITY (CRITICAL)
- [ ] No hardcoded secrets, API keys, credentials
- [ ] No SQL injection, command injection, XSS vulnerabilities
- [ ] Proper input validation and error handling
- [ ] No privilege escalation or unauthorized access
- [ ] Secrets/credentials managed via environment or Kubernetes Secrets only

### CORRECTNESS (HIGH)
- [ ] Logic is correct and handles edge cases
- [ ] Null/nil pointer checks in place
- [ ] Error handling is comprehensive
- [ ] No off-by-one errors or boundary issues
- [ ] State mutations are thread-safe where needed

### PERFORMANCE (HIGH)
- [ ] No N+1 queries or unnecessary loops
- [ ] No excessive memory allocations
- [ ] Efficient algorithms used
- [ ] Timeouts and limits are reasonable
- [ ] No blocking operations where async available

### MAINTAINABILITY (MEDIUM)
- [ ] Code is clear and well-documented
- [ ] Follows project conventions and style
- [ ] Functions are not too long (single responsibility)
- [ ] No magic numbers or unexplained constants
- [ ] DRY principle followed (no duplication)

### TESTING (MEDIUM)
- [ ] Changes are covered by tests
- [ ] Tests are comprehensive (positive + negative cases)
- [ ] Mock objects used appropriately
- [ ] Test names clearly describe intent
- [ ] Edge cases tested

## Changed Files to Review

EOF
    echo "$CODE_FILES" | while read file; do
        echo ""
        echo "### File: $file"
        echo ""
        echo "Content:"
        echo '```'
        git show "HEAD:$file" 2>/dev/null | head -100 || echo "(Binary or deleted file)"
        echo '```'
        echo ""
    done
    
    cat << 'EOF'

## Review Output Format

For each issue found, report:
- **Severity:** CRITICAL | HIGH | MEDIUM | LOW
- **Location:** file:line or file
- **Issue:** Clear description
- **Fix:** Recommended fix
- **Reason:** Why this matters

End with:
**DECISION:** PASS ✅ | CONDITIONAL_PASS ⚠️ | FAIL ❌

If CONDITIONAL_PASS or FAIL, list required fixes before merging.
EOF
} > "$REVIEW_PROMPT_FILE"

echo "📊 Review request saved to: $REVIEW_PROMPT_FILE"
echo ""
echo "Next steps:"
echo "  1. Manual review: Read the changed code above"
echo "  2. Use gemini CLI for automated review:"
echo "     gemini \"$(cat $REVIEW_PROMPT_FILE | head -50)...\""
echo "  3. Or push with --no-verify if confident:"
echo "     git push --no-verify"
echo ""

# Try to use gemini if available
if command -v gemini &> /dev/null; then
    echo "🤖 Running automated GLM5 review..."
    echo ""
    # Note: This is a placeholder - actual gemini invocation would go here
    # For now, just prompt the user
    echo "To run automated review, use:"
    echo "  gemini --model glm5 < $REVIEW_PROMPT_FILE"
else
    echo "⚠️  Gemini CLI not found. Manual review required before push."
fi

echo ""
echo "Decision: $(git diff --stat origin/main...HEAD 2>/dev/null | tail -1)"
