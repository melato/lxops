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
  For documentation about the config file, run "lxops help config".
  
  Commands with summaries in parenthesis are auxiliary utilities.
  They may change without maintaining backward compatibility.

commands:
  configure:
    short: configure an existing instance
    use: <config-file> ...
    examples:
    - configure -name c1 demo.yaml
  image:
    short: (utilities for {{.ServerType}} instances)
    commands:
      instances:
        short: list image aliases for containers
        long: |
          print a table of all intances with columns:
          - image aliases
          - instance name
  cloudconfig:
    short: apply cloud-config files to instances (new)
    use: "[instance]..."
    examples:
    - -ostype alpine -f myconfig.cfg mycontainer
    - -ostype debian mycontainer < myconfig.cfg
    long: |
      Uses the {{.ServerType}} API to apply the config to {{.ServerType}} instances.
      It does not lxops config files or properties.
      The following cloud-init modules (sections) are supported and applied in this order:
        - packages
        - write_files (defer: false)
        - users
        - runcmd
        - write_files (defer: true)
      See github.com/melato/cloudconfig.
      
      This command is marked as experimental, because it is new.
      It may turn out to be a core functionality of lxops.
  instance:
    short: ({{.ServerType}} instance utilities)
    commands:
      addresses:
        short: (export network addresses for all containers)
      cloudconfig:
        short: (applies cloud-config files to instances)
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
        short: (export hwaddr for all containers)
      info: 
        short: (print instance information)
        long: |
          Prints instance information, as reported by the {{.ServerType}} server.
          This is for informational purposes only.
      number:
        short: (assign numbers to containers)
        use: -first <number> [-a] [-r] [-project <project>] <container>...]
      network:
        short: (print container network addresses)
        use: <container>
      profiles:
        short: (print container profiles)
        use: <container>
      devices:
        short: (print container disk devices)
        use: <container>
      wait: 
        short: (wait until all the requested containers have an ipv4 address)
        use: <container>...
  publish:
    short: (create an image from an instance)
    use: <instance>[/<snapshot>]
    examples:
    - publish -c mariadb.yaml debian-mariadb-20260105-0931
    - publish -c mariadb.yaml debian-mariadb-20260105-0931/copy
    - publish -c mariadb.yaml --dry-run debian-mariadb-20260105-0931
    long: |
      publish is a wrapper for {{.ServerType}} publish.
      
      It extracts default image properties from the instance name and
      from the config file (that was used to build the instance).
      
      If there is no -release or -serial specified,
      it looks for a series of numbers and dashes in the instance name.
  create-devices:
    short: (create devices)
  create-profile:
    short: (create lxops profile for instance)
  delete:
    short: delete an instance
    use: <configfile>...
    long: |
      Delete a stopped instance and its profile.
      Do not touch its non-root filesystems.
  destroy:
    short: delete an instance and its filesystems
    use: <configfile>...
    long: |
      destroy is like delete, but it also destroys container filesystems
      that have the destroy flag set.  Other filesystems are left alone.
      Standalone devices without a filesystem are also left alone.
  config:
    short: (lxops file utilities)
    long: |
      config sub-commands examing lxops files.
      An lxops file is a yaml file starting with a version comment.
      The latest version is {{.ConfigVersion}}.
      Some earlier versions are also supported.
      Documentation of the latest config version is provided by the help commands.
    commands:
      convert:
        short: convert config files to the latest format
        use: <config-file>...
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
      packages:
        short: print packages to install
        use: <config-file>
      packages:
        short: print cloud-config files
        use: <config-file>
      script:
        short: print the body of a script
        use: <config-file> <script-name>
      includes:
        short: list included files
  i:
    short: (show information about an instance/config)
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
        long: |
          displays a table of filesystems specified in the config file.
          The columns are:
          FILESYSTEM: the filesystem identifier, as used in the config
          PATH: the filesystem path, with properties substituted
          PATTERN: The filesystem path, without properties substituted
          FLAGS: Filesystem flags: destroy, transient.
          
          see "help filesystem"
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
  extract:
    short: (extract devices from an image)
    use: <config-file>
    examples:
    - launch -name php php.yaml
    long: |
      The extract command copies files from an image to disk devices.
      It works for containers and images that use ZFS storage.

      Specifically, the command:
      - creates instance filesystems and device directories for an instance
      - launches a container from the image
      - finds the container root zfs filesystem origin
      - mounts this origin to a temporary directory
      - uses rsync to copy the lxops device directories from the
        temporary directory to the instance device directories.
      - shifts the uids and gids of the copied files, by the amount
        specified by the device-owner property (uid:gid)
      
      Example Use Case:
      Suppose config file example.yaml specifies:
        image: example-template
        device-template: example-template
        
      It also specifies a disk device, to be mounted on /var/log

      When you run "lxops launch example.yaml", it reports that the template
      device for /var/log is missing.
      
      You can create the template device using "extract":
        lxops extract -name extract-emplate example.yaml
      
      Then delete the example devices created by the previous incomplete launch:
        lxops destroy example.yaml
      
      and relaunch it properly:
        lxops launch example.yaml
                
      LIMITATIONS:      
      The "launch" should incorporate what the "extract" command does.
      
  create:
    short: (create container from image)
    use: <config-file>
    long: |
      This is an experimental version of launch that does not start the container.
      Its original purpose was to extract devices from an image.
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
    short: manage user properties
    long: |
      User properties in a file specified by the -properties flag,
      with a default location in the user's config directory.
      They can be managed by the these sub-commands, or edited manually.
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
    short: (documentation on lxops configuration)
    commands:
      topics:
        short: print documentation for various topics
        use: "[<topic>]"
        long: |
          If a topic is provided, print its content.
          Otherwise, print a list of topics.          
  copy-filesystems:
    short: (copy zfs filesystems from one instance to another)
    long: |
      Uses ssh and zfs send/receive to copy a snapshot of the instance filesystems
      between hosts.
      
      If a short snapshot name is provided, it is used in the zfs send commands.
      Otherwise, a new snapshot name is generated and created using "lxops snapshot":
        [ssh <from-host>] lxops snapshot -s <generated-snapshot> ...
      
      This command assumes that the config file is at the same path on the other host,
      and has the same filesystems.
      
