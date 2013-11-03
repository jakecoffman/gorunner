var app = angular.module("GoRunnerApp", [], function ($routeProvider) {
	$routeProvider.when('/jobs', {
		templateUrl: '/static/templates/jobs.html',
		controller: JobsCtl
	})
	.when('/jobs/:job', {
		templateUrl: '/static/templates/job.html',
		controller: JobCtl
	})
	.otherwise({
		redirectTo: '/jobs'
	});
});

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
		}
	}
});

app.controller('MainCtl', function ($scope, gorunner) {
	gorunner.getRecentRuns(function(data){
		$scope.recent = data;
	}, function(data) {
		$scope.recent = [];
	});
});

function JobsCtl($scope, gorunner) {
	gorunner.listJobs(function (data) {
		$scope.jobs = data;
	}, function () {
		alert("Error loading jobs");
	});
}

function JobCtl($scope, $routeParams, gorunner) {
	gorunner.getJob($routeParams.job, function(data){
		$scope.job = data;
	}, function(){
		alert("Error loading " + name);
	});
}
