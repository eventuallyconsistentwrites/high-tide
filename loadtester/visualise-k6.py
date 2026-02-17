import re
import sys
import argparse
import matplotlib.pyplot as plt
import numpy as np

def parse_time_to_ms(time_str):
    """Converts k6 time strings (e.g., 1.83s, 500ms, 59µs) to milliseconds."""
    if not time_str:
        return 0.0

    # Remove 's' suffix if just s
    val = 0.0
    unit = ""

    if 'ms' in time_str:
        val = float(time_str.replace('ms', ''))
        unit = 'ms'
    elif 'µs' in time_str:
        val = float(time_str.replace('µs', '')) / 1000.0
        unit = 'µs'
    elif 's' in time_str:
        # Handle minutes like 1m2s? k6 usually does 1m30s
        if 'm' in time_str:
            parts = time_str.split('m')
            minutes = float(parts[0])
            seconds = float(parts[1].replace('s', ''))
            val = (minutes * 60 + seconds) * 1000
        else:
            val = float(time_str.replace('s', '')) * 1000.0
            unit = 's'
    else:
        val = float(time_str) # Assume ms or raw?

    return val

def parse_log_file(filepath):
    """
    Parses a docker-compose log file containing multiple k6 workers.
    Returns a dictionary of worker_id -> metrics.
    """

    # Regex patterns
    # Matches: k6-worker-12 |     http_reqs......................: 16213  245.471587/s
    rps_pattern = re.compile(r'k6-worker-(\d+)\s+\|\s+http_reqs\.+: (\d+)\s+([\d\.]+)/s')

    # Matches: k6-worker-12 |     http_req_duration..............: avg=1.83s    min=59.58µs med=883.75ms max=11.79s p(90)=5.21s    p(95)=6.21s
    latency_pattern = re.compile(r'k6-worker-(\d+)\s+\|\s+http_req_duration\.+: avg=([\w\.]+) .*? med=([\w\.]+) .*? p\(95\)=([\w\.]+)')

    # Matches: k6-worker-12 |     http_req_failed................: 99.38% 16113 out of 16213
    fail_pattern = re.compile(r'k6-worker-(\d+)\s+\|\s+http_req_failed\.+: ([\d\.]+)%')

    workers = {}

    print(f"Parsing {filepath}...")

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                # remove ANSI color codes if present
                clean_line = re.sub(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])', '', line)

                # Check RPS
                rps_match = rps_pattern.search(clean_line)
                if rps_match:
                    wid = rps_match.group(1)
                    if wid not in workers: workers[wid] = {}
                    workers[wid]['rps'] = float(rps_match.group(3))
                    workers[wid]['count'] = int(rps_match.group(2))
                    continue

                # Check Latency
                lat_match = latency_pattern.search(clean_line)
                if lat_match:
                    wid = lat_match.group(1)
                    if wid not in workers: workers[wid] = {}
                    workers[wid]['avg_lat'] = parse_time_to_ms(lat_match.group(2))
                    workers[wid]['med_lat'] = parse_time_to_ms(lat_match.group(3))
                    workers[wid]['p95_lat'] = parse_time_to_ms(lat_match.group(4))
                    continue

                # Check Fail rate
                fail_match = fail_pattern.search(clean_line)
                if fail_match:
                    wid = fail_match.group(1)
                    if wid not in workers: workers[wid] = {}
                    workers[wid]['fail_rate'] = float(fail_match.group(2))
                    continue

    except Exception as e:
        print(f"Error reading file {filepath}: {e}")
        return None

    # Filter out workers that didn't complete (missing metrics)
    completed_workers = {k: v for k, v in workers.items() if 'rps' in v and 'avg_lat' in v}
    return completed_workers

def plot_comparison(results):
    """
    Generates bar charts comparing metrics across different files (scenarios).
    'results' is a dict: { 'filename': {worker_id: stats, ...} }
    """

    scenarios = list(results.keys())
    # Calculate aggregates for each scenario
    agg_rps = []
    avg_latency = []
    p95_latency = []

    for sc in scenarios:
        workers = results[sc]
        if not workers:
            agg_rps.append(0)
            avg_latency.append(0)
            p95_latency.append(0)
            continue

        # Sum RPS across all workers (Total System Throughput)
        total_rps = sum(w['rps'] for w in workers.values())

        # Average the latency averages (Macro average)
        macro_avg_lat = np.mean([w['avg_lat'] for w in workers.values()])
        macro_p95_lat = np.mean([w['p95_lat'] for w in workers.values()])

        agg_rps.append(total_rps)
        avg_latency.append(macro_avg_lat)
        p95_latency.append(macro_p95_lat)

    # Setup Plot
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))

    # 1. Total Throughput
    ax1.bar(scenarios, agg_rps, color=['#4285f4', '#34a853', '#ea4335'])
    ax1.set_title('Total System Throughput (RPS)')
    ax1.set_ylabel('Requests / Second (Aggregated)')
    ax1.grid(axis='y', linestyle='--', alpha=0.7)
    for i, v in enumerate(agg_rps):
        ax1.text(i, v + (max(agg_rps)*0.01), f"{int(v)}", ha='center', fontweight='bold')

    # 2. Latency Comparison
    x = np.arange(len(scenarios))
    width = 0.35

    rects1 = ax2.bar(x - width/2, avg_latency, width, label='Avg Latency', color='#fbbc05')
    rects2 = ax2.bar(x + width/2, p95_latency, width, label='p95 Latency', color='#ea4335')

    ax2.set_title('Response Latency (ms)')
    ax2.set_ylabel('Milliseconds')
    ax2.set_xticks(x)
    ax2.set_xticklabels(scenarios)
    ax2.legend()
    ax2.grid(axis='y', linestyle='--', alpha=0.7)

    plt.tight_layout()
    plt.show()

def main():
    parser = argparse.ArgumentParser(description="Visualize K6 Worker Logs")
    parser.add_argument('files', metavar='F', type=str, nargs='+', help='Log files to parse')
    args = parser.parse_args()

    results = {}

    for file_path in args.files:
        # Create a readable label from filename (e.g., "k6-worker-cmsmode.log" -> "cmsmode")
        label = file_path.split('-')[-1].replace('.log', '')
        if 'k6' not in label: # fallback if splitting failed simple logic
             label = file_path

        data = parse_log_file(file_path)
        if data:
            print(f"  -> Found {len(data)} completed workers in {label}")
            results[label] = data
        else:
            print(f"  -> No usable data found in {label}")

    if results:
        plot_comparison(results)
    else:
        print("No data to plot.")

if __name__ == "__main__":
    main()