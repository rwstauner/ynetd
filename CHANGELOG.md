# v0.4 (Unreleased)

- Make command optional (just a port forwarder).

# v0.3 (2017-12-09)

- Make output simpler and more consistent.
- Limit signal handling to HUP, INT, TERM so that the default
  behavior is used for any other signals (STOP, CONT, etc).
- Reap all child processes.

# v0.2 (2017-11-15)

- Output version with `--version` flag.
- Fix signal handling when command restarts.
- Automate test suite and release building.

# v0.0.1 (2017-10-23)

Initial Release
