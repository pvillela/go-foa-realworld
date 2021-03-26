# go-foa-realworld

This is a [RealWorld Example App](https://github.com/gothinkster/realworld) backend written in Golang in a
function-oriented architecture style.

## To-dos:

* Complete first pass at SFLs.
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
* Implement DAFs with SQLite.
* Rinse and repeat.
