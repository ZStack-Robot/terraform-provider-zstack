#!/usr/bin/env bash
# Batch test runner for terraform-provider-zstack
# Usage:
#   ./scripts/run_tests.sh                    # Run all tests
#   ./scripts/run_tests.sh datasource         # Run all data source tests
#   ./scripts/run_tests.sh resource            # Run all resource tests
#   ./scripts/run_tests.sh -run TestAccZStack  # Run tests matching pattern
#   ./scripts/run_tests.sh -f zone,images      # Run tests for specific resources
#   ./scripts/run_tests.sh -list               # List all available tests
#   ./scripts/run_tests.sh -gen-only           # Only regenerate testdata/env.json
#   ./scripts/run_tests.sh -skip-gen           # Skip regenerating testdata/env.json

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_PKG="./zstack/provider/"
ENV_FILE="$PROJECT_ROOT/.env.test"
REPORT_DIR="$PROJECT_ROOT/test-reports"
ENV_JSON="$PROJECT_ROOT/zstack/provider/testdata/env.json"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS] [CATEGORY]

Categories:
  datasource    Run all data source tests
  resource      Run all resource tests
  all           Run all tests (default)

Options:
  -run PATTERN   Run tests matching Go test pattern (e.g. -run TestAccZStackImage)
  -f RESOURCES   Comma-separated list of resources to test (e.g. -f zone,images,clusters)
  -list          List all available test functions
  -v             Verbose output (show all test log)
  -timeout DUR   Test timeout (default: 30m)
  -parallel N    Number of parallel tests (default: 1, sequential)
  -report        Generate HTML test report
  -dry-run       Show which tests would run without executing
  -gen-only      Only regenerate testdata/env.json, don't run tests
  -skip-gen      Skip regenerating testdata/env.json (use cached data)
  -h, --help     Show this help message

Examples:
  $(basename "$0")                          # Regenerate env data + run all tests
  $(basename "$0") datasource               # Run only data source tests
  $(basename "$0") resource                 # Run only resource tests
  $(basename "$0") -f zone,images           # Run zone and images tests
  $(basename "$0") -run TestAccCreateImage  # Run specific test by name
  $(basename "$0") -report                  # Run all tests with HTML report
  $(basename "$0") -gen-only                # Only regenerate testdata/env.json
  $(basename "$0") -skip-gen                # Skip env data generation, run tests
EOF
}

# Load environment
load_env() {
    if [[ -f "$ENV_FILE" ]]; then
        echo -e "${CYAN}Loading environment from $ENV_FILE${NC}"
        # shellcheck disable=SC1090
        source "$ENV_FILE"
    else
        echo -e "${RED}Error: $ENV_FILE not found.${NC}"
        echo "Copy .env.test.example to .env.test and configure your ZStack credentials."
        exit 1
    fi

    # Validate required vars
    if [[ -z "${ZSTACK_HOST:-}" ]]; then
        echo -e "${RED}Error: ZSTACK_HOST is not set${NC}"
        exit 1
    fi

    if [[ -z "${TF_ACC:-}" ]]; then
        export TF_ACC=1
    fi

    # Check auth method
    if [[ -n "${ZSTACK_ACCESS_KEY_ID:-}" && -n "${ZSTACK_ACCESS_KEY_SECRET:-}" ]]; then
        echo -e "${GREEN}Auth: AccessKey (${ZSTACK_ACCESS_KEY_ID:0:8}...)${NC}"
    elif [[ -n "${ZSTACK_ACCOUNT_NAME:-}" && -n "${ZSTACK_ACCOUNT_PASSWORD:-}" ]]; then
        echo -e "${GREEN}Auth: Account ($ZSTACK_ACCOUNT_NAME)${NC}"
    else
        echo -e "${RED}Error: No authentication configured. Set AK/SK or account/password in .env.test${NC}"
        exit 1
    fi

    echo -e "${GREEN}Target: ${ZSTACK_HOST}:${ZSTACK_PORT:-8080}${NC}"
    echo ""
}

# Generate testdata/env.json by querying the real ZStack environment
generate_env_data() {
    echo -e "${CYAN}Generating testdata/env.json from ZStack environment...${NC}"

    cd "$PROJECT_ROOT"
    if go run ./zstack/provider/testdata/generate_env.go; then
        echo -e "${GREEN}testdata/env.json generated successfully${NC}"
        echo ""
    else
        echo -e "${RED}Failed to generate testdata/env.json${NC}"
        exit 1
    fi
}

# List all test functions
list_tests() {
    echo -e "${CYAN}Available test functions:${NC}"
    echo ""
    echo -e "${YELLOW}=== Data Source Tests ===${NC}"
    grep -rh "^func Test" "$PROJECT_ROOT/zstack/provider/"data_source_*_test.go 2>/dev/null | sed 's/func \(Test[^(]*\).*/  \1/' | sort
    echo ""
    echo -e "${YELLOW}=== Resource Tests ===${NC}"
    grep -rh "^func Test" "$PROJECT_ROOT/zstack/provider/"resource_*_test.go 2>/dev/null | sed 's/func \(Test[^(]*\).*/  \1/' | sort
    echo ""
    total=$(grep -rch "^func Test" "$PROJECT_ROOT/zstack/provider/"*_test.go 2>/dev/null | paste -sd+ - | bc)
    echo -e "${CYAN}Total: ${total} tests${NC}"
}

