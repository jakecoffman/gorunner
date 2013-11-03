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
	var s = {};

	s.getRecentRuns = function (scope) {
		$http({
			method: "GET",
			url: "/runs"
		})
		.success(function (data) {
			scope.recent = data;
		})
		.error(function () {
			scope.recent = [];
		})
	};

	s.listJobs = function (scope) {
		$http({
			method: "GET",
			url: "/jobs"
		})
		.success(function (data) {
			scope.jobs = data;
		})
		.error(function () {
			alert("Error loading jobs");
		});
	};

	s.getJob = function(name, scope) {
		$http({
			method: "GET",
			url: "/jobs/" + name
		})
		.success(function(data){
			scope.job = data;
		})
		.error(function(){
			alert("Error loading " + name);
		})
	};

	return s;
});

app.controller('MainCtl', function ($scope, gorunner) {
	gorunner.getRecentRuns($scope);
});

function JobsCtl($scope, gorunner) {
	gorunner.listJobs($scope);
}

function JobCtl($scope, $routeParams, gorunner) {
	gorunner.getJob($routeParams.job, $scope);
}