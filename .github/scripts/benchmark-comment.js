const fs = require('fs');

module.exports = async ({ github, context }) => {
  const result = fs.readFileSync('benchstat.txt', 'utf8');
  const marker = '<!-- benchmark-comment -->';

  const owner = context.repo.owner;
  const repo = context.repo.repo;
  const issue_number = context.issue.number;

  const { data: comments } = await github.rest.issues.listComments({
    owner,
    repo,
    issue_number,
    per_page: 100
  });

  const existingComment = [...comments].reverse().find((comment) => {
    if (!comment || !comment.body) return false;
    return comment.body.includes(marker);
  });

  const body = `### 📊 Go Benchmark Result

${result}
${marker}
`;

  if (existingComment) {
    await github.rest.issues.updateComment({
      owner,
      repo,
      comment_id: existingComment.id,
      body
    });
    return;
  }

  await github.rest.issues.createComment({
    owner,
    repo,
    issue_number,
    body
  });
};
