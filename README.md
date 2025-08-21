# Sysy - Nginx Access Log converter to Prometheus Metrics

## Usage

### Nginx Configuration

Add Below Configuration to your existing nginx.conf. For example, check config directory in this repository.

```nginx
  log_format syslog_format '$http_host\t'
  '$request_method\t'
  '$status\t'
  '$request_completion\t'
  '$request_time\t'
  '$request_length\t'
  '$bytes_sent\t'
  '$request_uri';

  # Change <sysy-url> to the address of sysy app
  # (service name in docker or ip address)
  # and configured sysy port that run syslog server (default :5140)
  access_log syslog:server=<sysy-url>:<port>,nohostname syslog_format;
```

### Sysy Available Configuration

Configuration can be set using .env file if the code run from the source and via environment if using docker container. Here is the list available configuration:

- CONF_NGINX_TARGET_URL -> set nginx url target for stub_status (DEFAULT: http://nginx/nginx_status)
- CONF_SYSLOG_ADDR -> set syslog port number (DEFAULT: :5140)
- CONF_LOG_ENABLE -> to enable log mode in the app (DEFAULT: false)

For the example, you can see .env and docker-compose.yml in this repository. Please take a look that in the docker-compose.yml, the config in environment doesn't need "" (quote)!

### Grafana Dashboard

Import JSON that is available in this repository, inside the config directory.

<img src="./asset//grafana.png" alt="Grafana Dashboard"/>