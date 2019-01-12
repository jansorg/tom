## tom start

starts a new activity for the given project ands adds a list of optional tags

### Synopsis

starts a new activity for the given project ands adds a list of optional tags

```
tom start <project> [time shift into past] [+tag1 +tag2] [flags]
```

### Examples

```
start acme 15m +onsite
```

### Options

```
      --create-missing   
  -h, --help             help for start
      --notes string     Optional notes for the new time frame
      --stop-on-start    
```

### Options inherited from parent commands

```
  -c, --config string     config file (default is $HOME/.gotime.yaml)
      --data-dir string   data directory (default is $HOME/.gotime)
```

### SEE ALSO

* [tom](tom.md)	 - tom is a command line application to track time.

###### Auto generated by spf13/cobra on 7-Jan-2019