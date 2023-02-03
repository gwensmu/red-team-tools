# building 

`go build -o bin/$(basename $(pwd))`


### todo


https://cybervelia.com/?p=736#htoc-query-for-arguments

* introspection query
* query for query args
* query for mutation args

* compare this to fields called by the UI - the diff is interesting fruit
* intercept and enumerate queries called by the ui

'''
Therefore, the developers often think that the input data are the most important part of the request, as it contains values defined by the user. But because the queries they receive during the development cycles are predefined, because of the front-end implementation, they often forget they should also check for access controls in the fields of the output of the query.
'''

Using amplification:
* scan for nosql injection
* scan for sql injection
* scan for stored xss
* use pagination where possible
* scan for graphiql
