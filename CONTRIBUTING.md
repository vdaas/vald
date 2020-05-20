# Contributing to Vald

Vald is an open source project.

We appreciate your help!

## Table of Contents

- [Contributing Issue](#Contributing-Issue)
  - [Before filing an issue](#Before-filing-an-issue)
  - [Filing issues](#Filing-issues)
  - [Bug Report](#Bug-Report)
  - [Proposal](#Proposal)
  - [Feature Request](#Feature-Request)
  - [Security Issue Report](#Security-Issue-Report)
- [Contributing Code](#Contributing-Code)
  - [Before contributing code](#Before-contributing-code)
  - [Contributing code](#Contributing-code)
  - [Branch naming convention](#Branch-naming-convention)

## Contributing Issue

### Before filing an issue

If you are unsure whether you have found a bug, please consider asking in the [Vald Slack](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) first. If the behavior you are seeing is confirmed as a bug or issue, it can easily be re-raised in the issue tracker.

### Filing issues

Sensitive security-related issues should be reported to the security channel in the [Vald Slack](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) .

Otherwise, when filing an issue, make sure to answer these five questions:

1. What version of Vald are you using (vald docker image version)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?
5. What did you see instead?

Note:
When you'd like to contribute with the new feature of some large update influencing the current Vald, we recommend getting the agreement through the authors' design review of your ISSUE.


## Contributing Code

### Before contributing code

Follow these steps to make a contribution to any of our open source repositories:

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald)
2. Set your name and email (these should match the information on your submitted CLA)

        git config --global user.name "Firstname Lastname"
        git config --global user.email "your_email@example.com"


### Contributing code

1. Fork the repository ( https://github.com/vdaas/vald/fork )
2. Create your feature branch (git checkout -b [`[type]/[area]/[description]`](#Branch-naming-convention))
3. Commit your changes on your branch (git commit -am 'Add some feature')
4. Run tests (make test)
5. Push to the forked branch (git push origin my-new-feature)
6. Create new Pull Request


We favor pull requests with very small, single commits with a single purpose.

Your pull request is much more likely to be accepted if:

* Your pull request includes tests

* Your pull request includes benchmark results

* Your pull request is small and focused with a clear message that conveys the intent of your change


### Branch naming convention

Name your branches with prefixes: `[type]/[area]/[description]`

* `type` = feature, bug, refactoring, benchmark, security, documentation, dependencies, ci, ...
* `area` (\*) = gateway, meta, manager-backup, manager-replication, ...
* `description` = branch description. description must be hyphenated.

(\*) If you changed multiple areas, please list up each area with "-".
