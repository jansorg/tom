## tom edit project

edit properties of a project

```
tom edit project fullName | ID [flags]
```

### Options

```
  -h, --help                    help for project
      --hourly-rate string      Optional hourly rate which applies to this project and all subproject without hourly rate values
  -n, --name string             update the project name
      --name-delimiter string   Delimiter used in full project names (default "/")
      --note-required string    An optional flag to enforce a note for time entries of this project and all subprojects, where this setting is not turned off.
  -p, --parent string           update the parent. Use an empty ID to make it a top-level project. A project keeps all frames and subprojects when it's assigned to a new parent project.
```

### Options inherited from parent commands

```
      --backup-dir string   backup directory (default is $HOME/.tom/backup)
  -c, --config string       config file (default is $HOME/.tom/tom.yaml)
      --data-dir string     data directory (default is $HOME/.tom)
      --iso-dates           use ISO date format instead of a locale-specific format (default is false)
```

### SEE ALSO

* [tom edit](tom_edit.md)	 - edit properties of projects or frames

###### Auto generated by spf13/cobra on 19-Jul-2022
