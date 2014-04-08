function TriggersCtl($scope, Trigger) {
	$scope.triggers = Trigger.query();
	$scope.selected = [];
	$scope.gridOptions = {
		data: 'triggers',
		plugins: [new ngGridFlexibleHeightPlugin()],
		multiSelect: false,
		selectedItems: $scope.selected
	};

	$scope.deleteTrigger = function(name) {
		Trigger.delete({id: name});
		$scope.triggers = Trigger.query();
	};

	$scope.promptTrigger = function(){
		var name = prompt("Enter a name for the new trigger");
		if(name) {
			var trigger = new Trigger();
			trigger.name = name;
			trigger.$save();
			$scope.triggers = Trigger.query();
		}
	}
}

function TriggerCtl($scope, $routeParams, Trigger) {
	$scope.trigger = Trigger.get({id: $routeParams.trigger});

	$scope.saveTrigger = function() {
		Trigger.update({id: $scope.trigger.name, cron: $scope.trigger.schedule});
		window.location = "/#/triggers";
	};

	$scope.jobs = Trigger.listJobs({id: $routeParams.trigger});
}
