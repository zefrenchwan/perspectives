# Rules to use the DSL 

## File structure 

1. Comments starts with `##`. The rest of the line is ignored
2. New group starts at position 0 of a line. Groups may be separated with a \n but it is not mandatory. Anything within a group starts after a space or tab character (except of course the first line)


For instance: 
```
first group
second group

third group
    still third group
    third group again
```