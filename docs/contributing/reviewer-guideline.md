# Guideline for the Pull Request reviewers

This document serves as a guideline for those reviewing Pull Requests.

This guideline is intended for internal developers. A separate guideline for external contributors is planned.

## Purpose

The purpose of this guideline is to clarify what reviewers should be aware of and to optimize communication, thereby facilitating efficient development of Vald.

## Preparation

This guideline assumes that the content of the Pull Request is appropriate.

Consider splitting the Pull Request if the changes are extensive, have multiple intentions, or if the tests are insufficient. This makes it easier to review.

Regarding the amount of change, opinions may vary, but in this project, we aim to keep the implementation of logic under 1000 lines, excluding automatically generated code such as Protocol Buffers.

## Review request

In this project, merging a Pull Request requires approval from two people. If you want a specific person to review, assign them explicitly; otherwise, assign vdaas/maintainer for a random assignment.

### Request notification

People may notice requests through different means, but in this project, we prioritize the following methods of communication:

1. Speak on Zoom if it is during business hours and the reviewer is available on Zoom; we often talk directly.
1. Send a mention on Slack If the reviewer is not on Zoom, we send a mention on Slack.
1. Email from GitHub An email is automatically sent when someone is assigned.

If the reviewer does not notice the request, we follow the same priority for re-notifying.

#### When Requested for Review

Unless there are high-priority issues, you are in a meeting, or deeply focused on development, start the review immediately.

## Approve condition

- The content of the Pull Request is understood.
- The content can be understood at a glance.
  - If it cannot be understood at a glance, request additional comments.
- There are no deficiencies in the test cases.
- The quality is suitable for release.
  - Compatibility
  - Security
  - Performance

Regarding quality, while there is always room for improvement, if the quality is sufficient for release, it is acceptable to approve. If there is potential for improvement, suggest it in the comments, and let the implementer decide whether to address it in the same Pull Request or a separate one.

### What Reviewers Should Do

- If the title or description is unclear, ask questions before reviewing the code.
- If any part of the code is not understood after a brief review, ask for clarification or additional comments instead of trying to understand it yourself.
- Verify if the test cases are sufficient.
- Approve immediately if there are no comments.
- If interrupted, respond even if the review is incomplete.
- Ignore other reviewers' comments.

### What Reviewers Should Not Do

- Try to decipher something that is not immediately understandable.
- Close the browser or tab or start something else if it is not understood.
- Attempt to verify the implementation in detail.

### Communication Methods

All communications do not need to be completed on GitHub. Using Slack or Zoom can make the review process more efficient.

The following criteria often determine the choice of communication tool:

- On GitHub Mainly used for specific questions or comments about the code. If the exchange is unsuitable for comments, use Slack or Zoom.
- On Slack, used for brief exchanges. If the discussion is likely to be lengthy or complex, switch to Zoom.
- On Zoom Used when more complex communication is necessary.

Regardless of where communication takes place, ensure that the results are documented on GitHub or in the code as a record of the communication.

## Effect measurement

Once a month, a report is created at [Vald Issues](https://github.com/vdaas/vald/issues).

This report includes the time taken from when a Pull Request is opened to when it is merged and a summary of this data. Regularly reflect on this information.
