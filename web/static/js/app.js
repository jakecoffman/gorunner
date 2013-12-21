var app = angular.module("GoRunnerApp", ['ui.bootstrap', 'gorunnerServices', 'ngRoute', 'ui.ace'], function ($routeProvider) {
	$routeProvider.when('/jobs', {
		title: "jobs",
		templateUrl: '/static/templates/jobs.html',
		controller: JobsCtl
	})
	.when('/jobs/:job', {
		title: "job",
		templateUrl: '/static/templates/job.html',
		controller: JobCtl
	})
	.when('/tasks', {
		title: "tasks",
		templateUrl: '/static/templates/tasks.html',
		controller: TasksCtl
	})
	.when('/tasks/:task', {
		title: "task",
		templateUrl: '/static/templates/task.html',
		controller: TaskCtl
	})
	.when('/triggers', {
		title: 'triggers',
		templateUrl: '/static/templates/triggers.html',
		controller: TriggersCtl
	})
	.when('/triggers/:trigger', {
		title: 'trigger',
		templateUrl: '/static/templates/trigger.html',
		controller: TriggerCtl
	})
	.when('/runs', {
		title: 'runs',
		templateUrl: '/static/templates/runs.html',
		controller: RunsCtl
	})
	.when('/runs/:run', {
		title: 'run',
		templateUrl: '/static/templates/run.html',
		controller: RunCtl
	})
	.otherwise({
		redirectTo: '/jobs'
	});
});

app.filter('join', function(){
	return function(input) {
		if(input)
			return input.join(', ');
		else
			return "";
	};
});

app.run(['$location', '$rootScope', function($location, $rootScope) {
	$rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
		if(current.$$route) {
			$rootScope.title = current.$$route.title;
		}
	});
}]);

app.controller('MainCtl', function ($scope, $timeout, Run) {
	$scope.recent = [];

	var conn = new WebSocket("ws://localhost:8090/ws");
	conn.onclose = function(e) {
		console.log("Connection closed");
	};

	conn.onopen = function(e) {
		console.log("Conencted");
	};

	conn.onmessage = function(e){
		console.log("Got data");
		$scope.$apply(function(){
			$scope.recent = JSON.parse(e.data);
		});
	}

});
