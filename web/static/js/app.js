var app = angular.module("GoRunnerApp", [], function ($routeProvider) {
	$routeProvider.when('/jobs', {
		templateUrl: '/static/templates/jobs.html',
		controller: JobsCtl
	});
	$routeProvider.otherwise({
		redirectTo: '/jobs'
	});
});

app.controller('MainCtl', function ($scope, $http) {
	$scope.getRecentRuns = function () {
		$http({
			method: "GET",
			url: "/runs"
		})
			.success(function (data) {
				$scope.recent = data;
			})
			.error(function () {
				$scope.recent = [];
			})
	};
	$scope.getRecentRuns();
});

function JobsCtl($scope, $http) {
	$scope.listJobs = function () {
		$http({
			method: "GET",
			url: "/jobs"
		})
			.success(function (data) {
				$scope.jobs = data;
			})
			.error(function () {
				alert("Error loading jobs");
			});
	};

	$scope.listJobs();
}