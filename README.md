The project run a endless worker with the simulation docker running workflow.

The workflow first pull the docker image from docker hub. Then, the workflow will create the container with randomly generated container id which means for the same image every run will create a separate container instance. After container spin up correctly, the workflow will trigger the simulation via running the end point by http://xxxx/start.


Steps to run this project:
1) You need a cadence service running. 
2) Run the following command
```
```bash
export DOCKERHUB_TOKEN=VNXN9nbywVdaBOkE
export DOCKERHUB_USERNAME=matrinos
export DOCKER_NETWORK=matrinos-network
./bin/main
```

TODO
