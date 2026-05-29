# perspectives

Copyright zefrenchwan, 2025
MIT license

## What is it ? 

An event manager to manage incoming information. 
It registers incoming information and build a state for information over time. 

## Concepts

To manage real world information, this project uses concepts such as traits to regroup common behavior, objects, and links. 
Traits represent concepts, objects are instances of traits, and links represent relationships between them.
Links and instances are valid during a given period of time. 
A formal class system defines what is an element (to determine whether current element is a trait, or an object, or a link). 


For instance, Paul, Marie and John are objects of the trait PERSON, and they are linked by the link FRIEND_OF.
Another example is the link composition. 
For instance, Knows(subject=Paul, object=Likes(subject=Marie,object=John)) 
Paul's lifetime is an interval `[ now() - 18 years, +oo[` and this link FRIEND_OF has a lifetime too, for instance `[now() - 3 years, +oo[`. 
