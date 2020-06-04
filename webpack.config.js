const webpack = require('webpack');
const fs = require('fs');
const path = require('path');
const config = require('flarum-webpack-config')();

module.exports = [
  // flarum.core配置
  {
    entry: function () {
      const entries = {};
      for (const app of ['forum', 'admin']) {
        const file = path.resolve(process.cwd(), "view", "flarum", "js", app + '.js');
        if (fs.existsSync(file)) {
          entries[app] = file;
        }
      }
      return entries;
    }(),

    output: {
      path: path.resolve(process.cwd(), 'dist'),
      // library: 'module.exports',
      library: 'flarum.core',
      libraryTarget: 'assign',
      devtoolNamespace: require(path.resolve(process.cwd(), 'package.json')).name
    },
    module: config.module,
    devtool: config.devtool,
    devServer: {
      contentBase: path.join(__dirname, 'dist'),
      compress: true,
      host: '0.0.0.0',
      port: 9000
    },
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
  },
];
