#!/usr/bin/env python3
import subprocess
import sys
import re
import os

def run_command(cmd):
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    return result.stdout.strip(), result.returncode

def get_coverage():
    print("Running tests and generating coverage profile...")
    # We ignore the exit code because 'go: no such tool "covdata"' causes exit 1 
    # even if coverage.out is generated successfully.
    run_command("go test ./... -coverprofile=coverage.out")

    if not os.path.exists("coverage.out") or os.path.getsize("coverage.out") == 0:
        print("❌ Failed to generate coverage.out!")
        sys.exit(1)

    output, _ = run_command("go tool cover -func=coverage.out")
    lines = output.split('\n')
    total_line = [l for l in lines if 'total:' in l]
    if not total_line:
        print("❌ Could not find total coverage in output!")
        sys.exit(1)

    match = re.search(r'(\d+\.\d+)%', total_line[0])
    return float(match.group(1)) if match else 0.0

def get_changed_lines():
    # Fallback chain: origin/main -> main -> HEAD~1
    base = None
    for candidate in ["origin/main", "main", "HEAD~1"]:
        _, code = run_command(f"git rev-parse --verify {candidate}")
        if code == 0:
            base = candidate
            break

    if not base:
        print("⚠️ No base commit found for diff. Checking all local changes.")
        diff_cmd = "git diff -U0"
    else:
        print(f"ℹ️ Comparing against base: {base}")
        diff_cmd = f"git diff -U0 {base}...HEAD"

    diff_output, _ = run_command(diff_cmd)
    changed = {}
    current_file = None

    for line in diff_output.split('\n'):
        if line.startswith('+++ b/'):
            current_file = line[6:]
            changed[current_file] = set()
        elif line.startswith('@@') and current_file:
            match = re.search(r'\+(\d+)(?:,(\d+))?', line)
            if match:
                start = int(match.group(1))
                count = int(match.group(2)) if match.group(2) else 1
                if count == 0: continue
                for i in range(start, start + count):
                    changed[current_file].add(i)
    return changed

def calculate_incremental_coverage(changed_lines):
    if not changed_lines:
        print("ℹ️ No changed lines detected.")
        return 100.0

    covered_changed = 0
    total_changed = 0

    with open("coverage.out", "r") as f:
        next(f) # skip mode line
        for line in f:
            match = re.match(r'([^:]+):(\d+)\.\d+,(\d+)\.\d+ (\d+) (\d+)', line)
            if match:
                full_path, start, end, num_stmt, count = match.groups()
                rel_path = full_path.split('hyperagent/')[-1]

                if rel_path in changed_lines:
                    start, end, num_stmt, count = int(start), int(end), int(num_stmt), int(count)
                    for l in range(start, end + 1):
                        if l in changed_lines[rel_path]:
                            total_changed += 1
                            if count > 0:
                                covered_changed += 1

    if total_changed == 0: 
        print("ℹ️ No testable changed lines detected.")
        return 100.0
    return (covered_changed / total_changed) * 100

def main():
    abs_cov = get_coverage()
    changed = get_changed_lines()
    inc_cov = calculate_incremental_coverage(changed)

    print(f"\nAbsolute Coverage: {abs_cov:.2f}%")
    print(f"Incremental Coverage: {inc_cov:.2f}%")

    failed = False
    if abs_cov < 90.0:
        print(f"❌ Absolute coverage {abs_cov:.2f}% is below 90%!")
        failed = True
    if inc_cov < 95.0:
        print(f"❌ Incremental coverage {inc_cov:.2f}% is below 95%!")
        failed = True

    if failed:
        sys.exit(1)

    print("✅ Coverage health checks passed!")

if __name__ == '__main__':
    main()
