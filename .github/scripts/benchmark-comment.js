const fs = require('fs');

module.exports = async ({ github, context }) => {
  const result = fs.readFileSync('benchstat.txt', 'utf8');

  const body = `### 📊 Go Benchmark Result

${result}
`;

  await github.rest.issues.createComment({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
    body
  });
};
