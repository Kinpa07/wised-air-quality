Build the Docker image (the binary is compiled inside the image, so no local Go or C toolchain is required):

`docker build -t wised .`

Then run it:

`docker run -p 8080:8080 wised`

## Assignment

This repo is the backend for a take-home assignment. The task brief is in
[`docs/assignment.md`](docs/assignment.md).

A device-fleet simulator is provided as a Docker image (delivered as a TAR
archive) to exercise the service — it enrolls clients and streams sensor
readings (including resent duplicates). Load it and run it against your running
service:

```bash
docker load -i wised-interview-simulator.tar
docker run --rm wised-interview-simulator:latest
```

With no arguments it runs the standard full-fleet scenario against
`http://host.docker.internal:8080` (the service on your host). Pass flags to
customise the run — see the simulator image's own README for all options and the
Linux networking note.
