function RunsCtl($scope, Run) {
	$scope.runs = Run.query();
}

function RunCtl($scope, $routeParams, Run) {
	$scope.run = Run.get({id: $routeParams.run});
}
