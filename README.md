## unzar - extractor for Zip-Archiv (ZAR) files

A proprietary format by Peter Troxler. These files are DCL imploded with
some basic header.

#### Requires

To build:
 - golang 1.16

To install binary release:
 - See releases for builds for most major OS and architectures.

### Format Info
Reverse engineering of the format I found:

- n x m bytes - DCL Implode blob(s)
- n x Entry headers:
  - 1 byte - Filename length + 0x80
  - n bytes - Filename
  - 4 bytes - compressed blob length
- File footer: 7 bytes
  - 2 bytes - uint16(3) - unknown 16bit int, always "3"
  - 2 bytes - uint16(?) - unknown 16bit int, varies
  - 3 bytes - fileID - "\x50\x54\x26" - "PT&"
- EOF

### References
Archive format information: http://fileformats.archiveteam.org/wiki/ZAR_(Zip-Archiv)
DOS Archiver: https://www.sac.sk/download/pack/zip_ar26.zip

#### Author

- Kris Hunt (@ctfkris)