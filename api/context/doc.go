/*
Contains set of routines which implementing the context storage in GO language.

When to use it

You can use it whenever you need to pass the value(s) to nested calls where it impossible or inconvenient to pass
through function arguments.

Is this when

Imagine that you are writing the HTTP application endpoint. That endpoint mapped to call on some function "A()",
which then performs a calls to a functions B, C, D and they have some their sub-calls E, F, G, H as well, like:

  A(request, response) -> B(request) -> E() -> F() -> H()
                       -> C() -> G()
                       -> D()

Then you figuring out that there is a case where the warning could happen within function "E()" and this information is
significant for a response (i.e. you need to pass that warning to function "A()"). You can update function B and E with
the reference argument or update the E function to pass warning with the returning value - but all of this will require
appropriate modifications in "parent" functions "A()" and "B()". Alternatively you can use global package variable and
put the message there, but global package variable would stay the same for a different "go-routines" and the function
"A()" would not have a chance to distinguish the request relation to data.

Consider that "H()" function (or any more nested function) could have the similar cases or for an instance the "H()"
function locates in external package and works through the interface but it requires HTTP request parameters to perform
the appropriate work, consider that these parameters are not static. Would you like to support function arguments to
pass them?

If not convinced, try to consider the following case: you are developing the database driver in separate package, the
other packages will work with your driver through the interfaces. The database supports the transactions. There are at
least 2 packages (shipping and payment) which are using the DB operations during the checkout process (another package).
Shipping and payment packages should save the generated information into database only in case when the checkout was
successful and the order was created, otherwise the operations should be rolled back. All packages are developing
independently by different organizations.

  [checkout]: receipt() -> [db] startTransaction()
                        -> [shipping] createShipping(...) -> [db] insertData(...)
                        -> [payment] createPayment(...) -> [db] insertData(...)
                        -> [checkout] placeOrder(...) -> error()
                        -> [db] rollbackTransaction()

The same database driver is used for different operations and the couple DB operations happens simultaneously in a
different go routines (cms operation, another checkout, etc.). How the database driver could distinguish the operations
it should rollback on "rollbackTransaction()" call? The transaction ID should be passed to each possible
package and later provided for each possible functions working with DB? Should these package be aware of the transaction
happens or that fact should stay known only between database driver and checkout packages?

One way or another you can omit this architecture questions nightmare with the contexts.


How it works

There is call stack you can get in GO application at any moment by call to runtime.CallersFrames(). It returns the
frames information created by go for each function call. It looks as follow:

      {PC}:{Entry}  - {File}:{Line}                - {Function}
  -------------------------------------------------------------------------
  17392491:17392432 - /Users/user/go/sample.go:421 - main.H
  17392415:17392384 - /Users/user/go/sample.go:417 - main.F
  17392367:17392336 - /Users/user/go/sample.go:413 - main.E
  17392319:17392288 - /Users/user/go/sample.go:409 - main.B
  17392260:17392208 - /Users/user/go/sample.go:405 - main.A
  17392768:17392688 - /Users/user/go/sample.go:425 - main.main
  16949355:16948832 - /usr/local/go/src/runtime/proc.go:200 - runtime.main

So each stack frame have "Entry" field which is the pointer to a function entity. That pointer is unique and persistent
for each function you have in application.

That is how it changes with the context usage.

  17392491:17392432 - /Users/user/go/sample.go:421 - main.H
  17392415:17392384 - /Users/user/go/sample.go:417 - main.F
  17392367:17392336 - /Users/user/go/sample.go:413 - main.E
  17392319:17392288 - /Users/user/go/sample.go:409 - main.B
  17392260:17392208 - /Users/user/go/sample.go:405 - main.A
  17409360:17409280 - /Users/user/go/sample.go:425 - main.main.func1          <- c.RunInContext(func() { A() }, nil)
  17401636:17401600 - /Users/user/go/context.go:150 - context.glob..func139   <- proxiesDict[17401636] = 139
  17389962:17389840 - /Users/user/go/context.go:302 - context.proxyLoop       <- proxiesDict[17389840] = 256
  17389915:17389840 - /Users/user/go/context.go:300 - context.proxyLoop       <- proxiesDict[17389840] = 256
  17392034:17391696 - /Users/user/go/context.go:386 - context.RunInContext    <- proxiesStart
  17392741:17392688 - /Users/user/go/sample.go:425 - main.main
  16949355:16948832 - /usr/local/go/src/runtime/proc.go:200 - runtime.main

You can see a set of call appeared in stack, it allows us to bind unique number/key and use it then in global variable.
So in the application we wrapping the call to function "A()" with context.RunInContext(func() { A() }, nil) which looks
for a free value in global "contexts" variable - it doing it with the regular loop from 0 to uint max value, once the
free key found in "contexts" variable it puts there the new "map[string]interface{}" instance which is the context.
RunInContext have deffer function which will release the key after the given function finishes. In the given example
this key is "256+256+139=651".

As you can see the context package have 255 similar function and knows their pointers, each pointer is mapped to a
sequenced number from 1 to 255, if the index is bigger than 255 it calls the recursive function "proxyLoop" reducing
the index on 256 and doing this as much while the rest could be mapped into existing functions (i.e. less then
proxyBase which equals length of "proxies" variable, currently it is 255 proxy functions).

Example:
    import "context"
    import "fmt"

    func A() {
            fmt.Printf("func A: x=%d\n", GetContext()["x"])
    }

    func B() {
          dict := GetContext()
          fmt.Printf("func B: x=%d\n", dict["x"])

            if val, ok := dict["x"].(int); ok {
                dict["x"] = val * 2
            }

            A()
    }

    func C() {
          dict := GetContext()
          fmt.Printf("func C: x=%d\n", dict["x"])
        dict["x"] = 5

            B()
    }

    func main() {
            RunInContext(C, map[string]interface{} {"x": -1})
    }


Result:
    func C: x=-1
    func B: x=5
    func A: x=10

*/
package context
