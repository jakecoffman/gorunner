var app = angular.module("GoRunnerApp", [], function ($routeProvider) {
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
	.otherwise({
		redirectTo: '/jobs'
	});
});

app.run(['$location', '$rootScope', function($location, $rootScope) {
    $rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
        $rootScope.title = current.$$route.title;
    });
}]);

app.factory('gorunner', function($http){
	return {
		getRecentRuns: function (success, failure) {
			$http({
				method: "GET",
				url: "/runs"
			})
			.success(success)
			.error(failure)
		},

		listJobs: function (success, failure) {
			$http({
				method: "GET",
				url: "/jobs"
			})
			.success(success)
			.error(failure);
		},

		getJob: function(name, success, failure) {
			$http({
				method: "GET",
				url: "/jobs/" + name
			})
			.success(success)
			.error(failure)
		},

		listTasks: function(success, failure) {
			$http({
				method: "GET",
				url: "/tasks"
			})
			.success(success)
			.error(failure);
		},

		getTask: function(name, success, failure) {
			$http({
				method: "GET",
				url: "/tasks/" + name
			})
			.success(success)
			.error(failure);
		}
	}
});

app.controller('MainCtl', function ($scope, gorunner) {
	$scope.page = "jobs";
	gorunner.getRecentRuns(function(data){
		$scope.recent = data;
	}, function(data) {
		$scope.recent = [];
	});
});

function JobsCtl($scope, gorunner) {
	$scope.page = "jobs";
	gorunner.listJobs(function (data) {
		$scope.jobs = data;
	}, function () {
		alert("Error loading jobs");
	});
}

function JobCtl($scope, $routeParams, gorunner) {
	$scope.page = "job";
	gorunner.getJob($routeParams.job, function(data){
		$scope.job = data;
	}, function(){
		alert("Error loading " + name);
	});
}

function TasksCtl($scope, gorunner) {
	$scope.page = "tasks";
	gorunner.listTasks(function(data) {
		$scope.tasks = data;
	}, function(data) {
		alert("Error loading tasks");
	})
}

function TaskCtl($scope, $routeParams, gorunner) {
	$scope.page = "task";
	gorunner.getTask($routeParams.task, function(data) {
		$scope.task = data;
	}, function(data) {
		alert("Error loading task " + $routeParams.task)
	})
}
