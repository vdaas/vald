# Contributing to Vald

Vald is an open source project.

We appreciate your help!

## Table of Contents

- [Contributing Issue](#Contributing-Issue)
  - [Bug Report](#Bug-Report)
  - [Proposal](#Proposal)
  - [Feature Request](#Feature-Request)
  - [Security Issue Report](#Security-Issue-Report)
- [Contributing Source Code](#Contributing-Source-Code)
  - [Before contributing source code](#Before-contributing-source-code)
  - [How to contribute source code](#How-to-contribute-source-code)
  - [Branch naming convention](#Branch-naming-convention)

## Contributing Issue

We use [Github Issues](https://github.com/vdaas/vald/issues) to track issues within this repository.
If you can determine the problem you are facing is a bug or issue, you can easily submit the issues.

If you are unsure whether you have found a bug or security-related issues, please consider asking in the [Vald Slack](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) first.
If the behavior you are seeing is confirmed as a bug or issue, it can easily be re-raised in the issue tracker.

### Bug Report

A bug is a demonstrable problem that produces an incorrect result or behaves in unintended ways.<br>

Bug reports are helpful for developers who maintain the Vald project.<br>
A good bug report should not leave others needing to constantly ask you for more information.<br>
Please try to write as detailed as possible in your bug report.

When filing an issue, make sure to answer these five questions:

1. What version of Vald are you using (vald docker image version)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?
5. What did you see instead?

Please submit the bug report [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=)

### Proposal

The Vald is being developed based on the design-driven process.<br>
The significant change to the library or the architecture should be discussed first.

We may ask for detailed documentation before your proposal is accepted.<br>
Your proposal will be reviewed, discussed, and decided whether it is accepted or declined.

Please submit the proposal [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=)

### Feature Request

Feature request is welcome.

Before opening an issue, please make sure your idea fits the project.<br>
You can request the feature and convince the project maintainers to accept your feature request.

Please provide the problem and solution associated with the feature request in as much detail as possible.

NOTE: If you’d like to contribute to the new feature which may affect the current Vald architecture or design, you should discuss it with the Vald team first.

Please submit the feature request [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=)

### Security Issue Report

The Vald team and community have a serious concern about security issues.

Sensitive security-related issues should be reported to the security channel in [Vald Slack](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA).

We appreciate your efforts to disclose your findings.<br>
If the security issue is caused by a third-party module, we will contact the module owner and ask for a fix.

We will consider using another third-party module if the vulnerable module is not actively maintained anymore.

Please submit the security issue report [here](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fsecurity%2C+priority%2Fmedium%2C+team%2Fcore%2C+team%2Fsre&template=security_issue_report.md&title=)

## Contributing Source Code

### Before contributing source code

Follow these steps to make a contribution to any of our open source repositories:

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald)
2. Set your name and email (these should match the information on your submitted CLA)

   ```bash
   git config --global user.name "Firstname Lastname"
   git config --global user.email "your_email@example.com"
   ```

### How to contribute source code

1. Fork the repository ( https://github.com/vdaas/vald/fork )
2. Create your feature branch (git checkout -b [`[type]/[area]/[description]`](#Branch-naming-convention))
3. Commit your changes on your branch (git commit -am 'Add some feature')
4. Run tests (make test)
5. Push to the forked branch (git push origin my-new-feature)
6. Create new Pull Request

Each pull request and commit should be small enough to contain only one purpose.

Your pull request is much more likely to be accepted if:

- Your pull request includes tests

- Your pull request includes benchmark results

- Your pull request is small and focused with a clear message that conveys the intent of your change

### Branch naming convention

Name your branches with prefixes: `[type]/[area]/[description]`

| Field       | Explanation                           | Naming Rule                                                                                                               |
| :---------- | :------------------------------------ | :------------------------------------------------------------------------------------------------------------------------ |
| type        | The PR type                           | The type of PR can be feature, bug, refactoring, benchmark, security, documentation, dependencies, ci, test, or etc...    |
| area        | Area of context                       | The area of PR can be gateway, agent, agent-sidecar, lb-gateway, or etc...                                                |
| description | Summarized description of your branch | The description must be hyphenated. Please use [a-zA-Z0-9] and hyphen as characters, and do not use any other characters. |

(\*) If you changed multiple areas, please list each area with "-".

For example, when you add a new feature for internal/servers, the name of the branch will be `feature/internal/add-newfeature-for-servers`.
