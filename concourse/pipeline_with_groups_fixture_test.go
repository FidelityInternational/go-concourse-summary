package summary_test

const examplePipeline = `[
    {
        "id": 2,
        "name": "cf-example-pipeline",
        "paused": false,
        "public": true,
        "groups": [
            {
                "name": "test-group",
                "jobs": [
                    "testJob1",
                    "testJob2"
                ]
            }
        ],
        "team_name": "main"
    }
]`

const examplePipelineJobs = `[
    {
        "id": 695,
        "name": "testJob1",
        "pipeline_name": "cf-example-pipeline",
        "team_name": "main",
        "next_build": null,
        "finished_build": {
            "id": 691532,
            "team_name": "main",
            "name": "4",
            "status": "succeeded",
            "job_name": "testJob1",
            "api_url": "/api/v1/builds/691532",
            "pipeline_name": "cf-example-pipeline",
            "start_time": 1525965370,
            "end_time": 1525965398
        },
        "inputs": [
            {
                "name": "testInput1",
                "resource": "testResource1",
                "trigger": false
            },
            {
                "name": "testInput2",
                "resource": "testResoure2",
                "trigger": false
            }
        ],
        "outputs": [],
        "groups": [
            "test-group"
        ]
    },
    {
        "id": 696,
        "name": "testJob1",
        "pipeline_name": "cf-example-pipeline",
        "team_name": "main",
        "next_build": null,
        "finished_build": {
            "id": 691297,
            "team_name": "main",
            "name": "5",
            "status": "succeeded",
            "job_name": "testJob1",
            "api_url": "/api/v1/builds/691297",
            "pipeline_name": "cf-example-pipeline",
            "start_time": 1525965125,
            "end_time": 1525965137
        },
        "inputs": [
            {
                "name": "testInput1",
                "resource": "testResource1",
                "trigger": false
            },
            {
                "name": "testInput2",
                "resource": "testResource2",
                "trigger": false
            }
        ],
        "outputs": [],
        "groups": [
            "test-group"
        ]
    }
]`
