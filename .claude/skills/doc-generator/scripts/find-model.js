/**
 * 自动寻找 model 目录下第一层子目录中的 .go 文件（排除 request/response 子目录）
 * 路径结构：server/internal/model/xxx/*.go（但排除 xxx/request/*.go 和 xxx/response/*.go）
 */

const fs = require('fs');
const path = require('path');

/**
 * 向上查找包含 server 目录的路径
 */
function findServerDirectory(startPath) {
  let currentPath = startPath;
  const maxDepth = 5;

  for (let i = 0; i < maxDepth; i++) {
    const serverPath = path.join(currentPath, 'server');
    if (fs.existsSync(serverPath) && fs.statSync(serverPath).isDirectory()) {
      return serverPath;
    }
    const parentPath = path.dirname(currentPath);
    if (parentPath === currentPath) break;
    currentPath = parentPath;
  }

  return null;
}

/**
 * 搜索 model 目录下第一层子目录中的 .go 文件
 * 排除 request, response, common, enum, system 目录
 * @param {string} serverPath - server 目录路径
 * @returns {string[]} - .go 文件的绝对路径数组
 */
function findModelFiles(serverPath) {
  const modelPath = path.join(serverPath, 'internal', 'model');

  if (!fs.existsSync(modelPath)) {
    console.error('Model directory not found:', modelPath);
    return [];
  }

  console.log('Found model directory:', modelPath);

  const modelFiles = [];
  const skipDirs = ['common', 'enum', 'system', 'request', 'response'];

  // 读取 model 目录下的第一层子目录
  const entries = fs.readdirSync(modelPath, { withFileTypes: true });

  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    if (skipDirs.includes(entry.name)) continue;

    const subDir = path.join(modelPath, entry.name);

    // 只读取子目录中的文件，不递归
    const subEntries = fs.readdirSync(subDir, { withFileTypes: true });

    for (const subEntry of subEntries) {
      if (subEntry.isFile() && subEntry.name.endsWith('.go') && !subEntry.name.endsWith('_test.go')) {
        modelFiles.push(path.join(subDir, subEntry.name).replace(/\\/g, '/'));
      }
    }
  }

  return [...new Set(modelFiles)].sort();
}

/**
 * 主函数
 */
function main() {
  const startPath = process.cwd();
  const serverPath = findServerDirectory(startPath);

  if (!serverPath) {
    console.error('Could not find server directory');
    return [];
  }

  const files = findModelFiles(serverPath);
  console.log(JSON.stringify(files, null, 2));
}

// 导出函数
module.exports = { findModelFiles, findServerDirectory };

if (require.main === module) {
  main();
}