lxops has some support for {{.ServerType}} projects.

If you do not specify a project anywhere in a config file or the command line,
lxops will use the current project.

To determine what the current project is, it reads {{.ServerType}}
client configuration (usually in ~/.config/...).

If you specify a project on the command line, lxops will use that project,
instead of the current project.

If a config file has the "project" field set.