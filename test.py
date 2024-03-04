import requests
import datetime
import time
import re

# Define the proxy server URL and parameters
url = "http://localhost:5000/proxyGet/"
params = {"app_name": "portainer-server", "app_port": "9000", "app_endpoint": "api/system/status"}

# Set the total number of requests to perform
total_requests = 500

# Perform 500 requests to the proxy server
start_time = time.time()
for _ in range(total_requests):
    try:
        response = requests.get(url, params=params, timeout=5)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")

# Add a small delay to ensure log captures the end of requests
time.sleep(0.05)
end_time = time.time()

# Calculate the total duration of the Python script
duration = end_time - start_time

# Convert start and end times to datetime objects for log analysis
start_datetime = datetime.datetime.utcfromtimestamp(start_time)
end_datetime = datetime.datetime.utcfromtimestamp(end_time)

# Read and analyze the log file within the time range of the script execution
log_file_path = "proxy.log"
with open(log_file_path, "r") as file:
    log_lines = file.readlines()

# Initialize variables for log analysis
total_requests_in_log = 0
total_duration_in_log = 0
successful_requests_in_log = 0
failed_requests_in_log = 0
other_log_entries = []

# Iterate through each line in the log file
for line in log_lines:
    # Extract timestamp from the log line
    timestamp_str = line.split(" ")[1] + " " + line.split(" ")[2]
    log_time = datetime.datetime.strptime(timestamp_str, "%Y/%m/%d %H:%M:%S.%f")

    # Check if the log line is within the time range of the script execution
    if start_datetime <= log_time <= end_datetime:
        # Check different log patterns and perform corresponding actions
        if "Proxy request duration:" in line:
            # Extract duration in seconds from the log line
            duration_in_log = float(line.split("duration:")[1].strip().replace("ms", "e-3").replace("µs", "e-6"))
            total_duration_in_log += duration_in_log
        elif "Proxy request status code:" in line:
            # Extract status code from the log line
            status_code = int(line.split("status code:")[1].strip())
            if status_code == 200:
                successful_requests_in_log += 1
            else:
                failed_requests_in_log += 1
        elif re.match(r".*http://.*", line):
            # Count total requests based on the specified pattern
            total_requests_in_log += 1
        else:
            # Store lines that don't match the expected patterns in a temporary list
            other_log_entries.append(line)

# Display detailed statistics and log analysis
print(f"Total Requests (Python Script): {total_requests}")
print(f"Total Time (Python Script): {duration:.3f} seconds")
print(f"Total Requests (from log file): {total_requests_in_log}")
print(f"Successful Requests (from log file): {successful_requests_in_log}")
print(f"Failed Requests (from log file): {failed_requests_in_log}")
print(f"Other Log Entries (from log file): {other_log_entries}")
print(f"Total Time (from log file): {total_duration_in_log:.3f} seconds")

# Calculate and display average time for successful requests in milliseconds
if successful_requests_in_log > 0:
    average_time_in_log = total_duration_in_log / successful_requests_in_log * 1000
    print(f"Average Time for Successful Requests (from log file): {average_time_in_log:.3f} milliseconds")
else:
    print("No successful requests in the log file.")