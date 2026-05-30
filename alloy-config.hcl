local.file_match "app_logs" {
  path_targets = [
    {
      __path__ = "/var/log/app/*.log",
      job      = "campus_marketplace",
    },
  ]
}

loki.source.file "app_logs" {
  targets    = local.file_match.app_logs.targets
  forward_to = [loki.process.process_logs.receiver]
}


loki.process "process_logs" {
  stage.json {
    expressions = {
      level     = "level",
      timestamp = "timestamp",
      caller    = "caller",
      msg       = "msg",
    }
  }

  stage.labels {
    values = {
      level = "",
    }
  }

  stage.timestamp {
    source = "timestamp"
    format = "2006-01-02T15:04:05.000Z0700"
  }

  forward_to = [loki.write.to_loki.receiver]
}

loki.write "to_loki" {
  endpoint {
    url = "http://loki:3100/loki/api/v1/push"
  }
}
