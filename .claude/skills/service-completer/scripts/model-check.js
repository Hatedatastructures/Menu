import fs from 'fs'
import path from 'path'

const MODEL_DIR = path.join(process.cwd(), 'server', 'internal', 'model')

function getModelFiles() {
  if (!fs.existsSync(MODEL_DIR)) {
    console.error(`Model directory not found: ${MODEL_DIR}`)
    return []
  }

  const entries = fs.readdirSync(MODEL_DIR, { withFileTypes: true })
  const models = []

  for (const entry of entries) {
    if (entry.isFile() && entry.name.endsWith('.go') && !entry.name.startsWith('_')) {
      models.push(entry.name.replace('.go', ''))
    }
  }

  return models
}

const models = getModelFiles()
console.log(JSON.stringify(models, null, 2))
