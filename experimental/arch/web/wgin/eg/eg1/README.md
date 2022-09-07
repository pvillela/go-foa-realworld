This exemplifies the use of an adapter to transform a servicce flow into a Gin HandlerFunc.

This example was used as a precursor to eg2.  It was used to generate a valid JWT token that is used in eg2.

The service flow is a simplistic login function, but it is used here as a regular service flow, not a login function.

The dummyAuthenticator function always authenticates the request.  One of its purposes is to log a valid JWT token that is used in eg2.

The test examples focus on exercising the error handler.
