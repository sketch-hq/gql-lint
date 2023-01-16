Hi there! We’re thrilled that you’d like to contribute to this project.

Contributions to this project are released to the public under the [MIT license](LICENSE).

Please note that this project comes with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

### Did you find a bug?

If the bug is a security vulnerability that affects Sketch **do not** open up a GitHub issue. If the security vulnerability affects Sketch applications, please contact Sketch following [Sketch’s Responsible Disclosure Policy](https://www.sketch.com/security/disclosure/). If the security vulnerability affects other applications, please report it through their appropriate secure channels.

For any other bug, first, make sure the bug isn’t already reported by searching for it on GitHub under [Issues](https://github.com/sketch-hq/gql-lint/issues).

If you’re unable to find an open issue addressing the problem, open a new one. Be sure to include a title and clear description, as much relevant information as possible, and an error trace or something similar. If a written description is not enough to describe the problem, feel free to attach screenshots or videos explaining the problem.

## Do you have doubts?

If you have doubts or there's anything you want to talk about that does not make sense to add it as a new issue, you can do it writing in the [discussions area](https://github.com/sketch-hq/gql-lint/discussions).

## Missing a feature?

You can *request* a new feature by submitting an issue to our GitHub Repository. If you would like to *implement* a new feature, please consider the size of the change in order to determine the right steps to proceed:

- For a **Major Feature**, the first step is to open an issue and outline your proposal so that it can be discussed. This process allows us to better coordinate our efforts, avoid duplicated work, and help you to craft the change so that it’s successfully accepted as part of the project.
    
    **Note**: Adding a new topic to the documentation, or significantly re-writing a topic, counts as a major feature.
    
- **Small Features** can be crafted and directly [submitted as a Pull Request](about:blank#submitting-a-pull-request).

## The Developer Certificate of Origin (DCO)

The Developer Certificate of Origin (DCO) is a lightweight way for contributors to certify that they wrote or otherwise have the right to submit the code they are contributing to the project. Here is the full [text of the DCO](https://developercertificate.org/), reformatted for readability:

> By making a contribution to this project, I certify that:
> 
> 
> (a) The contribution was created in whole or in part by me and I have the right to submit it under the open source license indicated in the file; or
> 
> (b) The contribution is based upon previous work that, to the best of my knowledge, is covered under an appropriate open source license and I have the right under that license to submit that work with modifications, whether created in whole or in part by me, under the same open source license (unless I am permitted to submit under a different license), as indicated in the file; or
> 
> (c) The contribution was provided directly to me by some other person who certified (a), (b) or (c) and I have not modified it.
> 
> (d) I understand and agree that this project and the contribution are public and that a record of the contribution (including all personal information I submit with it, including my sign-off) is maintained indefinitely and may be redistributed consistent with this project or the open source license(s) involved.
> 

Contributors *sign-off* that they adhere to these requirements by adding a `Signed-off-by` line to commit messages.

```
This is my commit message

Signed-off-by: Random J Developer <random@developer.example.org>
```

Git even has a `-s` command line option to append this automatically to your commit message:

```
$ git commit -s -m 'This is my commit message'
```

## Submitting a pull request

1. Fork and clone the repository
2. Configure and install the dependencies
3. Make sure the tests pass on your machine
4. Create a new branch: `git checkout -b my-branch-name`
5. Make your change, add tests, and make sure the tests still pass.
6. Push your fork and submit a pull request
7. Pat yourself on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of getting your pull request approved:

- Follow the same style guide you observe in the current source code
- Write tests
- Write comments in the code, whenever you are doing something that is not straightforward or the intention may not seem obvious
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests
- Write a good commit message that explains briefly the purpose of your change. Also, take the chance to properly explain your intention in the PR description
- Keep an eye on the GH PR, and make sure all webhooks are passing as expected

## Code guidelines

It’s important to keep quality and style consistency across the application. We have different rules and use different tools to guarantee this consistency

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
