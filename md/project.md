## Projects
lxops has some support for Incus projects.

If you do not specify a project anywhere in a config file or the command line,
lxops will use the current project.

To determine what the current project is, it reads Incus
client configuration (usually in ~/.config/...).

If you specify a project on the command line, lxops will use that project,
instead of the current project.

I don't use multiple projects, so multi-project functionality is not sufficiently tested.