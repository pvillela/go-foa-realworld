# go-foa-realworld

This is a [RealWorld Example App](https://github.com/gothinkster/realworld) backend written in Golang in a
function-oriented architecture style.

## To-dos:

* Review model structs, modify them as appropriate, create any necessary corresponding rpc structs.
* For each Sfl, add a CoreFl that does not use the external input/output structs. The Sfl calls BFs to convert to/from
  the input/output structs. At least for now, put Sfl and CoreFl in the same file.
* Replace Make higher-order methods with straight methods and take advantage of partial application. When needed, use a
  Prep function to create an Aug struct that is the receiver of the straight method.
* Consider defining an ArticleDafs interface that bundles all the DAFs for Article; similarly for other entities.
