var app = angular.module("GoRunnerApp", ['ui.bootstrap', 'gorunnerServices', 'ngRoute'], function ($routeProvider) {
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
	$scope.refreshRuns = function() {
		$scope.recent = Run.query({'offset': 0, 'length': 20});
	};

	$scope.refreshRunsEvery = function(millis) {
		$scope.refreshRuns();
		$timeout(function(){
			$scope.refreshRunsEvery(millis);
		}, millis);
	};

	$scope.refreshRuns();
	$scope.refreshRunsEvery(3000);
});
