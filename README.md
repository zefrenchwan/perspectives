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

### Traits define common behavior for instances

Traits are concepts that define common behavior for instances.
They represent abstract ideas or categories that can be applied to multiple instances. Traits provide a way to group similar objects together and define their shared characteristics.
A trait is **not** explicitly used to instantiate an object; instead, we apply a **duck typing** principle.
This duck typing allows us to treat objects that share common behavior as instances of a trait based on their current content, without explicitly defining their type.
Traits may apply to multiple instances, and instances may belong to multiple traits.

### Instances are linked to traits during a given period

Traits apply to an instance **during a given period of time**, if any.
For instance, a person can be a student during their school years, an employee during their working years, and a parent during their parenting years.

### Instances and traits share links during periods of time

Links represent relationships between any element within the system, **including other links**.
Links and instances are valid during a given period of time.

For instance, Paul, Marie, and John are instances of the trait `PERSON`, and they are connected by the link `FRIEND_OF`.
Another example is link composition.
For instance: `Knows(Subject=Paul, Object=Likes(Subject=Marie, Object=John))`
Paul's lifetime is an interval `[now() - 18 years, +oo[` and the `FRIEND_OF` link has a lifetime too, for instance `[now() - 3 years, +oo[`.

#### Example: Knowledge Representation
`Link:Works(Subject:TraitWorker, Object:TraitJob)`

#### Example: Inheritance Tree
`Link:Extends(Subject:TraitDessert, Object:TraitFood)`

#### Example: Knowledge Management
`Link:Knows(Subject=Paul, Object=Likes(Subject=Marie, Object=John))`