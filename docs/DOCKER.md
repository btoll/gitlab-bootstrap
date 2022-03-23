## Build

```
$ docker build -t gitlab-client:test .
```

## Run

### Help

```
$ docker run \
    --rm \
    --init gitlab-client:test \
    -h
Usage of ./gitlab-client:
  -destroy
        Should destroy all projects listed in the given Gitlab config file.
  -file string
        Path to GitLab config file (json or yaml).
  -user string
        List everything for the given user.
```

### Create

```
$ docker run \
    --rm \
    --init \
    -e GITLAB_API_PRIVATE_TOKEN="$GITLAB_API_PRIVATE_TOKEN" \
    -v $(pwd)/examples:/build/examples \
    gitlab-client:test \
    -file examples/gitlab.yaml
```

### Destroy

```
$ docker run \
    --rm \
    --init \
    -e GITLAB_API_PRIVATE_TOKEN="$GITLAB_API_PRIVATE_TOKEN" \
    -v $(pwd)/examples:/build/examples \
    gitlab-client:test \
    -file examples/gitlab.yaml \
    -destroy
```

