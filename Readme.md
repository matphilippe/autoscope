# autoscope

`autoscope` automates the inclusion of **scopes** in your [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

1. You define rules such that "all files under directory `my-dir/` are in the module `myModule`".
2. If a commit changes a file within `myModule`, then the scope `myModule` will be appended to your commit scope.

## But Why ?

Conventional commits are awesome: they allow you to leverage your git logs to support your workflow
However, they have drawbacks: a typo could lead to very different interpretation of a commit's message.

`autoscope` helps with typos and omissions around the commit message's scope.

Some teams/users may not even care about scopes, but in the context of `monorepos`, they can play a pivotal role.

Consider a repo structured as follows:

```text
README.md
docs/
modules/
  A/
  B/
```

where `A`, and `B` are subprojects you'd like to release and version independently.

In a PR, you make a `fix` to `A`, and a `feat!` to be, expecting a patch bump and major bump to `A` and `B` respectively.

- If you merge the PR by rebasing, your log will show the `fix` commit on `A`, and the `feat!` on `B`. Come back tomorrow, you'll still understand what happened there.
- If you merge the PR with a squash then rebase, you will lose information: a commit now indicates both a `fix` and a `feat!` and touches `A` and `B`. What do you bump ?

Scopes are useful, but it puts work and focus on the dev, this tool lifts some of the work.

## Installation

```bash
go install github.com/matphilippe/autoscope@latest
```

Or build from source:

```bash
git clone https://github.com/matphilippe/autoscope
cd autoscope
go build
```

## Configuration

Drop a yaml file named `.autoscope.yaml` at the root of your git repository.
You can use the environment variable `AUTOSCOPE_CONFIG_FILE` to overwrite the location.

The file declares a list of `modules` for instance:

```yaml
modules:
    - name: api
      files: src/api/**/*
    - name: db
      files: src/db/**/*
    - filesRe: modules/(?P<scope>\w+)/.*
```

We support two kinds:

1. Named Glob: all files matching a glob belong to the module:

    ```yaml
    - name: docs
      files: docs/**/*
    ```

2. Captured with a Regexp: a regexp with a named capture group named `scope` is matched against file paths. The scope is extracted.

    ```yaml
    - filesRe: modules/(?P<scope>\w+)/.* # You can only do one match here
    ```

## Usage as Git Hook

The indented usage is as a git `commit-msg` hook.

Copy the script at [hooks/commit-msg](./hooks/commit-msg) to `.git/hooks/commit-msg`.
Make it executable:

```bash
chmod +x .git/hooks/commit-msg
```

Enjoy!
