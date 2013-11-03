import requests
import json
import unittest


class GoRunnerAPI(object):
    def __init__(self, host):
        self.host = host

    def list_jobs(self):
        r = requests.get("%s/jobs" % self.host)
        self._raise_if_status_not(r, 200)
        return r.json()

    def list_job_names(self):
        jobs = self.list_jobs()
        return [job['name'] for job in jobs]

    def add_job(self, name):
        r = requests.post("%s/jobs" % self.host, data=json.dumps({'name': name}))
        self._raise_if_status_not(r, 201)

    def get_job(self, name):
        r = requests.get("%s/jobs/%s" % (self.host, name))
        self._raise_if_status_not(r, 200)
        return r.json()

    def delete_job(self, name):
        r = requests.delete("%s/jobs/%s" % (self.host, name))
        self._raise_if_status_not(r, 200)

    def add_task_to_job(self, task, job):
        r = requests.post("%s/jobs/%s/tasks" % (self.host, job), json.dumps({'task': task}))
        self._raise_if_status_not(r, 201)

    def list_tasks(self):
        r = requests.get("%s/tasks" % self.host)
        self._raise_if_status_not(r, 200)
        return r.json()

    def list_task_names(self):
        tasks = self.list_tasks()
        return [task['name'] for task in tasks]

    def add_task(self, name):
        r = requests.post("%s/tasks" % self.host, json.dumps({'name': name}))
        self._raise_if_status_not(r, 201)

    def get_task(self, name):
        r = requests.get("%s/tasks/%s" % (self.host, name))
        self._raise_if_status_not(r, 200)
        return r.json()

    def delete_task(self, name):
        r = requests.delete("%s/tasks/%s" % (self.host, name))
        self._raise_if_status_not(r, 200)

    def list_runs(self):
        r = requests.get("%s/runs" % self.host)
        self._raise_if_status_not(r, 200)
        return r.json()

    def list_run_ids(self):
        runs = self.list_runs()
        return [run['uuid'] for run in runs]

    def run_job(self, name):
        r = requests.post("%s/runs" % self.host, json.dumps({'job': name}))
        self._raise_if_status_not(r, 201)
        return r.json()

    def _raise_if_status_not(self, r, status):
        if r.status_code != status:
            raise Exception(r.text)


class TestGoAPI(unittest.TestCase):
    def setUp(self):
        self.api = GoRunnerAPI("http://localhost:8090")

        self.test_job = "test_job999"
        self.test_task = "test_task999"

    def tearDown(self):
        self.api.delete_job(self.test_job)
        self.api.delete_task(self.test_task)

    def test_jobs(self):
        self.crud_test(self.api.list_job_names, self.api.delete_job, self.api.add_job, self.api.get_job)

    def test_tasks(self):
        self.crud_test(self.api.list_task_names, self.api.delete_task, self.api.add_task, self.api.get_task)

    def test_adding_job_with_no_name(self):
        try:
            self.api.add_job("")
            self.fail()
        except Exception:
            pass

    def test_adding_job_with_no_payload(self):
        try:
            requests.post("%s/jobs" % self.api.host)
            self.fail()
        except Exception:
            pass

    def test_add_task_to_job(self):
        self.api.add_job(self.test_job)
        self.api.add_task(self.test_task)
        self.api.add_task_to_job(self.test_task, self.test_job)

        job = self.api.get_job(self.test_job)
        self.assertIn(self.test_task, job['tasks'])

        self.api.delete_job(self.test_job)
        self.api.delete_task(self.test_task)

    def test_runs(self):
        self.api.add_job(self.test_job)
        self.api.add_task(self.test_task)
        self.api.add_task_to_job(self.test_task, self.test_job)
        uuid = self.api.run_job(self.test_job)['uuid']
        runs = self.api.list_run_ids()
        self.assertIn(uuid, runs)

        print json.dumps(self.api.list_runs(), indent=4)

    def crud_test(self, list_names, delete, add, get):
        test_name = "test999"

        names = list_names()
        if test_name in names:
            delete(test_name)
            names = list_names()
        self.assertNotIn(test_name, names)

        add(test_name)
        names = list_names()
        self.assertIn(test_name, names)

        thing = get(test_name)
        self.assertEqual(test_name, thing['name'])

        delete(test_name)
        names = list_names()
        self.assertNotIn(test_name, names)

if __name__ == "__main__":
    unittest.main()
