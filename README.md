# perspectives

Copyright zefrenchwan, 2025-2026.
MIT license

## What is it?

An event manager that registers incoming data to build a bitemporal history of information.

## Introduction

Events arrive about information changes: creating new elements, deleting some, or changing others.
This requires a **bitemporal model**: information arrives at a given time, and it changes what we know about the history of specific elements.
For instance, an event gives us information today about a person's wedding that happened three years ago.

| Name        | Description                       | Example                                      |
|-------------|-----------------------------------|----------------------------------------------|
| Event date  | when we receive the event         | we learned today Paul got married 3 years ago |
| Change date | when the action actually happened | Paul got married 3 years ago                 |

### States are data that vary over time

Basically, we store values as time-dependent content.
For instance, a given person has a name (Paul) assumed to be constant, and an address that may change over time.
The complete information is stored as a **state object**:
1. A global activity: the period of time during which the content object is valid. For instance, Paul's activity would be his lifetime.
2. Primitive values with a validity period. For instance, Paul's address is "123 Street" from `now() - 18 years` to `+oo`.

Accepted primitive types are strings, bools, integers, floats, and times. 
Other types are not accepted.

### Links are dynamic relationships that allow reification

Before detailing the model, let us pick some examples and determine what would be a good design. 

#### The time dependency
First one is obvious : time dependency.
Links depend on time.
For instance, Marie likes Paul since 8 years ago.
A link has a lifetime too. 

#### Roles more than predicates

"John likes Tiramisu" (an italian dessert). 
One might say : "easy, it is Likes(John, Tiramisu)".
But what about "John went to Paris by car"?

We use **roles** within links. 
A role for a link is the semantic meaning of related elements to that link. 
*Went(subject=John, object=Paris, role=destination, mode=car)* 
allows to include more than just predicates with union of possible options. 
Went is the kind of the link, usually a verb. 
It may be a time-dependent relation, too. 
For instance, *President* of a country. 
This example will explain why we have a time for the link, and one for the role. 
France is a country, and since 1958 (at least), it has a president.

```
France == president since 1958 ==> between 2017 and 2027 -> Macron
                                   between 2012 and 2017 -> Hollande
                                   ...
```

#### Links should be combined within other links

As an example, *Knows(subject=Marie, object=Likes(subject=John, object=Tiramisu))*.
We want links of links, so we use **reification**.
It means that links may be linked to other links for given roles. 
