const fs = require('fs');
const path = require('path');

const docsPath = path.join(process.cwd(), 'docs', '需求.md');

if (!fs.existsSync(docsPath)) {
  console.log('NOT_EXIST');
} else {
  const stats = fs.statSync(docsPath);
  if (stats.size === 0) {
    console.log('EMPTY');
  } else {
    const content = fs.readFileSync(docsPath, 'utf-8');
    const lines = content.split('\n').slice(0, 5).join('\n');
    console.log('EXISTS');
    console.log(lines);
  }
}
