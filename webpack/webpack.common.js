const path = require('path');
const webpack = require('webpack');
const { merge } = require('webpack-merge');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const utils = require('./utils.js');
const environment = require('./environment');

const getTsLoaderRule = env => {
  const rules = [
    {
      loader: 'thread-loader',
      options: {
        // There should be 1 cpu for the fork-ts-checker-webpack-plugin.
        // The value may need to be adjusted (e.g. to 1) in some CI environments,
        // as cpus() may report more cores than what are available to the build.
        workers: require('os').cpus().length - 1,
      },
    },
    {
      loader: 'ts-loader',
      options: {
        transpileOnly: true,
        happyPackMode: true,
      },
    },
  ];
  return rules;
};

module.exports = async options => {
  const development = options.env === 'development';
  return merge(
    {
      cache: {
        // 1. Set cache type to filesystem
        type: 'filesystem',
        cacheDirectory: path.resolve(__dirname, '../target/webpack'),
        buildDependencies: {
          // 2. Add your config as buildDependency to get cache invalidation on config change
          config: [
            __filename,
            path.resolve(__dirname, `webpack.${development ? 'dev' : 'prod'}.js`),
            path.resolve(__dirname, 'environment.js'),
            path.resolve(__dirname, 'utils.js'),
            path.resolve(__dirname, '../postcss.config.js'),
            path.resolve(__dirname, '../tsconfig.json'),
          ],
        },
      },
      resolve: {
        extensions: ['.js', '.jsx', '.ts', '.tsx', '.json'],
        modules: ['node_modules'],
        alias: utils.mapTypescriptAliasToWebpackAlias(),
        fallback: {
          path: require.resolve('path-browserify'),
        },
      },
      module: {
        rules: [
          {
            test: /\.tsx?$/,
            use: getTsLoaderRule(options.env),
            include: [utils.root('./webapp/app')],
            exclude: [utils.root('node_modules')],
          },
          /*
       ,
       Disabled due to https://github.com/jhipster/generator-jhipster/issues/16116
       Can be enabled with @reduxjs/toolkit@>1.6.1 
      {
        enforce: 'pre',
        test: /\.jsx?$/,
        loader: 'source-map-loader'
      }
      */
        ],
      },
      stats: {
        children: false,
      },
      optimization: {
        splitChunks: {
          cacheGroups: {
            commons: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              chunks: 'all',
            },
          },
        },
      },
      plugins: [
        new webpack.EnvironmentPlugin({
          // react-jhipster requires LOG_LEVEL config.
          LOG_LEVEL: development ? 'info' : 'error',
        }),
        new webpack.DefinePlugin({
          DEVELOPMENT: JSON.stringify(development),
          VERSION: JSON.stringify(environment.VERSION),
          SERVER_API_URL: JSON.stringify(environment.SERVER_API_URL),
        }),
        new ForkTsCheckerWebpackPlugin(),
      ],
    }
    // jhipster-needle-add-webpack-config - JHipster will add custom config
  );
};
