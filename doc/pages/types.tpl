# lxops configuration types

{{range $i, $type := strings.List "Config" "Pattern" "HostPath" "Filesystem" "Device"}}
- [{{$type}}]({{$type}}.md)
{{- end}}

