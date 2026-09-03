import fs from 'fs'
import path from 'path'

const ERROR_CODE_FILE = path.join(process.cwd(), 'server', 'internal', 'error', 'error_code.go')

function getErrorCodes() {
  if (!fs.existsSync(ERROR_CODE_FILE)) {
    console.error(`Error code file not found: ${ERROR_CODE_FILE}`)
    return {}
  }

  const content = fs.readFileSync(ERROR_CODE_FILE, 'utf-8')
  const errorCodes = {}

  const constRegex = /(\w+)\s*=\s*0x[0-9A-F]+\s*\/\/\s*(.+)/g
  let match
  while ((match = constRegex.exec(content)) !== null) {
    errorCodes[match[1]] = match[2]
  }

  return errorCodes
}

const codes = getErrorCodes()
console.log(JSON.stringify(codes, null, 2))
