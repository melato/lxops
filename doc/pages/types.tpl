# lxops configuration types

{{range $i, $type := strings.List "Config" "Pattern" "HostPath" "Filesystem" "Device"}}
- [{{$type}}]({{$type}}.md)
{{- end}}

# Conditional configuration
*include* and *cloud-config-files* paths in lxops files go through variable substitution
and are used only if all referenced variables exist.

The characters "|;,:" are reserved.  They are not allowed in these path names.
