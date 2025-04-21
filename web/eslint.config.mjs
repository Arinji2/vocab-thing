import { tanstackConfig } from '@tanstack/eslint-config'

export default [
  ...tanstackConfig,

  {
    rules: {
      'import/order': 'off',
      'import/consistent-type-specifier-style': 'off',
      'sort-imports': 'off',
    },
  },
]
