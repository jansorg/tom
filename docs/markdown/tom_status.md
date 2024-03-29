## tom status

Displays when the current project was started and the time spent...

```
tom status [flags]
```

### Options

```
  -d, --delimiter string        Delimiter to separate flags on the same line. Only used when --format is specified. (default "\t")
  -f, --format string           Properties to print for each active frame. Possible values: id,projectID,projectName,projectFullName,projectParentID,startTime
  -h, --help                    help for status
      --name-delimiter string   Delimiter used in the full project name (default "/")
  -v, --verbose                 Print details about the currently stored projects, tags and frames
```

### Options inherited from parent commands

```
      --backup-dir string   backup directory (default is $HOME/.tom/backup)
  -c, --config string       config file (default is $HOME/.tom/tom.yaml)
      --data-dir string     data directory (default is $HOME/.tom)
      --iso-dates           use ISO date format instead of a locale-specific format (default is false)
```

### SEE ALSO

* [tom](tom.md)	 - tom is a command line application to track time.
* [tom status projects](tom_status_projects.md)	 - Prints project status

###### Auto generated by spf13/cobra on 19-Jul-2022
