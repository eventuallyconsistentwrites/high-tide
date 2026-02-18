import json
import re
import argparse
import matplotlib.pyplot as plt
from datetime import datetime
import collections

def parse_iso8601(time_str):
    """
    Parses timestamps like "2026-02-17T21:10:17.594995917Z".
    Handles the nanosecond precision (9 digits) by truncating to microseconds (6 digits)
    which Python's datetime can handle.
    """
    try:
        # Handle 'Z' usually meaning UTC
        time_str = time_str.replace('Z', '+00:00')

        # If there are fractional seconds
        if '.' in time_str:
            main_part, frac_part = time_str.split('.')
            # Extract timezone if present after fractional part
            if '+' in frac_part:
                frac, tz = frac_part.split('+')
                # Truncate to 6 digits (microseconds)
                frac = frac[:6]
                clean_str = f"{main_part}.{frac}+{tz}"
            elif '-' in frac_part: # unlikely for ISO but possible
                frac, tz = frac_part.split('-')
                frac = frac[:6]
                clean_str = f"{main_part}.{frac}-{tz}"
            else:
                # No timezone or implied UTC
                frac = frac_part[:6]
                clean_str = f"{main_part}.{frac}"
        else:
            clean_str = time_str

        return datetime.fromisoformat(clean_str)
    except ValueError:
        return None

def parse_server_log(filepath):
    """
    Parses the high-tide-server log file.
    Returns a sorted list of datetime objects representing request timestamps.
    """
    # Regex to find the JSON part after the pipe
    # Matches: high-tide-server-1  | {"time":...}
    log_pattern = re.compile(r'^.*?\|\s+(\{.*\})$')

    timestamps = []

    print(f"Parsing {filepath}...")

    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            for line in f:
                match = log_pattern.search(line)
                if match:
                    json_str = match.group(1)
                    try:
                        data = json.loads(json_str)

                        # We only care about request events
                        if data.get('msg') == 'checking rate limit':
                            ts_str = data.get('time')
                            if ts_str:
                                dt = parse_iso8601(ts_str)
                                if dt:
                                    timestamps.append(dt)

                    except json.JSONDecodeError:
                        continue # Skip malformed lines

    except Exception as e:
        print(f"Error reading {filepath}: {e}")
        return []

    timestamps.sort()
    return timestamps

def calculate_rps(timestamps):
    """
    Converts a list of timestamps into a dictionary of { relative_second: count }.
    """
    if not timestamps:
        return {}, 0

    start_time = timestamps[0]
    duration_seconds = (timestamps[-1] - start_time).total_seconds()

    # Bucket into seconds
    buckets = collections.defaultdict(int)

    for ts in timestamps:
        # Calculate seconds since start of this specific log
        delta = (ts - start_time).total_seconds()
        sec_bucket = int(delta)
        buckets[sec_bucket] += 1

    return buckets, duration_seconds

def plot_server_rps(results):
    """
    Plots the RPS curves for multiple files.
    """
    plt.figure(figsize=(14, 7))

    colors = ['#4285f4', '#34a853', '#ea4335', '#fbbc05'] # Google colors

    for i, (label, timestamps) in enumerate(results.items()):
        if not timestamps:
            continue

        rps_data, duration = calculate_rps(timestamps)

        # Prepare X and Y lists
        # We limit to reasonable duration (e.g. 90s) to avoid long tails of idle server
        x = sorted(rps_data.keys())
        y = [rps_data[k] for k in x]

        # Sliding window average for smoothing (optional, makes graphs readable)
        window_size = 3
        y_smooth = []
        for j in range(len(y)):
            start = max(0, j - window_size + 1)
            window = y[start:j+1]
            y_smooth.append(sum(window) / len(window))

        plt.plot(x, y_smooth, label=f"{label} (Avg: {int(len(timestamps)/duration)} RPS)", linewidth=2, color=colors[i % len(colors)])

    plt.title('Server-Side Traffic (Requests Processed/Checked)')
    plt.xlabel('Time Elapsed (seconds)')
    plt.ylabel('Requests Per Second (RPS)')
    plt.legend()
    plt.grid(True, which='both', linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.show()

def main():
    parser = argparse.ArgumentParser(description="Visualize High-Tide-Server Logs")
    parser.add_argument('files', metavar='F', type=str, nargs='+', help='Server Log files to parse')
    args = parser.parse_args()

    results = {}

    for file_path in args.files:
        # Generate label from filename
        # e.g. "server-cmsmode.log" -> "cmsmode"
        label = file_path.split('-')[-1].replace('.log', '')
        if 'server' in label:
            label = file_path

        data = parse_server_log(file_path)
        if data:
            print(f"  -> Found {len(data)} requests in {label}")
            results[label] = data
        else:
            print(f"  -> No request data found in {label}")

    if results:
        plot_server_rps(results)
    else:
        print("No data to plot.")

if __name__ == "__main__":
    main()