# perspectives

Copyright zefrenchwan, 2025-2026.
MIT license

## What is it?

An event manager that registers incoming data to build a bitemporal history of information.

## Concepts

Events arrive about information changes: creating new elements, deleting some, or changing others.
This requires a **bitemporal model**: information arrives at a given time, and it changes what we know about the history of specific elements.
For instance, an event gives us information today about a person's wedding that happened three years ago.

| Name        | Description                       | Example                                      |
|-------------|-----------------------------------|----------------------------------------------|
| Event date  | when we receive the event         | we learned today Paul got married 3 years ago |
| Change date | when the action actually happened | Paul got married 3 years ago                 |

### Contents are data that vary over time

Basically, we store values as time-dependent content.
For instance, a given person has a name (Paul) assumed to be constant, and an address that may change over time.
The complete information is stored as a **content object**:
1. A global activity: the period of time during which the content object is valid. For instance, Paul's activity would be his lifetime.
2. Primitive values with a validity period. For instance, Paul's address is "123 Street" from `now() - 18 years` to `+oo`.

Accepted primitive types are strings, bools, integers, floats, and times. 
Other types are not accepted.

### Instances are identified objects with time-varying content

An instance is something one may identify, but that changes over time.
The changing part is a version of the content.
Identity is what we use to say: *this is this specific instance and not another*.
