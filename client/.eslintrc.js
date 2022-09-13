const config = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: './tsconfig.json',
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/eslint-recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:import/typescript',
    'prettier',
    'react-app',
  ],

  rules: {
    semi: ['error', 'always'],
    curly: ['error'],
    quotes: ['error', 'single', { avoidEscape: true }],
    'no-duplicate-imports': ['error'],
    'dot-notation': 'error',
    'no-return-await': 'error',
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/no-unnecessary-condition': 'error',
    '@typescript-eslint/no-unused-vars': [
      'error',
      { ignoreRestSiblings: true },
    ],
    '@typescript-eslint/no-inferrable-types': [
      'warn',
      { ignoreParameters: true },
    ],
    'sort-imports': ['warn', { ignoreDeclarationSort: true }],
    'import/order': [
      'warn',
      {
        alphabetize: {
          order: 'asc',
          caseInsensitive: true,
        },
      },
    ],
    'tsdoc/syntax': 'warn',
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',
    'react/jsx-curly-brace-presence': [
      'error',
      { props: 'never', children: 'never' },
    ],
    'react/jsx-sort-props': [
      'warn',
      { shorthandFirst: true, multiline: 'last', noSortAlphabetically: true },
    ],
    'react/jsx-no-useless-fragment': ['error', { allowExpressions: true }],
    'react/jsx-key': 'error',
  },

  plugins: [
    '@typescript-eslint/eslint-plugin',
    'import',
    'eslint-plugin-tsdoc',
    'react-hooks',
  ],

  overrides: [
    {
      files: ['*.js'],
      rules: {
        '@typescript-eslint/no-var-requires': 'off',
      },
    },
  ],
};

module.exports = config;
