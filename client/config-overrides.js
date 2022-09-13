const {
  addLessLoader,
  adjustStyleLoaders,
  fixBabelImports,
  override,
  overrideDevServer,
} = require('customize-cra');
const { getThemeVariables } = require('antd/dist/theme');

module.exports = {
  webpack: override(
    fixBabelImports('import', {
      libraryName: 'antd',
      libraryDirectory: 'es',
      style: true,
    }),

    addLessLoader({
      lessOptions: {
        javascriptEnabled: true,
        modifyVars: getThemeVariables({
          '@link-color': '#a7864f',
          '@background-color': '#f0f0f0',
        }),
        cssModules: {
          localIdentName:
            process.env.NODE_ENV === 'development'
              ? '[local]--[hash:base64:5]'
              : '[hash:base64:5]',
        },
      },
    }),
    adjustStyleLoaders(({ use: [, , postcss] }) => {
      const postcssOptions = postcss.options;
      postcss.options = { postcssOptions };
    }),
  ),
};
