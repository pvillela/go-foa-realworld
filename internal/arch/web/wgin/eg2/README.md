This exemplifies the use of an adapter to transform a servicce flow into a Gin HandlerFunc.

The service flow is a simplistic login function, but it is used here as a regular service flow, not a login function.

The dummyAuthenticator function just checks if the JWT token in the HTTP header is valid.

The test examples send valid or invalid authentication tokens as headers.
