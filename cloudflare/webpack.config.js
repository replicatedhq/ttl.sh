const path = require("path")
const webpack = require('webpack');

module.exports = {
  entry: {
    bundle: path.join(__dirname, "./src/index.ts"),
  },

  output: {
    filename: "bundle.js",
    path: path.join(__dirname, "dist"),
  },

  mode: process.env.NODE_ENV || "development",

  plugins: [
    // Ignore all locale files of moment.js
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
  ],

  watchOptions: {
    ignored: /node_modules|dist|\.js/g,
  },

  devtool: "cheap-module-source-map",

  resolve: {
    extensions: [".ts", ".tsx", ".js", ".json"],
    plugins: [],
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loader: "awesome-typescript-loader",
      },
    ],
  },
}
