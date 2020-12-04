```
Time Flies (tf) is a tool for budgeting focus time.

Usage:
  tf [command]

Available Commands:
  edit        Edit the log file.
  help        Help about any command
  tidy        Reformats log to spark joy.
  todo        List TODO entries.
  tots        Output weekly focus totals.

Flags:
  -c, --config string   Config file. JSON serialization of pkg/cmd/Config.
  -f, --focus string    Focus on a particular label group.
  -g, --group strings   Group entries by labels.
  -h, --help            help for tf
  -l, --log string      Log file.
  -p, --period string   Aggregation period.

Use "tf [command] --help" for more information about a command.
```
# Overview

Time Flies (`tf`) is a text-oriented tool for keeping track of focus time. The center of `tf` is the log, a plaintext file into which daily activities are written, one per line. The log is as a record with metadata (time spent, focus category, etc...) and a TODO list.

## Example Log

```
Date: Nov 23 2020

find other controllers with the same defect http://bug/1234   ## cat=customer sub=ops f=1h
prepare for eng review http://link/my-design                  ## cat=primary sub=thing-one f=3h
vision for my tool https://chat/group/abcd/message/1234       ## cat=community sub=docs
provide customer with advice on the thing http://bug/5678     ## cat=customer sub=ops
read about https://github.com/interesting-project             ## cat=community sub=group
create release coordination doc http://doc/release-the-thing  ## cat=primary sub=thing-two f=1h
meet about repeated bug http://bug/9012                       ## cat=customer sub=ops t=45m
# review http://doc/really-interesting-idea                   ##
# look at http://bug/3456 boxes within boxes                  ##
# bug report cleanup                                          ##
```

## Philosophy

# Commands

## tots

## edit

## tidy

## todo

# Customization
