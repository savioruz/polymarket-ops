# /polymarket-copy — Execute a copy trade

Load modes/copy.md, modes/_shared.md, and config/profile.yml.

Run the full pre-flight checklist before any order.
Calculate position size based on profile settings.
Present the trade summary and wait for CONFIRM from the user.
Log to data/tracker.tsv after execution.

IMPORTANT: Never execute without user confirmation unless auto_execute=true in profile.
