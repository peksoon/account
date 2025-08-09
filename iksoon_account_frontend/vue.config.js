const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: [
    '@toast-ui/calendar'
  ],
  devServer: {
    port: 3000
  },
  configureWebpack: {
    plugins: [
      new (require('webpack')).DefinePlugin({
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false
      })
    ]
  }
});
