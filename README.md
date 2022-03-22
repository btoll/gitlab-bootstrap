# GitLab Client

## Requirements

- [GitLab API token] with `api` scope
- A pre-existing group
    + The current API has [problems creating new groups]

## Notes

- To keep it simple, the given project's `name` will also be used as its `path`.
    + https://docs.gitlab.com/ee/api/index.html#namespaced-path-encoding

- All projects are created from a default template.
    + https://gitlab.com/gitlab-org/project-templates

## Examples

### Creating Projects

```
$ go run main.go group.go user.go project.go -file examples/gitlab.yaml
```

### Deleting Projects

To teardown what was setup when creating the projects, simply pass the same config file with the `-destroy` flag.

Or, pass another file or make your changes in the same one.  Pick your poison.

```
$ go run main.go group.go user.go project.go -file examples/gitlab.yaml -destroy
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
    - name: bar
      tpl_name: android
      visibility: public
    - name: quux
      tpl_name: dotnetcore
      visibility: public
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

[GitLab API token]: https://docs.gitlab.com/ee/security/token_overview.html
[problems creating new groups]: https://gitlab.com/gitlab-org/gitlab/-/issues/244345
[`go-gitlab`]: https://github.com/xanzy/go-gitlab

