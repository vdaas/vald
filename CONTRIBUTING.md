# Contributing to Vald

Vald is an open source project.

We appreciate your help!

## Before filing an issue

If you are unsure whether you have found a bug, please consider asking in the [vald slack](https://vald-community.slack.com/messages/CN1RDC8NB) first. If
the behavior you are seeing is confirmed as a bug or issue, it can easily be re-raised in the issue tracker.

## Filing issues

Sensitive security-related issues should be reported to [vald slack security](https://vald-community.slack.com/messages/CN1RDC8NB) .

Otherwise, when filing an issue, make sure to answer these five questions:

1. What version of Vald are you using (vald docker image version)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?
5. What did you see instead?

## Before contributing code

Follow these steps to make a contribution to any of our open source repositories:

1. Ensure that you have completed our [CLA Agreement](https://cla-assistant.io/vdaas/vald)
2. Set your name and email (these should match the information on your submitted CLA)

        git config --global user.name "Firstname Lastname"
        git config --global user.email "your_email@example.com"


## Contributing code

1. Fork the repository ( https://github.com/vdaas/vald/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes on your branch (git commit -am 'Add some feature')
4. Run tests (make test)
5. Push to the forked branch (git push origin my-new-feature)
6. Create new Pull Request


We favor pull requests with very small, single commits with a single purpose.

Your pull request is much more likely to be accepted if:

* Your pull request includes tests

* Your pull request includes benchmark results

* Your pull request is small and focused with a clear message that conveys the intent of your change
