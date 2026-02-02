module.exports = async ({ github, context }) => {
  await github.rest.issues.createComment({
    issue_number: context.issue.number,
    owner: context.repo.owner,
    repo: context.repo.repo,
    body: '❌ The GolangCI-Lint check failed. Please check the logs above for details.'
  });
};
