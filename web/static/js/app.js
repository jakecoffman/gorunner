var app = angular.module("GoRunnerApp", ['ui.bootstrap', 'gorunnerServices'], function ($routeProvider) {
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
	.when('/triggers', {
		title: 'triggers',
		templateUrl: '/static/templates/triggers.html',
		controller: TriggersCtl
	})
	.when('/triggers/:trigger', {
		title: 'trigger',
		templateUrl: '/static/templates/trigger.html',
		controller: TriggerCtl
	})
	.when('/runs', {
		title: 'runs',
		templateUrl: '/static/templates/runs.html',
		controller: RunsCtl
	})
	.when('/runs/:run', {
		title: 'run',
		templateUrl: '/static/templates/run.html',
		controller: RunCtl
	})
	.otherwise({
		redirectTo: '/jobs'
	});
});

app.filter('join', function(){
	return function(input) {
		if(input)
			return input.join(', ');
		else
			return "";
	};
});

app.run(['$location', '$rootScope', function($location, $rootScope) {
	$rootScope.$on('$routeChangeSuccess', function (event, current, previous) {
		if(current.$$route) {
			$rootScope.title = current.$$route.title;
		}
	});
}]);

app.controller('MainCtl', function ($scope, $timeout, Run) {
	$scope.recent = [];
	$scope.refreshRuns = function() {
		$scope.recent = Run.query({'offset': 0, 'length': 20});
	};

	$scope.refreshRunsEvery = function(millis) {
		$scope.refreshRuns();
		$timeout(function(){
			$scope.refreshRunsEvery(millis);
		}, millis);
	};

	$scope.refreshRuns();
	$scope.refreshRunsEvery(3000);
});

function JobsCtl($scope, Job, Run) {
	$scope.jobs = Job.query();

	$scope.quickRun = function(job) {
		var run = new Run();
		run.job = job;
		run.$save();
		$scope.refreshRuns();
	};

	$scope.promptJob = function() {
		var name = prompt("Enter name of job:");
		if(name) {
			var job = new Job();
			job.name = name;
			job.$save();
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

function RunsCtl($scope, Run) {
	$scope.runs = Run.query();
}

function RunCtl($scope, $routeParams, Run) {
	$scope.run = Run.get({id: $routeParams.run});
}

function TriggersCtl($scope, Trigger) {
	$scope.triggers = Trigger.query();

	$scope.deleteTrigger = function(name) {
		Trigger.$delete({id: name});
		$scope.refreshTriggers();
	};

	$scope.promptTrigger = function(){
		var name = prompt("Enter a name for the new trigger");
		if(name) {
			var trigger = new Trigger();
			trigger.name = name;
			trigger.$save();
			$scope.refreshTriggers();
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