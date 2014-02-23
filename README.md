gorunner
========

gorunner is an attempt to create a continuous integration web server written in golang.
As development occurs, I hope reusable code and interfaces develop since go is a
relatively new language.

Installation instructions
----

	$ cd $GOPATH/src
	$ go get github.com/jakecoffman/gorunner
	$ cd github.com/jakecoffman/gorunner
	$ go run main.go
	
Or, you know, whatever works. 

Screenshots
----

Jobs are the bucket of things, such as triggers and tasks, that tell the system
when to execute and what to do.
![Jobs page](http://www.coffshire.com/static/gorunner/jobs.png "Jobs page")

Tasks are editable in the browser.
![Jobs page](http://www.coffshire.com/static/gorunner/task.png "A task")

When a job executes, it creates a run which records the output of the tasks. 
![Runs page](http://www.coffshire.com/static/gorunner/runs.png "Runs page")

Technologies
----

* Go (golang)
* Javascript
  * Angularjs
  * Websockets
