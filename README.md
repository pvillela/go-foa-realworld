# go-foa-realworld

This is a [RealWorld Example App](https://github.com/gothinkster/realworld) backend written in Golang in a
function-oriented architecture style.

Some of the code was based on https://github.com/err0r500/go-realworld-clean.

## To-dos:

v Create web handlers for non-POST methods
- Test handler adapter
- Create Default ReqCtxExtractor
- Implement correct request context logic with Context. Initially, ReqCtx is only used as an intermediary to get username and is not being stored in Context.
v Complete first pass at SFLs. 
v Review use of pointers, clean it up. 
v Review model in terms of relationships and how they should be represented.
v Redo GetBy... functions to use FindFirst
v Complete review and clean-up of model
  x Embed Comment in Article; no need for separate Comment DAFs with mapdb
  v No need for this for mapdb, but we want to do it for consistency with how we might model it
    with a real database (relational or NoSQL): Separate Favorited entity with its own DAFs
  v Add followers count field to User
* Review FLs and SFLs, move functionality to BFs, model, or rpc as appropriate.
    - Start with article filtering logic
v Review model and move functionality to BFs and DAFs as appropriate.
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
