function JobsCtl($scope, Job, Run) {
	$scope.jobs = Job.query();
	$scope.selected = [];
	$scope.gridOptions = {
		data: 'jobs',
		plugins: [new ngGridFlexibleHeightPlugin()],
		multiSelect: false,
		selectedItems: $scope.selected,
		columnDefs: [
			{field: 'name', displayName: 'Name'},
			{field: 'status', displayName: 'Status'},
			{field: 'tasks', displayName: 'Tasks', cellTemplate: '/static/gridTemplates/count.html'},
			{field: 'triggers', displayName: 'Triggers', cellTemplate: '/static/gridTemplates/count.html'}
		]
	};

	$scope.quickRun = function(name) {
		console.log(name);
		var run = new Run();
		run.job = name;
		run.$save();
	};

	$scope.promptJob = function() {
		var name = prompt("Enter name of job:");
		if(name) {
			Job.addJob({id: name});
			window.location = "/#/jobs/" + name;
		}
	}
}

function JobCtl($scope, $routeParams, Job, Task, Trigger) {
	$scope.refreshJob = function(){
		$scope.job = Job.get({id: $routeParams.job})
	};

	$scope.refreshTasks = function(){
		$scope.tasks = Task.query();
	};

	$scope.refreshTriggers = function() {
		$scope.triggers = Trigger.query();
	};

	$scope.removeTask = function(idx) {
		Job.removeTask({tidx: idx, id: $routeParams.job});
		$scope.refreshJob();
	};

	$scope.removeTrigger = function(name) {
		Job.removeTrigger({trigger: name, id: $routeParams.job});
		$scope.refreshJob()
	};

	$scope.addTaskToJob = function(task) {
		Job.addTask({task: task, id: $routeParams.job});
		$scope.refreshJob();
	};

	$scope.addTriggerToJob = function(trigger) {
		Job.addTrigger({trigger: trigger, id: $routeParams.job});
		$scope.refreshJob();
	};

	$scope.deleteJob = function() {
		Job.$delete({id: $routeParams.job});
		window.location = "/#/jobs";
	};

	$scope.refreshJob();
	$scope.refreshTasks();
	$scope.refreshTriggers();
}
