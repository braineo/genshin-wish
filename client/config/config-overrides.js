const { override, fixBabelImports, addLessLoader } = require("customize-cra");
const { getThemeVariables } = require("antd/dist/theme");

module.exports = {
  webpack: override(
    fixBabelImports("import", {
      libraryName: "antd",
      libraryDirectory: "es",
      style: true,
    }),

    addLessLoader({
      lessOptions: {
        javascriptEnabled: true,
        modifyVars: getThemeVariables({
        }),
        cssModules: {
          localIdentName:
            process.env.NODE_ENV === "development"
              ? "[local]--[hash:base64:5]"
              : "[hash:base64:5]",
        },
      },
    })
  ),
};
