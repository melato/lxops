short: manage {{.ServerType}} containers together with attached ZFS disk devices
long: |
  lxops launches {{.ServerType}} containers and creates or clones ZFS filesystem devices for them.
  lxops launches an "instance" by:
    - Creating or cloning a set of ZFS filesystems
    - Creating and initializing a set of sub-directories under these filesystems
    - Creating an {{.ServerType}} profile with disk devices for these directories
    - Launching or copying an {{.ServerType}} container with this profile
    
  lxops can also install packages, create users, setup .ssh/authorized_keys for users,
  push files from the host to the container, attach profiles, and run scripts.
  
  One of its goals is to separate the container OS files from user files,
  so that the container can be upgraded by swapping its OS with a new one,
  instead of upgrading the OS in place.
  Such rebuilding  can be done by copying a template container
  whie keeping the existing container disk devices.
  
  The template container can be upgraded manually, using the OS upgrade procedure,
  or relaunched from scratch.
  
  A Yaml configuration file provides the recipe for how the container should be created.
  It can include other config files, so that common configuration
  can be reused across instances.
  
  Devices are attached to the container via an instance profile.
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
  instance:
    short: {{.ServerType}} instance utilities
    commands:
      addresses:
        short: export network addresses for all containers
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
    short: config .yaml utilities
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
    short: stop, delete, launch
    use: <config-file> ...
    long: |
      Rebuild stops, deletes, and relaunches the container.
      It preserves the previous hwaddr from the container,
      so the new container should have the same IP addresses as before.
  rename:
    short: rename an instance
    use: <configFile> <newname>
    long: Renames the container, its filesystems, and its devices profile
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
  copy-filesystems:
    short: copy zfs filesystems from instance to another.
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