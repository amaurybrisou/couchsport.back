'use strict'

const webpack = require('webpack')
const merge = require('webpack-merge')
const baseConfig = require('./webpack.config.base')

const HOST = 'localhost'
const PORT = 8081

module.exports = merge(baseConfig, {
  mode: 'development',
  resolve: {
    alias: {
      vue: 'vue/dist/vue.js'
    }
  },
  output: {
    publicPath: '/'
  },
  devServer: {
    clientLogLevel: 'warning',
    hot: true,
    contentBase: 'static',
    compress: true,
    host: HOST,
    port: PORT,
    open: true,
    overlay: { warnings: false, errors: true },
    historyApiFallback: true,
    quiet: true,
    watchOptions: {
      poll: true
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        secure: false
      }
    }
  },

  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ]
      }, {
        test: /\.styl(us)?$/,
        use: [
          'vue-style-loader',
          'css-loader',
          'stylus-loader'
        ]
      }
    ]
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin()
  ]
})
