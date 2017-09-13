package summary_test

const jobsPayload = `[
  {
    "id": 1,
    "name": "testJob1",
    "url": "/test1.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "started"
    }
  },
  {
    "id": 2,
    "name": "testJob2",
    "url": "/test2.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "succeeded"
    }
  },
  {
    "id": 3,
    "name": "testJob3",
    "url": "/test3.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "failed"
    }
  },
  {
    "id": 4,
    "name": "testJob4",
    "url": "/test4.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "errored"
    }
  },
  {
    "id": 5,
    "name": "testJob5",
    "url": "/test5.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "aborted"
    }
  },
  {
    "id": 6,
    "name": "testJob6",
    "url": "/test6.job.url",
    "paused": false,
    "team_name": "main",
    "finished_build": {
      "id": 1,
      "status": "pending"
    }
  }
]`

const pipelinesPayload = `[
  {
    "id": 1,
    "name": "test1",
    "url": "/test1.url",
    "paused": false,
    "public": true,
    "team_name": "main"
  }
]`
