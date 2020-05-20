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

### Bug Report

A bug is a demonstrable problem which produce incorrect result or to behave in unintended ways.
Bug reports are helpful for developers who maintain the Vald project.

A good bug report should not leave others needing to ask you for more information constantly.
Please try to write as detail as possible in your bug report.

Things to check before reporting a bug:
- Check your environment

Bug may occur depending on different version of Vald, Kubernetes and etc.

- Identify the bug

Developers need to reproduce the bug in the developers environment.
So please identify the bug and clarify how it occurs.
Please also clarify the expected behavior too.

- Use the Github issue search

The bug may be reported before.
So please check the Github issue before reporting.

Please submit [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=)

### Proposal

The Vald is being developed based on the design-driven process.
The significant change to the library or the architecture should be discussed first.

We may need detailed documentation before your proposal is implemented.
Your proposal will be reviewed and discussed and decide whether it is accepted or declined.

Things to check before submitting your proposal:



Please submit [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=)

### Feature Request

Feature request is welcome.

Before opening an issue, please make sure your idea fits the project.
You can strongly request the feature and convince the project maintainers to accept your feature request.

Please provide the problem and solution associated with the feature request as detail as possible.

Things to check before submitting a new feature:

Please submit [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=))

### Security Issue Report

The Vald team and community take serious concern about security issues.

We appreciate your efforts to disclose your findings.
If the security issue is caused by third-part module, we will notify the person or team maintaining the module.

Please submit [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fsecurity%2C+priority%2Fmedium%2C+team%2Fcore%2C+team%2Fsre&template=security_issue_report.md&title=)

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
