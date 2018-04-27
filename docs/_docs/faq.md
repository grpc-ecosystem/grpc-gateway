---
category: documentation
---

# FAQ

> Why are the models in the swagger specification prefixed with the last part of the proto package name?

The reason to generate the prefixes is that we don't have a guaranteed unique namespace. If two packages produce different Foo messages then we will have trouble.

> Why not strip the prefix?

When a message is added which happens to conflict with another message (e.g. by importing a message with the same name from a different package) it will break code that is very far away from the code that changed. This is in an effort to adhere to the [principle of least astonishment](https://en.wikipedia.org/wiki/Principle_of_least_astonishment).
