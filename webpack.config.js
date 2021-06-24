const webpack = require('webpack');
const fs = require('fs');
const path = require('path');
const config = require('flarum-webpack-config')();
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

let devServer = {
  publicPath: '/webpack/static/flarum/',
  // contentBase: path.join(__dirname, 'dist'),
  // compress: true,
  host: '0.0.0.0',
  port: 9000,
  headers: {
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
    "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
  }
};
let OUTPUT_PATH = path.resolve(process.cwd(), 'static', 'flarum');
let FLARUM_DIR = path.resolve(process.cwd(), 'view', 'flarum');
let STATIC_DIR = path.resolve(process.cwd(), 'view');
let EXT_DIR = path.resolve(process.cwd(), "view", "extensions");

module.exports = [
  // flarum.core配置
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const file = path.resolve(STATIC_DIR, app + '.js');
        if (fs.existsSync(file)) {
          entries[app] = file;
        }
      }
      return entries;
    }(),
    plugins: config.plugins,

    output: {
      path: OUTPUT_PATH,
      library: 'flarum.core',
      libraryTarget: 'assign',
      devtoolNamespace: require(path.resolve(process.cwd(), 'package.json')).name
    },
    module: config.module,
    devtool: config.devtool,
    devServer: devServer,
    resolve: config.resolve
  },
  // flarum的一些扩展功能
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const file = path.resolve(EXT_DIR, app + ".js");
        entries[`${app}_ext`] = file;
      }
      return entries;
    }(),

    module: config.module,
    externals: config.externals,
    devtool: config.devtool,

    output: {
      path: OUTPUT_PATH,
    }
  },
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const file = path.resolve(EXT_DIR, app + '.less');
        if (fs.existsSync(file)) {
          entries[app] = file;
        }
      }
      return entries;
    }(),

    module: {
      rules: [{
        test: /\.less$/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
            options: {
              publicPath: ''
            }
          },
          'css-loader',
          {
            loader: 'less-loader',
            options: {
              sourceMap: true,
              lessOptions: {
                paths: [
                  path.resolve(__dirname, 'node_modules/components-font-awesome/less/'),
                  path.resolve(__dirname, 'node_modules/bootstrap/less/'),
                  path.resolve(FLARUM_DIR, 'less'),
                  path.resolve(FLARUM_DIR, 'less', 'common'),
                ],
              },
            },
          },
        ],
      },
      {
        test: /\.(eot|svg|ttf|woff|woff2|png)\w*/,
        loader: 'file-loader',
        options: {
          name: './fonts/[name].[ext]',
        },
      },
      ]
    },
    plugins: [new MiniCssExtractPlugin()],
    output: {
      path: OUTPUT_PATH,
    },
  },
];