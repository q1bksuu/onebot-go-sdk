const fs = require('fs');

module.exports = async ({ github, context }) => {
  try {
    const coverage = fs.readFileSync('./coverage.out', 'utf8');
    const lines = coverage.split('\n');
    let totalStmt = 0;
    for (const line of lines) {
      if (line.includes('total:')) {
        totalStmt = line.split('\t')[1];
        break;
      }
    }
    await github.rest.issues.createComment({
      issue_number: context.issue.number,
      owner: context.repo.owner,
      repo: context.repo.repo,
      body: `📊 **Code Coverage Report**

Coverage: ${totalStmt || 'N/A'}

Full report: [View Coverage](${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${process.env.GITHUB_RUN_ID})`
    });
  } catch (error) {
    console.log('Could not generate coverage comment:', error);
  }
};
