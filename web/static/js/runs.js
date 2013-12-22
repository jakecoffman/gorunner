app.filter('datetime', function(){
	return function(input){
		var dt = new Date(Date.parse(input));
		return dt.toLocaleString();
	}
});

function RunsCtl($scope, Run) {
	$scope.runs = Run.query();
	$scope.selected = [];

	$scope.blah = function(data) {
		console.log(data);
		return data;
	};

	$scope.gridOptions = {
		data: 'runs',
		plugins: [new ngGridFlexibleHeightPlugin()],
		multiSelect: false,
		selectedItems: $scope.selected,
		columnDefs: [
			{field: 'uuid', displayName: 'UUID'},
			{field: 'job.name', displayName: 'Job'},
			{field: 'tasks', displayName: 'Tasks', cellTemplate: '/static/gridTemplates/count.html'},
			{field: 'status', displayName: 'Status'},
			{field: 'start', displayName: 'Start', cellFilter: 'datetime'},
			{field: 'end', displayName: 'End', cellFilter: 'datetime'}
		]
	};
}

function RunCtl($scope, $routeParams, $timeout, Run) {
	var update = function() {
		Run.get({id: $routeParams.run}, function(data) {
			$scope.run = data;
			if (data.status == "Running") {
				$timeout(update, 3000);
			}
		});
	};
	$scope.run |= {};
	update();
	
}
