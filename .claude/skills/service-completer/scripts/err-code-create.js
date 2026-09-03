import fs from 'fs'
import path from 'path'
import readline from 'readline'

const ERROR_CODE_FILE = path.join(process.cwd(), 'server', 'internal', 'error', 'error_code.go')

const args = process.argv.slice(2)
if (args.length < 3) {
  console.error('Usage: node err-code-create.js <CODE_NAME> <HEX_VALUE> <MESSAGE>')
  process.exit(1)
}

const [codeName, hexValue, message] = args

async function insertErrorCode() {
  if (!fs.existsSync(ERROR_CODE_FILE)) {
    console.error(`Error code file not found: ${ERROR_CODE_FILE}`)
    return
  }

  const content = fs.readFileSync(ERROR_CODE_FILE, 'utf-8')
  const lines = content.split('\n')

  let insertIndex = -1
  let inModuleSection = false
  let lastConstIndex = -1

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]

    if (line.includes('// 模块业务错误码')) {
      inModuleSection = true
    }

    if (inModuleSection && line.includes('const (')) {
      lastConstIndex = i
    }

    if (inModuleSection && /^\s*\)/.test(line) && lastConstIndex >= 0) {
      insertIndex = i
      break
    }
  }

  if (insertIndex === -1) {
    console.error('Could not find insertion point in error code file')
    return
  }

  const newLine = `\t${codeName} = ${hexValue} // ${message}`

  lines.splice(insertIndex, 0, newLine)

  fs.writeFileSync(ERROR_CODE_FILE, lines.join('\n'))
  console.log(`Inserted error code: ${codeName}`)
}

insertErrorCode()
