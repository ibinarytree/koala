package main

var config_template = `port: 8080
prometheus:
  switch_on: true
  port: 8081
service_name: {{.Package.Name}}
register:
  switch_on: false
log:
  level: debug
  path: ./logs/
limit:
  switch_on: true
  qps: 50000
`
