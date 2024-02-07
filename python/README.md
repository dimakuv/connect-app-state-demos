# Python checkpoint/restore example

This directory contains a Makefile and a manifest template for running a simple
checkpoint/restore Python application in Gramine.

The Python application (`main.py`) creates a Person object (with name and age)
and goes into an infinite loop which increments the person's age every second.
At any point, a user can send a `SIGINT` (Ctrl-C) or a `SIGTERM` (`kill <pid>`)
signal, and the application will dump its Person object (aka application's
state) to a file and terminate.

To restore the application's state, the application must be started with
`RESTORE_FROM_FILE=<filename> ./python3 main.py`.

# Building

- Run `make` in the directory.
- Add `DEBUG=1` to build with Gramine log level set to `all`.
- Add `SGX=1` to build with Gramine-SGX support.

# Running

- Natively:
  1. Fresh start: `python3 main.py --name="Dmitrii Kuvaiskii" --age=35 --file=dump.dat`
  2. Terminate via Ctrl-C (or `SIGTERM` from another terminal)
  3. Restore from checkpoint: `RESTORE_FROM_FILE=dump.dat python3 main.py`

- With `gramine-direct` (non SGX mode):
  1. Fresh start: `gramine-direct python3 main.py --name="Dmitrii Kuvaiskii" --age=35 --file=dump.dat`
  2. Terminate via `SIGTERM` from another terminal (note that Ctrl-C aka `SIGINT` is not supported)
  3. Restore from checkpoint: `RESTORE_FROM_FILE=dump.dat gramine-direct python3 main.py`

- With `gramine-sgx` (SGX mode):
  1. Fresh start: `gramine-sgx python3 main.py --name="Dmitrii Kuvaiskii" --age=35 --file=dump.enc.dat`
  2. Terminate via `SIGTERM` from another terminal (note that Ctrl-C aka `SIGINT` is not supported)
  3. Restore from checkpoint: `RESTORE_FROM_FILE=dump.enc.dat gramine-sgx python3 main.py`

# How does it work internally?

We register a SIGINT and SIGTERM signal handler. Then the main thread goes into
an inifinite loop, incrementing the person's age every second.

As soon as the signal arrives, the main thread is interrupted and the control is
handed to the `handler()` function. This function performs process termination
-- as part of termination it serializes the internal app state and dumps into
the file.

# TODOs

Need a monotonic counter service:
- Upon checkpoint, the application must atomically increment-and-fetch the
  counter, and store the fetched value in the header of the file. Upon restore,
  the application must atomically fetch-and-increment the counter, and verify
  the value in the header of the loaded file against the fetched value.
- The service must run on e.g. localhost and be untrusted, for demo purposes.
  For example, a Redis server with a well-known key name.
