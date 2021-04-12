# go-foa-realworld

This is a [RealWorld Example App](https://github.com/gothinkster/realworld) backend written in Golang in a
function-oriented architecture style.

Some of the code was based on https://github.com/err0r500/go-realworld-clean.

## To-dos:

v Complete first pass at SFLs.
v Review use of pointers, clean it up.
* Review model in terms of relationships and how they should be represented.
* Review model and move functionality to BFs as appropriate.
* Review FLs and SFLs, move functionality to BFs, model, or rpc as appropriate.
* Review basic error handling, including fs/errors.go.
* Create request adapter pattern that includes:
    - last-resort error handling
    - encompasses both query parameters and JSON payloads
* Create pattern for platform-independent unmarshallers.
* Create platform-independent unmarshallers for rpc structs.
* Implement authentication.
* Guerrilla-test some SFLs.
* Implement test suite.
* Implement DAFs with SQLite. The current in-memory persistence based on the implementation in
  https://github.com/err0r500/go-realworld-clean is inherently broken as it doesn't support optimistic concurrency and
  allows dirty writes.
* Rinse and repeat.
