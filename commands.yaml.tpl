short: manage {{.ServerType}} instances together with attached ZFS disk devices
long: |
  lxops launches {{.ServerType}} containers from config (lxops) files that 
  specify how the instance should be launched and configured.
  
  Goals:
  - launch instances from a config file
  - configure an instance, using cloud-config files
  - manage instance filesystems along with the instance
  - rebuild an instance from a new image, while preserving its filesystems and configuration

  Devices are attached to an instance via an instance profile.
  For information about the config file, run "lxops config -h".
commands:
  configure:
    short: configure an existing container
    use: <config-file> ...
    examples:
    - configure -name c1 demo.yaml
  image:
    commands:
      instances:
        short: list image aliases for containers
        long: |
          print a table of all intances with columns:
          - image aliases
          - instance name
  cloudconfig:
    short: (applies one cloud-config file to instances)
    use: "[instance]..."
    examples:
    - -ostype alpine -f myconfig.cfg mycontainer
    - -ostype alpine mycontainer < myconfig.cfg
    long: |
      Uses the {{.ServerType}} API to apply the config to {{.ServerType}} instances.
      See also the "instance cloudconfig" command.
      The following cloud-init modules (sections) are supported and applied in this order:
        - packages
        - write_files (defer: false)
        - users
        - runcmd
        - write_files (defer: true)
      See github.com/melato/cloudconfig
  instance:
    short: {{.ServerType}} instance utilities
    commands:
      addresses:
        short: export network addresses for all containers
      cloudconfig:
        short: applies cloud-config files to instances
        use: "[cloud-config-file]..."
        examples:
        - -ostype alpine -i mycontainer -f config.cfg
        - -ostype alpine -i mycontainer < config.cfg
        - -ostype alpine -i mycontainer config.cfg...
        - -ostype alpine -f config.cfg instance...
        long: |
          configures {{.ServerType}} instances by applying cloud-init config to them,
          using the {{.ServerType}} API.
          Can either apply one cloud-init file to multiple instances
          or apply multiple cloud-init files to one instance.

          The following cloud-init modules (sections) are supported and applied in this order:
            - packages
            - write_files (defer: false)
            - users
            - runcmd
            - write_files (defer: true)
          See github.com/melato/cloudconfig
      hwaddr:
        short: export hwaddr for all containers
      number:
        short: assign numbers to containers
        use: -first <number> [-a] [-r] [-project <project>] <container>...]
      network:
        short: print container network addresses
        use: <container>
      profiles:
        short: print container profiles
        use: <container>
      devices:
        short: print container disk devices
        use: <container>
      publish:
        short: publish an instance into an image
        use: <instance> <snapshot> <alias>
      wait: 
        short: wait until all the requested containers have an ipv4 address
        use: <container>...
  create-devices:
    short: create devices
  create-profile:
    short: create lxops profile for instance
  delete:
    short: delete a container
    use: <configfile>...
    long: |
      delete a stopped container and its profile.
  destroy:
    short: delete a container and its filesystems
    use: <configfile>...
    long: |
      destroy is like delete, but it also destroys container filesystems
      that have the destroy flag set.  Other filesystems are left alone.
  config:
    short: lxops file utilities
    long: |
      config sub-commands examing lxops files.
      An lxops file is a yaml file starting with a version comment.
      The latest version is {{.ConfigVersion}}.
      Some earlier versions are also supported.
      Documentation of the latest config version is provided by the help commands.
    commands:
      formats:
        short: print supported config formats        
      parse:
        short: parse a config file
        use: <config-file>
      print:
        short: parse and print a config file
        use: <config-file>
      properties:
        short: print config file properties
        use: <config-file>
      script:
        short: print the body of a script
        use: <config-file> <script-name>
      includes:
        short: list included files
  i:
    short: show information about an instance/config
    commands:
      project:
        short: print instance project
        use: <container>
      description:
        short: print instance description
        use: <config-file>
        examples:
        - test.yaml
      devices:
        short: print instance devices
        use: <config-file>
      filesystems:
        short: print instance filesystems
        use: <config-file>
      properties:
        short: print instance properties
        use: <config-file>
      verify:
        short: verify instance config
        use: <config-file> ...
        examples:
        - verify *.yaml
  launch:
    short: launch an instance
    use: <config-file> ...
    examples:
    - launch php.yaml
  profile:
    short: profile utilities
    commands:
      apply:
        short: apply the config profiles containers
        use: <config-file> ...
      diff:
        short: compare container profiles with config
        use: <config-file> ...
      exists:
        short: check if a profile exists
        use: <profile>
      export:
        short: export profiles to yaml files
        use: <profile> ...
      import:
        short: import profiles from yaml files
        use: <file> ...
        long: |
          the name of the profile is the last element of the file path
      list:
        short: list config profiles
        use: <config-file>
      reorder:
        short: reorder container profiles to match config order
        use: <config-file> ...
  rebuild:
    short: rebuild an instance
    use: <config-file> ...
    long: |
      Rebuild replaces the instance image with the one specified in the config file,
      preserving the instance configuration.
      It also applies the cloud-config files specified in the config file.
      The image will be left in the Running state.
      See {{.ServerType}} rebuild.
  rename:
    short: rename an instance and its filesystems
    use: <configFile> <newname>
    long: Renames the instance, its filesystems, and its devices profile
  snapshot:
    short: snapshot instance filesystems
  rollback:
    short: rollback instance filesystems
  property:
    short: manage global properties
    long: |
      Properties can be located in:
      - Global Properties File
      - Instance Properties, inside the config .yaml file
      - Command Line
      Command line properties override instance and global properties.
      Instance properties override global properties.
    commands:
      list:
        short: list global property value
      file:
        short: print the filename of the global properties
      set:
        short: set a global property
        use: <key> <value>
      get:
        short: get a global property
        use: <key>
  ostypes:
    short: (list supported ostypes)
  export:
    short: export instance filesystems
    use: <config.yaml>
    long: |
      export the filesystems of an instance to tar.gz files
  import:
    short: import instance filesystems
    use: <config.yaml>
    long: |
      import the filesystems of an instance from tar.gz files
  help:
    short: documentation on lxops configuration
  copy-filesystems:
    short: (copy zfs filesystems from instance to another)
    long: |
      Uses ssh and zfs send/receive to copy a snapshot of the instance filesystems
      between hosts.
      If a short snapshot name is provided, it is used in the zfs send commands.
      Otherwise, a new snapshot is generated using "lxops snapshot".
      
      		[ssh <from-host>] lxops snapshot -s <generated-snapshot> ...

        , using the current date,
      and a snapshot is created in the source filesystems.
      
      Assumes that the config file is at the same path on the other host,
      and has the same filesystems.
