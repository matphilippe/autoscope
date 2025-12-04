# SVScope

SVScope is a QOL tool aimed at facilitating the use of conventional commits and tooling around the practice.

The TL;DR:

1. You get to define a tree/forest of modules in your repo: All files under directory `X/` are in module `my-module`.
2. If a commit touches a file in `X`, then the scope `my-module` will be appended to your commit scope, if it's not there already.

Why this matters?

1. Someone decided your module should be named: `terraform-provider-jebediah-db`, turning `fix: fixed the bug` into `fix(terraform-provider-jebediah-db): fixed the bug`. You are going to make a typo there occasionally.
2. Scope correctness actually matters if you would like to understand version bumps after merges. In practice, people are using PRs with some merge strategy relevant to their usecases.
   Squash commits are popular, with a commit message listing the previous commits, but you will not know which individual commits changed which files. Scopes are there for you, by putting the relevant info in the header.

We recommend using this as a `commit-msg` git hook :)

## Installation

```bash
go install github.com/mphilippe/svscope@latest
```

Or build from source:

```bash
git clone https://github.com/mphilippe/svscope
cd svscope
go build
```

## Usage as Git Hook

Create `.git/hooks/commit-msg`:

```bash
#!/bin/sh
svscope "$1"
```

Make it executable:

```bash
chmod +x .git/hooks/commit-msg
```

Now when you commit, scopes will be automatically added based on the files you changed!

## How It Works

The tool has minimal git interaction - it only runs `git diff --cached --name-only` to get the list of staged files.

The real work happens in:
1. **String transformation**: Parsing and reconstructing conventional commit messages with regex
2. **Pattern matching**: Using glob patterns (via doublestar) and regex to match files to scopes
3. **Scope assignment**: Deduplicating and appending scopes to the commit message

## Configuration

Obviously a yaml file: `.svscope.yaml` or `.svscope.yml`

```yaml
modules:
  - <Module | RegexpModule > # Either regexps with a capture group, or a module def.
```

A module def is simply:

```
name: my-scope
files: modules/my-scope/**/*
```

This will match a change to `modules/my-scope/folder/file.md` to the scope `my-scope`.

A RegexpModule is slightly more clever:

```
filesRe: modules/(?<scope>\w+)/.*
```

This will match a change to e.g. `modules/heey/main.tf` to the scope `heey`

## Limitations

You are in charge of a sensible structure for the module. The tool will apply all scopes it sees.

If two modules have a non-empty intersection, one should be including the entirety of the other.