# Build test run pattern from resource filter
build_filter_pattern() {
    local resources="$1"
    local patterns=()

    IFS=',' read -ra ITEMS <<< "$resources"
    for item in "${ITEMS[@]}"; do
        item=$(echo "$item" | xargs) # trim whitespace
        # Match test function names containing the resource name (case insensitive-ish)
        patterns+=("TestAcc.*${item}")
    done

    # Join with | for Go test -run regex
    local joined
    joined=$(IFS='|'; echo "${patterns[*]}")
    echo "$joined"
}

# Run tests
run_tests() {
    local run_pattern="${1:-}"
    local verbose="${2:-false}"
    local timeout="${3:-30m}"
    local parallel="${4:-1}"
    local generate_report="${5:-false}"
    local dry_run="${6:-false}"

    local go_test_args=(
        "go" "test"
        "$TEST_PKG"
        "-timeout" "$timeout"
        "-parallel" "$parallel"
        "-count=1"
    )

    if [[ -n "$run_pattern" ]]; then
        go_test_args+=("-run" "$run_pattern")
    fi

    if [[ "$verbose" == "true" ]]; then
        go_test_args+=("-v")
    else
        go_test_args+=("-v")  # Always use -v for progress visibility
    fi

    if [[ "$dry_run" == "true" ]]; then
        echo -e "${YELLOW}Dry run - would execute:${NC}"
        echo "  cd $PROJECT_ROOT && ${go_test_args[*]}"
        if [[ -n "$run_pattern" ]]; then
            echo ""
            echo -e "${CYAN}Matching tests:${NC}"
            grep -rh "^func Test" "$PROJECT_ROOT/zstack/provider/"*_test.go 2>/dev/null \
                | sed 's/func \(Test[^(]*\).*/\1/' \
                | grep -iE "$run_pattern" \
                | sort \
                | sed 's/^/  /'
        fi
        return 0
    fi

    echo -e "${CYAN}Running tests...${NC}"
    echo -e "${YELLOW}Command: ${go_test_args[*]}${NC}"
    echo ""

    cd "$PROJECT_ROOT"

    local start_time
    start_time=$(date +%s)

    if [[ "$generate_report" == "true" ]]; then
        mkdir -p "$REPORT_DIR"
        local report_file="$REPORT_DIR/test-report-$(date +%Y%m%d-%H%M%S).log"
        echo -e "${CYAN}Test output will be saved to: $report_file${NC}"

        "${go_test_args[@]}" -json 2>&1 | tee "$report_file"
        local exit_code=${PIPESTATUS[0]}
    else
        "${go_test_args[@]}" 2>&1
        local exit_code=$?
    fi

    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))

    echo ""
    echo "─────────────────────────────────────────"
    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}ALL TESTS PASSED${NC} (${duration}s)"
    else
        echo -e "${RED}SOME TESTS FAILED${NC} (${duration}s, exit code: $exit_code)"
    fi
    echo "─────────────────────────────────────────"

    return $exit_code
}

# Main
main() {
    local run_pattern=""
    local category=""
    local verbose="false"
    local timeout="30m"
    local parallel="1"
    local generate_report="false"
    local dry_run="false"
    local gen_only="false"
    local skip_gen="false"

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -run)
                run_pattern="$2"
                shift 2
                ;;
            -f)
                run_pattern=$(build_filter_pattern "$2")
                shift 2
                ;;
            -list)
                list_tests
                exit 0
                ;;
            -v)
                verbose="true"
                shift
                ;;
            -timeout)
                timeout="$2"
                shift 2
                ;;
            -parallel)
                parallel="$2"
                shift 2
                ;;
            -report)
                generate_report="true"
                shift
                ;;
            -dry-run)
                dry_run="true"
                shift
                ;;
            -gen-only)
                gen_only="true"
                shift
                ;;
            -skip-gen)
                skip_gen="true"
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            datasource|datasources)
                category="datasource"
                shift
                ;;
            resource|resources)
                category="resource"
                shift
                ;;
            all)
                category=""
                shift
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                usage
                exit 1
                ;;
        esac
    done

    # Build pattern from category if no explicit -run
    if [[ -z "$run_pattern" && -n "$category" ]]; then
        case "$category" in
            datasource)
                run_pattern="TestAcc.*DataSource"
                ;;
            resource)
                run_pattern="TestAccCreate|TestAccUpdate|TestAccDelete|TestAccImport|TestAcc.*Resource"
                ;;
        esac
    fi

    load_env

    # Generate testdata/env.json unless skipped
    if [[ "$skip_gen" != "true" ]]; then
        generate_env_data
    else
        if [[ ! -f "$ENV_JSON" ]]; then
            echo -e "${YELLOW}Warning: $ENV_JSON not found. Run without -skip-gen first.${NC}"
            exit 1
        fi
        echo -e "${YELLOW}Skipping env data generation (using cached $ENV_JSON)${NC}"
        echo ""
    fi

    # If gen-only, stop here
    if [[ "$gen_only" == "true" ]]; then
        echo -e "${GREEN}Done. Only generated testdata/env.json.${NC}"
        exit 0
    fi

    run_tests "$run_pattern" "$verbose" "$timeout" "$parallel" "$generate_report" "$dry_run"
}

main "$@"
