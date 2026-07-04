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
| Record date | when we receive the event         | we learned today Paul got married 3 years ago |
| Actual date | when the action actually happened | Paul got married 3 years ago                 |

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

# Concepts

We extensively use time and a bitemporal model. 


## Periods to ease time management

The *age* of a person is basically a function that takes a date and returns an integer. 
Altought it is true, we usually think in term of periods : when we had that job, when we got married, etc. 
A **period** is a finite union of time intervals.
For instance, the period `[2020-01-01, 2020-01-02]` is a period of two days.
For instance, a president's mandate may be `[2016, 2020] UNION [2024, 2028]`.

## Dynamic mappings

Consider `T` the time dimension, it is a totally ordererd set.
Given a set `A`, we can define : 
* a function `f : T -> A` that associates a value to a time. 
* a relation `R : T -> P(A)` that associates a set of values to a time. Mathematically, a relation is a subset of the cartesian product `T x A` with elements of `A` that may occur many times for the same element of `T`. 

### Functions of time and time dependent relations

We distinguish **dynamic relations** and **dynamic functions**, both being **dynamic mappings**.
Mapping is then the *general* term, whereas *function* is the mathematical approach. 

| Name | Description                                                           | Example |
| --- |-----------------------------------------------------------------------| --- |
| Dynamic mapping | Association between elements over time, in general                    | |
| Dynamic relation | Given a period, a set of elements                                     | Friends at a party |
| Dynamic function | Function that takes a moment as parameter and returns an unique value | CEO of a company |



### The image sets of mappings 

To describe an information, we use dynamic mappings of values. 
Values may be : 
* **primitive** and then describe information such as local content : an age, a name, an address. 
* **references** and then point to identifiable information : a person, a company, a product. 

| Type of mapping | Type of variable | Description of typical use                                           |
|-----------------|------------------|----------------------------------------------------------------------|
| Function        | Primitive        | Define time dependent unique attribute (age)                         |
| Function        | Reference        | Define time dependent unique reference (husband)                     |
| Relation        | Primitive        | Define time dependent multiple attributes (hobbies, favorite movies) |
| Relation        | Reference | Define time dependent set of references (friends)                    |


## Bitemporal model : terminology and implications

Based on [Martin's Fowler bitemporal description](https://martinfowler.com/articles/bitemporal-history.html), we will use this terminolgy : 
1. **record date** : the date an event notified a change happened
2. **actual date** : the date the actual change happened

For instance, a record arrives now (record date), informing the system about the birth of X (actual date), 3 days ago.
Still from [this source](https://martinfowler.com/articles/bitemporal-history.html), usual SQL terms are confusing, and we agree. 


