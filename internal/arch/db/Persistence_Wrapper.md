# Persistence Wrapper

When generics are supported in Go, we will define the following types to support the clean integration of the domain
model with persistence frameworks:

**RecCtx[T]**

```go
// RecCtx is a type that holds platform-specific database record context information,
// e.g., an optimistic locking token and/or a record ID.  DAFs may accept this type as
// a parameter or return this type, together with domain entity types.
// This type is parameterized to provide type safety, i.e., to prevent passing a RecCtx[U]
// on a call that involves entity type V.
type RecCtx[T any] struct {
    Rc interface{}
}

// Example A1 -- DAF signature
func PersonUpdateDaf(person Person, rc RecCtx) (Person, RecCTx, error)

// Example A2 -- usage
// With separation of entity and RecCtx
person, recCtx, err := PersonReadDaf(name)
if err != nil { return err }
person = SomeBusinessFunctionBf(person)
person, recCtx, err = PersonUpdateDaf(person, recCtx)
```

**Pw[T]**

```go
// Pw wraps a domain entity and RecCtx together.  It can be returned or accepted by a
// DAF as an alternative to using RecCtx and the entity type separately.  This is most
// useful when there are multiple entity objects involved as inputs or outputs of a DAF.
// The type parameter T can either be a domain entity type or the pointer type thereof,
// depending on whether the DAF returns / receives by value or by pointer.
type Pw[T any] struct {
    RecCtx[T]
    Entity T
}

// Helper method
func (s Pw[T]) Copy(t T) Pw[T] {
    s.Entity = t
    return s
}

// Example B1 -- prefer the style of Example 1 above
func PersonUpdateDaf(pwPerson Pw[Person]) (Pw[Person], error)

// Example B2 -- usage where separating entity and RecCtx (see A2 above) would be 
// preferable
pwPerson, err := PersonReadDaf(name)
if err != nil { return err }
person := pwPerson.Entity
person = SomeBusinessFunctionBf(person)
pwPerson = pwPerson.Copy(person)
pwPerson, err = PersonUpdateDaf(pwPerson)

// Example B3 -- here it makes sense to use Pw
func ReadRecentUsagesDaf() ([]Pw[Usage], error)
```

Our philosophy is to define entity types that are relatively cheap to pass by value and to use value semantics for all
function parameters and method receivers. Thus, most of our functions do not mutate their receivers or parameters. This
makes for code that is easier to understand and where mutations are localised to explicit assignments. Not exactly pure
functional programming, but a practical idiom for an imperative language like Go. It also reduces garbage collection
pressure, often resulting in faster execution.

In most cases, this is achievable. In cases where the entity type has fields that can be large strings, we can either
use slice fields or `*string` fields as alternatives, and thus keep the cost of passing the entity by value low. In
cases where there is no practical way to design the entity so that passing it by value is cheap, we can fall back on
using pointers, e.g., a DAF can return a pointer and the type T in Pw[T] can be a pointer type.

## Alternative considered

Instead of using the above types and idioms, we could have DAFs return and consume a persistence wrapper *interface*
instead of a struct or separate entity and the database record context. This alternative approach looks like:

```go
type Pw[T any] interface {
    Entity() T // equivalent to the Entity field in the above struct
    Copy(t T) Pw[T] // creates a new instance with the same RecCtx and t as the T part
}
```

This alternative was tried and abandoned because:

- The resulting code at the point of usage is no cleaner than when using the struct version of Pw[T].
- The interface has to be implemented by some struct type S0 anyway and the code required to implement the required
  methods on S0 to implement the interface would be no easier than the code required to produce the struct version of
  Pw[T] from S0.
