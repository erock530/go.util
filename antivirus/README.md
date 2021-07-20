This library was designed to interface with mcafee uvscan CLI command tool.

##Implementing
This code requires the uvscan software be installed on your machine.

Once imported, the antivirus scanner object can be declared and managed by
running
```
avscanner, err := GetAVScanner()
```

Scanning files can be performed on a given system file path by running
```
avscanner.Scan(FILE_PATH)
```

##Basic configuration

Required configs and their defaults:
```
VIRUS\_SCAN\_ENABLED (false): Enable virus scanning when true
VIRUS\_SCANNER\_PATH (/usr/local/bin/uvscan): Path for virus scanner
VIRUS\_SCANNER\_DAT\_PATH (/usr/local/uvscan/): Path for directory containing dat files
MAXFILESIZE (10): Examine only files smaller than the specified number of MB
MEMSIZE (1000): File size (in KB) to load into memory for scanning, limited by
a maximum file size.
NOCOMP (true): Do no scan self extracting executables by default.
```

All other configuration parameters can be passed in 
```
ADDITIONAL\_PARAMETERS: List of additional parameters and values
e.g. ADDITIONAL\_PARAMETERS='[mime,selected,xmlpath,PATHNAME,unzip]'
```

##Speeding up antivirus run times

When running this program, it is recommended that the .dat files are copied
over to memory and referenced there. Assuming the dat files are located in
"/usr/local/uvscan/uvscan-dats/",  this can be done with by executing the
following
```
cp -Rv /usr/local/uvscan/uvscan-dats /dev/shm/
VIRUS_SCANNER_DAT_PATH='/dev/shm/uvscan-dat/'
```
