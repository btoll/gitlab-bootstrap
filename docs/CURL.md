## cURL

## Groups

### List

```
$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/groups
```

### Create
```
$ curl -X POST --data '{"path": "gl-test", "name": "derp"}' --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" --header "Content-Type: application/json" https://gitlab.com/api/v4/groups
```

> This is not currently working, getting a 403 forbidden response.

### Transfer Project to Group

```
$ curl -X POST --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/groups/:id/projects/:project_id

$ curl -X POST --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/groups/45212731/projects/10924637
```

This is not currently working, getting a 403 forbidden response.

However, this does work:

https://docs.gitlab.com/ee/api/projects.html#transfer-a-project-to-a-new-namespace

```
$ curl -X PUT -d "namespace=gl-test-group2" -d "default_branch=master" --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/projects/:project_id/transfer

$ curl -X PUT -d "namespace=gl-test-group2" -d "default_branch=master" --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/projects/33860506/transfer
```

The Group namespace is found in the URL of the group: https://docs.gitlab.com/ee/user/group/#namespaces

## Projects

### List All

```
$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/projects
```

### List By User

```
$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/users/:user_id/projects

$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/users/10924637/projects
```

### Get Project ID By User

```
$ curl -X GET --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" https://gitlab.com/api/v4/users/10924637/projects 2> /dev/null | jq '.[0].id'
```

### Create Project with Template

```
$ curl -X POST --data '{"path": "herp", "name": "herp", "template_name": "hugo"}' --header "PRIVATE-TOKEN: $GITLAB_API_PRIVATE_TOKEN" --header "Content-Type: application/json" https://gitlab.com/api/v4/projects
```

