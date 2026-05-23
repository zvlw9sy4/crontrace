# crontrace

> Parses and visualizes cron schedules with next-run predictions and conflict detection.

---

## Installation

```bash
go install github.com/yourusername/crontrace@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/crontrace.git && cd crontrace && go build ./...
```

---

## Usage

```bash
# Analyze a single cron expression
crontrace parse "*/5 * * * *"

# Load and visualize schedules from a file
crontrace analyze --file schedules.txt

# Detect conflicts across multiple expressions
crontrace conflicts "0 * * * *" "*/30 * * * *" "0 0 * * *"
```

**Example output:**

```
Expression : */5 * * * *
Description: Every 5 minutes
Next runs  :
  [1] 2024-11-01 14:05:00
  [2] 2024-11-01 14:10:00
  [3] 2024-11-01 14:15:00

⚠ Conflict detected: "0 * * * *" and "*/30 * * * *" overlap at 14:30, 15:00, ...
```

---

## Features

- Parse standard 5-field cron expressions
- Predict the next N scheduled run times
- Detect scheduling conflicts between multiple jobs
- Human-readable schedule descriptions

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

This project is licensed under the [MIT License](LICENSE).