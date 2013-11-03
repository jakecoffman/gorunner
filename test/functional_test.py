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
        return [job['Name'] for job in jobs]

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

    def _raise_if_status_not(self, r, status):
        if r.status_code != status:
            raise Exception(r.text)

class TestGoAPI(unittest.TestCase):
    def setUp(self):
        self.api = GoRunnerAPI("http://localhost:8090")

    def tearDown(self):
        pass

    def test_API(self):
        test_job = "test999"

        names = self.api.list_job_names()
        if test_job in names:
            self.api.delete_job(test_job)
            names = self.api.list_job_names()
        self.assertNotIn(test_job, names)

        self.api.add_job(test_job)
        names = self.api.list_job_names()
        self.assertIn(test_job, names)

        # job = self.api.get_job(test_job)
        # self.assertEqual(test_job, job['Name'])

        self.api.delete_job(test_job)
        names = self.api.list_job_names()
        self.assertNotIn(test_job, names)

if __name__ == "__main__":
    unittest.main()
