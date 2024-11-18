<h1>
  <span>cscheck</span>
</h1>

**Using `cscheck` for checking and get metrics for your Cobalt Strike teamserver.**

## Overview

`cscheck` uses the agscript ([Aggressor Script](https://hstechdocs.helpsystems.com/manuals/cobaltstrike/current/userguide/content/topics/agressor_script.htm)) engine to execute CNA scripts and get status, metrics from teamserver.

## Features
* Administration availability check 
* Checking the number of Listeners
* Get the names of all Listeners

You can see CNA scripts in `scripts` dir.

URL     | Description |
---------|-------------|
/healthz | Exposes Cobalt Strike teamserver admin status (down or up). |
/metrics | Exposes Cobalt Strike teamserver metrics in Prometheus exporter format. |

## Installation

The `cscheck` listens on HTTP port 8000 by default. See the `--help` output for more options.

### Prerequisites

* Docker on local machine
* golang compiler
* Make
### Build Docker image
1. Clone repo.
```bash
git clone https://github.com/3ayazaya/cscheck
cd cscheck
```

2. Build Docker image with `cscheck`.

```bash
make build-image
```

3. Run Docker image

```bash
docker run --rm -it \
    -p 8000:8000 \ 
    3ayazaya/cscheck:1.0 \
    -password changeme \
    -ip 192.168.20.150 \
    -port 50050 \
    -user checker \
    -bind 0.0.0.0:8000
```

4. Check `csckeck` status.

```bash
curl -v http://127.0.0.1/healthz
```
Prometheus exporter metrics.
```bash
curl -v http://127.0.0.1/metrics
```