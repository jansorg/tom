<div style="color:red">This project is still in an early phase. It's usable and working but not yet recommended for production.</div>

# gotime

gotime is a command line application to track time and to simplify office work.
It's able to track time, to creat HTML and PDF reports and to create invoice drafts for a few web-based invoicing application

## Basic usage
A typical session looks like this:
```bash
gotime start acme
gotime stop
gotime report --month 0
```

## Documentation
Documentation about the command line is available at [docs/markdown](./docs/markdown/gotime.md)

## Data model
The data is stored in a few JSON files on disk. It's easy to backup and still fast.

### Projects
gotime supports nested projects.
The simplest form is a project without any subprojects.

If you need to bill a single client for several distinct project then it's better to use subprojects.
For example:
```bash
gotime create project client1 client1/web client1/backend
gotime start client1/web
gotime stop
gotime report --split project -p client
```

This will create a report on all projects which belong to client1 with the tracked time per project.

## Tracking time
### Start
### Stop
### Cancel
### View

## Reporting

### Date and time filters
### Splitting options

### Plain Text Reports

### HTML Reports

### PDF Reports

## Create Invoices
### Create invoices at sevdesk.com

## Import from other tools

### Import data from Watson
```bash
gotime import watson
```

### Import data from Fanurio
This needs a custom CSV export (to be documented).
```bash
gotime import fanurio complete-export.csv
```

## License
To be decided