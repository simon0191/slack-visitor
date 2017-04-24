const UglifyJsPlugin = require('webpack/lib/optimize/UglifyJsPlugin'),
  CompressionPlugin = require('compression-webpack-plugin'),
  DefinePlugin = require('webpack/lib/DefinePlugin'),
  webpackConfig = require("./webpack.config.base"),
  helpers = require("./helpers"),
  env = require('../environment/prod.env.js');

webpackConfig.entry["main.min"] = helpers.root("/src/main.ts");

webpackConfig.plugins = [...webpackConfig.plugins,
  new UglifyJsPlugin({
    include: /\.min\.js$/,
    minimize: true
  }),
  new CompressionPlugin({
    asset: "[path].gz[query]",
    test: /\.min\.js$/
  }),
  new DefinePlugin({
    'process.env': env
  })
];

module.exports = webpackConfig;
