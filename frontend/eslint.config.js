import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'

export default [
  js.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  {
    rules: {
      // Vue 3 相关
      'vue/multi-word-component-names': 'off',
      'vue/no-unused-vars': 'warn',
      'vue/require-default-prop': 'off',
      // JS 通用
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
      'prefer-const': 'warn',
    },
  },
  {
    ignores: ['dist/**', 'node_modules/**', '*.config.js'],
  },
]
