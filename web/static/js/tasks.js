
function TasksCtl($scope, Task) {
	$scope.tasks = Task.query();

	$scope.promptTask = function() {
		var name = prompt("Enter name of task:");
		if(name) {
			var newTask = new Task();
			newTask.name = name;
			newTask.$save();
			window.location = "/#/tasks/" + name;
		}
	};
}

function TaskCtl($scope, $routeParams, Task) {
	$scope.task = Task.get({id: $routeParams.task});
	$scope.jobs = Task.jobs({id: $routeParams.task});

	$scope.saveTask = function() {
		Task.update({id: $routeParams.task, script: $scope.task.script});
		window.location = "/#/tasks";
	};

	$scope.deleteTask = function() {
		Task.$delete({id: $routeParams.task});
	}
}
