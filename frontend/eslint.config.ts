import { globalIgnores } from 'eslint/config'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import pluginVue from 'eslint-plugin-vue'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'

// Flat config (ESLint 9+). Order matters — later entries override earlier ones:
//   1. which files to lint
//   2. paths to ignore entirely
//   3. Vue's essential rules (template/SFC correctness)
//   4. the TS-aware recommended ruleset (no-explicit-any, no-unused, etc.)
//   5. skipFormatting LAST — turns off every rule that overlaps Prettier, so
//      ESLint never fights the formatter.
export default defineConfigWithVueTs(
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  globalIgnores(['**/dist/**', '**/node_modules/**']),
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  skipFormatting,
)
