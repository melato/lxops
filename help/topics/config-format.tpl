lxops supports multiple configuration file formats.  There are currently two supported formats:
  - lxops-v1 - The latest format.
  - lxdops - An older format that lxdops used

backward compatibility is maintained by using migrators that convert a format to a newer format.
- "lxdops" files are converted to "lxops-v1" files.
- If and when there is a lxops-v2 format, there could be an lxops-v1 migrator that converts lxops-v1 files to lxops-v2 files.

Format migrators are chained, so lxdops files will be converted to lxops-v1 files and then to lxops-v2 files.
Therefore, all previous formats should be supported.

Format migrators convert raw bytes to raw bytes.  They do not need to depend on lxops data types.
Instead, they can unmarshal the content to generic yaml data types, manipulate them, and marshal them back to binary datay.
