- Add --auto-start to immediately start the command
  (which will then be stopped according to stop-after, etc).

# v0.13 (2019-10-04)

- Pass local and remote addresses to "exec:" commands.
- Don't error for destination addresses of less than 5 characters.

# v0.12 (2019-08-18)

- Allow destination address to come from a command.
- Don't retry connections for non-temporary address errors (missing port, etc).

# v0.11 (2019-02-19)

- Wait until all clients disconnect before starting any "stop-after" timers.
- Fix cli usage description for -proxy-sep.

# v0.10 (2018-10-14)

- Use SIGTERM as default stop signal.
- Do not handle signals that are already ignored
  (for example: SIGINT when not run in the foreground).
- Make starting/stopping output simpler and more consistent.

# v0.9 (2018-03-19)

- Remove deprecated -listen option
- Remove deprecated JSON spellings in config file

# v0.8 (2018-01-07)

- Parse "-config" file as yaml.
  Keys should be lowercased with underscores ("stop_after").
  The JSON spellings ("StopAfter") are deprecated.

# v0.7 (2017-12-30)

- Add "-wait-after-start" option for services that aren't quite ready at
  the moment their port is open.

# v0.6 (2017-12-29)

- Set and signal process group for service commands (unix).

# v0.5 (2017-12-26)

- Allow "-proxy" to accept addresses like "interface:eth0:5001"
  and translate that to all addresses on that interface.
- Add "-stop-after" and "-stop-signal" to signal the command to stop
  after a period of inactivity.
- Exit immediately when listen fails for any address.

# v0.4 (2017-12-21)

- Add "-config" flag for json config file.
  This enables configuring multiple services.
- Deprecate "-listen".
  Instead pairs of addresses can be specified with -proxy "from to".
- Add "-proxy-sep" to configure the separator character used with "-proxy".
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
