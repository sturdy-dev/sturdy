module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
    jest: true,
  },

  extends: [
    'plugin:@typescript-eslint/recommended',
    'eslint:recommended',
    'plugin:vue/vue3-recommended',
    'prettier',
  ],

  plugins: ['prettier', '@typescript-eslint', 'file-progress'],

  parser: 'vue-eslint-parser',

  parserOptions: {
    ecmaVersion: 2021,
    parser: '@typescript-eslint/parser',
  },

  ignorePatterns: ['dist', 'node_modules'],

  rules: {
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    // todo: research if this rule is important and maybe enable it
    '@typescript-eslint/explicit-module-boundary-types': 'off',
    'file-progress/activate': 1,

    'vue/multi-word-component-names': 'off',
  },
}
