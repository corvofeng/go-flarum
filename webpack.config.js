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

module.exports = [
  // flarum.core配置
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const file = path.resolve(FLARUM_DIR, "js", app + '.js');
        if (fs.existsSync(file)) {
          entries[app] = file;
        }
      }
      return entries;
    }(),

    output: {
      path: OUTPUT_PATH,
      library: 'flarum.core',
      libraryTarget: 'assign',
      devtoolNamespace: require(path.resolve(process.cwd(), 'package.json')).name
    },
    module: config.module,
    devtool: config.devtool,
    devServer: devServer
  },
  // flarum的一些扩展功能
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const extDir = path.resolve(process.cwd(), "view", "extensions");
        const files = fs.readdirSync(extDir)
        files.forEach((f) => {
          entries[`${f}_${app}`] = path.resolve(extDir, f, "js", app + ".js");
        });
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
        const file = path.resolve(FLARUM_DIR, "less", app + '.less');
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
          MiniCssExtractPlugin.loader,
          'css-loader',
          {
            loader: 'less-loader',
            options: {
              prependData: (loaderContext) => {
                const { resourcePath, rootContext } = loaderContext;
                const relativePath = path.relative(rootContext, resourcePath);

                const variable_data = fs.readFileSync(path.resolve(FLARUM_DIR, 'less', 'common', 'variables.less'), 'utf-8')
                const mixin_data = fs.readFileSync(path.resolve(FLARUM_DIR, 'less', 'common', 'mixins.less'), 'utf-8')

                return variable_data + mixin_data;
              },
              appendData: (loaderContext) => {
                const webfont = path.resolve(FLARUM_DIR, 'node_modules/components-font-awesome/webfonts')
                // return `@fa-font-path: "${webfont}";`;
                return `@fa-font-path: "../webfonts";`;
              },

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
          // limit: 50000,
          // mimetype: 'application/font-woff',

          // // Output below the fonts directory
          name: './fonts/[name].[ext]',
          // Tweak publicPath to fix CSS lookups to take
          // the directory into account.
          // publicPath: 'fonts',
        },
      },
      ]
    },
    plugins: [new MiniCssExtractPlugin()],
    output: {
      path: OUTPUT_PATH,
    }
  },
];