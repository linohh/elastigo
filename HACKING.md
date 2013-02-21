
Testing
-----------------

To run tests, this library loads data into an elasticsearch server as a one-time load.  Then unit tests run against that.

See core/test_test.go.   The data set should remain the same as it pulls a known set of github archive data.

usage:
	
	# one-time load of test data
	cd core
    go test -v -host myelasticsearch.domain  -loaddata 
    
    # run unit tests for core
    go test -v -host myelasticsearch.domain

    cd ../search
    go test -v -host myelasticsearch.domain 

Clean out the Elasticsearch index:
    	
    # using https://github.com/jkbr/httpie 
    http -v DELETE http://localhost:9200/github