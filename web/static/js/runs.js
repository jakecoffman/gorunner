function RunsCtl($scope, Run) {
	$scope.runs = Run.query();
}

function RunCtl($scope, $routeParams, $timeout, Run) {
	var update = function() {
		Run.get({id: $routeParams.run}, function(data) {
			$scope.run = data
			if (data.status == "Running") {
				$timeout(update, 3000);
			}
		});
	}
	$scope.run |= {};
	update();
	
}
