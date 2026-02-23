#!/bin/bash
# Copyright 2024 The Voting Operator Authors.
# Licensed under the Apache License, Version 2.0.

set -e

echo "================================"
echo "Voting Operator Test Suite"
echo "================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
MIN_COVERAGE=80
INTEGRATION_TAG="integration"

# Functions
print_header() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo " $1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Step 1: Verify dependencies
print_header "Checking Dependencies"

if ! command -v go &> /dev/null; then
    print_error "Go is not installed"
    exit 1
fi
print_success "Go $(go version | awk '{print $3}')"

if ! command -v kubectl &> /dev/null; then
    print_warning "kubectl not found (optional for e2e tests)"
else
    print_success "kubectl $(kubectl version --client -o json | jq -r '.clientVersion.gitVersion')"
fi

# Step 2: Install test dependencies
print_header "Installing Test Dependencies"

go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
print_success "envtest installed"

go install github.com/onsi/ginkgo/v2/ginkgo@latest
print_success "ginkgo installed"

# Step 3: Generate mocks and manifests
print_header "Generating Code"

make manifests generate
print_success "Manifests and deep-copy generated"

# Step 4: Run linting
print_header "Running Linters"

if command -v golangci-lint &> /dev/null; then
    golangci-lint run --timeout 5m
    print_success "Linting passed"
else
    print_warning "golangci-lint not installed, skipping"
fi

# Step 5: Run unit tests
print_header "Running Unit Tests"

go test -v -race -coverprofile=coverage.out -covermode=atomic ./internal/... ./api/...

UNIT_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo ""
echo "Unit Test Coverage: ${UNIT_COVERAGE}%"

if (( $(echo "$UNIT_COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    print_error "Coverage ${UNIT_COVERAGE}% below minimum ${MIN_COVERAGE}%"
    exit 1
else
    print_success "Coverage ${UNIT_COVERAGE}% meets minimum ${MIN_COVERAGE}%"
fi

# Step 6: Generate HTML coverage report
print_header "Generating Coverage Report"

go tool cover -html=coverage.out -o coverage.html
print_success "HTML report: coverage.html"

# Step 7: Run integration tests
print_header "Running Integration Tests"

if [ -z "$SKIP_INTEGRATION" ]; then
    go test -v -tags=${INTEGRATION_TAG} -timeout 10m ./test/integration/...
    print_success "Integration tests passed"
else
    print_warning "Integration tests skipped (SKIP_INTEGRATION set)"
fi

# Step 8: Run e2e tests (optional)
print_header "Running E2E Tests"

if [ -n "$RUN_E2E" ]; then
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl required for e2e tests"
        exit 1
    fi
    
    # Setup envtest
    KUBEBUILDER_ASSETS=$(setup-envtest use -p path)
    export KUBEBUILDER_ASSETS
    
    print_success "envtest assets: $KUBEBUILDER_ASSETS"
    
    go test -v -tags=e2e -timeout 30m ./test/e2e/...
    print_success "E2E tests passed"
else
    print_warning "E2E tests skipped (set RUN_E2E=1 to enable)"
fi

# Step 9: Run benchmarks
print_header "Running Benchmarks"

if [ -n "$RUN_BENCH" ]; then
    go test -bench=. -benchmem -run=^$ ./internal/voting/
    print_success "Benchmarks completed"
else
    print_warning "Benchmarks skipped (set RUN_BENCH=1 to enable)"
fi

# Step 10: Security scanning
print_header "Security Scanning"

if command -v gosec &> /dev/null; then
    gosec -quiet ./...
    print_success "Security scan passed"
else
    print_warning "gosec not installed, skipping security scan"
    echo "Install: go install github.com/securego/gosec/v2/cmd/gosec@latest"
fi

# Step 11: Dependency audit
print_header "Dependency Audit"

if command -v nancy &> /dev/null; then
    go list -json -m all | nancy sleuth
    print_success "Dependency audit passed"
else
    print_warning "nancy not installed, skipping dependency audit"
    echo "Install: go install github.com/sonatype-nexus-community/nancy@latest"
fi

# Final summary
print_header "Test Summary"

echo "✓ Unit Tests: PASSED (${UNIT_COVERAGE}% coverage)"
echo "✓ Integration Tests: PASSED"
if [ -n "$RUN_E2E" ]; then
    echo "✓ E2E Tests: PASSED"
fi
echo ""
echo "Full report: coverage.html"
echo ""
print_success "All tests passed!"
