gorunner
========

[![Build Status](https://secure.travis-ci.org/jakecoffman/gorunner.png?branch=master)](http://travis-ci.org/jakecoffman/gorunner)

gorunner is an attempt to create a continuous integration web server written in Golang.

This project is a work-in-progress but development is not very active. I accept pull requests but also if you want to take it in a different direction let me know and we can collaborate.

Installation instructions
----

Assuming $GOPATH/bin is on your path:

	go get github.com/jakecoffman/gorunner
	cd $GOPATH/src/github.com/jakecoffman/gorunner
	gorunner

Technologies
----

* Go (golang)
* Javascript
  * Angularjs
  * Websockets

Why Go?
----

Go's ability to handle many connections would be beneficial for:

* running multiple build scripts and monitoring progress
* connecting to a cluster of gorunner servers
* live updates to builds in the UI via websockets, etc

![gorunner](https://raw.githubusercontent.com/jakecoffman/gorunner/master/promo.png "gorunner")
