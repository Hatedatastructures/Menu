const fs = require('fs');
const path = require('path');

const content = process.argv[2] || '';
const docsDir = path.join(process.cwd(), 'docs');

if (!fs.existsSync(docsDir)) {
  fs.mkdirSync(docsDir, { recursive: true });
}

const docsPath = path.join(docsDir, '需求.md');
fs.writeFileSync(docsPath, content, 'utf-8');
console.log(docsPath);
