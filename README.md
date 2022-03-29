# GitLab Client

## What does this do?

The goal of this project is a simple one: quickly and easily create as many `GitLab` projects as needed within an **existing** group.

## What does this support?

This project **only** supports the *creation* of the following `GitLab` objects:

- [Projects]
- [Invites]
- [Issues]
- [Releases]

Its intent is to **only** create them with **only** the required configuration parameters.  If you need something more customized, you'll have to do that yourself.

## Requirements

- [`GitLab` API token] with `api` scope
- A pre-existing group
    + The current API has [problems creating new groups]

## Notes

- To keep it simple, the given project's `name` will also be used as its `path`.
    + https://docs.gitlab.com/ee/api/index.html#namespaced-path-encoding

- All projects are created from a default template.
    + https://gitlab.com/gitlab-org/project-templates

- Sending invites will automatically add the invitee as a member to the project, if they have already created a `GitLab` account.  Otherwise, the invite will be pending.

## Fields

### Creating [Projects]

A list of `Project` objects composed of:

- `name` (string)
- `tpl_name` (string)
- `visibility` (string)
- [`invites`](#invites) (list of `Invites`)
- [`issues`](#issues) (list of `Issues`)
- [`releases`](#releases) (list of `Releases`)

### [`invites`]

A list of `Invite` objects composed of:

- `access_level` (string)
    + These values mostly map directly to the [Members API values].
        - `None`
        - `Minimal`
        - `Guest`
        - `Reporter`
        - `Developer` (default)
        - `Maintainer`
        - `Owner`
- `email` (string)

### [`issues`]

A list of `Issue` objects composed of:

- `title` (string)
- `type` (string)
    + These values map directly to the `GitLab` API values.
        - `Incident`
        - `Issue` (default)
        - `TestCase`

### [`releases`]

A list of `Release` objects composed of:

- `name` (string)
- `ref` (string)
- `tag_name` (string)

> For full examples in both `yaml` and `json`, see the `examples/` directory.

## Examples

### Creating Projects

```
$ gitlab-client -file examples/gitlab.yaml
$ gitlab-client -file examples/gitlab.json
```

### Deleting Projects

To teardown what was setup when creating the projects, simply pass the same config file with the `-destroy` flag.

Or, pass another file or make your changes in the same one.  Pick your poison.

```
$ gitlab-client -file examples/gitlab.yaml -destroy
```

> This will **not** ask for confirmation.

## Config Example

The tool expects an array of `group` objects.  Each `group` object consists of one or more `project`s.

`gitlab.yaml`

```
---
- name: gl-group
  projects:
    - name: foo
      tpl_name: hugo
      visibility: public
      invites:
        - email: foobar@example.com
          access_level: Developer
      issues:
        - title: yo
          type: TestCase
        - title: humdinger
          type: Incident
      releases:
        - name: test1
          ref: master
          tag_name: test1.0
        - name: test2
          ref: master
          tag_name: test2.0
    - name: bar
      tpl_name: android
      visibility: public
    - name: quux
      tpl_name: dotnetcore
      visibility: public
      invites:
        - email: btoll@example.com
          access_level: Guest
        - email: kilgore-trout@example.com
          access_level: Maintainer
---
```

> In addition to `yaml`, `json` is supported.

## Testing

If you don't want to compile, you can use `go run`:

```
$ go run main.go group.go project.go user.go -file examples/gitlab.yaml
$ go run main.go group.go project.go user.go -user btoll
```

## Acknowledgments

This project uses the [`go-gitlab`] client library.

[Projects]: https://docs.gitlab.com/ee/api/projects.html
[Invites]: https://docs.gitlab.com/ee/api/invitations.html
[`invites`]: https://docs.gitlab.com/ee/api/invitations.html
[Issues]: https://docs.gitlab.com/ee/api/issues.html
[`issues`]: https://docs.gitlab.com/ee/api/issues.html
[Releases]: https://docs.gitlab.com/ee/api/releases/
[`releases`]: https://docs.gitlab.com/ee/api/releases/
[`GitLab` API token]: https://docs.gitlab.com/ee/security/token_overview.html
[problems creating new groups]: https://gitlab.com/gitlab-org/gitlab/-/issues/244345
[Members API values]: https://docs.gitlab.com/ee/development/permissions.html#members
[`go-gitlab`]: https://github.com/xanzy/go-gitlab

