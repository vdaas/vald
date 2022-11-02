# Contributing to Vald

Thank you for your interest in Vald, and thank you for investing your time in contributing to Vald!

Read our [Code Of Conduct](https://github.com/vdaas/vald/blob/main/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

In this guide, you will get an idea of how to contribute to Vald.

If you are not a developer, don't worry, some contributions don't require writing a single line of code.

## New contributor guide

Please read the [README](https://github.com/vdaas/vald/blob/main/README.md) to get an overview of Vald.

We welcome you to contribute to Vald to make Vald better. We accept the following types of contributions:

- Issue
  - Bug report
  - Feature request
  - Proposal
  - Security issue report
- Pull request

Please also feel free to ask anything on [Vald Slack channel](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA) :)

If you decided to contribute source code changes, you may need to install the following tools manually before making changes.

- [git](https://git-scm.com/)
- [go](https://go.dev/)

Also, you may need to run the following command under [Vald repository](https://github.com/vdaas/vald) to install the necessary packages.

```bash
make init
make tools/install
make ngt/install
```

## Issue

[Issues](https://github.com/vdaas/vald/issues) are used to track tasks that contributors can help with.

An issue can be a bug report, feature request, or security issue report.

We welcome you to give us ideas to improve Vald or report any issue that exists in Vald by creating an issue on the Vald repository.

Please find [here](../contributing/issue.md) for more details about the issue, and please also feel free to ask anything on [Vald Slack channel](https://join.slack.com/t/vald-community/shared_invite/zt-db2ky9o4-R_9p2sVp8xRwztVa8gfnPA).

### Propose a new issue

If you have found something that should be updated in Vald, search open issues to see if someone else has already talked about it.

If it is something new, please open a new issue using the following template:

- [Bug report](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fbug%2C+priority%2Fmedium%2C+team%2Fcore&template=bug_report.md&title=)
- [Feature request](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Ffeature%2C+priority%2Flow%2C+team%2Fcore&template=feature_request.md&title=)
- [Security issue report](https://github.com/vdaas/vald/issues/new?assignees=&labels=type%2Fsecurity%2C+priority%2Fmedium%2C+team%2Fcore%2C+team%2Fsre&template=security_issue_report.md&title=)

Please fill in the information written on the template to help us to understand what you think.

### Solve an issue

Please find our [opening issues](https://github.com/vdaas/vald/issues) and find one that interested you.

You can find the issue using the label `type/*`. For example, you can find all the bug issues using the label `type/bug`.

## Make changes

Before making changes on Vald, please follow these steps to contribute to any of our open-source repositories:

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald)
2. Set your name and email (these should match the information on your submitted CLA)

   ```bash
   git config --global user.name "Firstname Lastname"
   git config --global user.email "your_email@example.com"
   ```

### How to make changes

1. Comment on the issue and say you are working on it to avoid conflict with others also working on the same issue.
2. Fork the repository. ( https://github.com/vdaas/vald/fork )
3. Create your feature branch. ( git checkout -b [`[type]/[area]/[description]`](#Branch-naming-convention) )
4. Make code changes. Please follow the design or requirement discussed on the issue.
5. Commit your changes to your branch. ( git commit -am 'Add some feature' )
6. Push to the forked branch. ( git push origin my-new-feature )
7. Create a new pull request against the Vald repository. Please also mention the issue on the pull request.
8. Wait for the code review. Resolve any issue/questions raised by reviewers until it is merged.

### Pull request

After you have finished making code changes, you may create a pull request to ask for accepting the changes.

Each pull request and commit should be small enough to contain only one purpose, for easier review and tracking.
Please fill in the description on the pull request and write down the overview of the changes.

If you are solving an issue, please also link the pull request to the issue.

Vald team will review the pull request. We may ask for changes or questions we have before the pull request is merged.

After the pull request is merged, your changes will be applied to Vald, and the changes will be included in the next Vald release.

### Branch naming convention

Name your branches with prefixes: `[type]/[area]/[description]`

| Field       | Explanation                           | Naming Rule                                                                                                               |
| :---------- | :------------------------------------ | :------------------------------------------------------------------------------------------------------------------------ |
| type        | The PR type                           | The type of PR can be a feature, bug, refactoring, benchmark, security, documentation, dependencies, ci, test, etc...    |
| area        | Area of context                       | The area of PR can be gateway, agent, agent-sidecar, lb-gateway, etc...                                                |
| description | Summarized description of your branch | The description must be hyphenated. Please use [a-zA-Z0-9] and a hyphen as characters, and do not use any other characters. |

(\*) If you changed multiple areas, please list each area with "-".

For example, when you add a new feature for internal/servers, the name of the branch will be `feature/internal/add-newfeature-for-servers`.
