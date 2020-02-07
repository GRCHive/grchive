const base = require('./base.config.js')
const merge = require('webpack-merge')
const TerserPlugin = require('terser-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const webpack = require('webpack')

module.exports = merge(base, {
    mode: 'production',
    optimization: {
        minimize: true,
        minimizer: [new TerserPlugin({
            parallel: true,
        })],
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'gen/[name].[contenthash].css',
            chunkFilename: 'gen/[id].[chunkhash].css',
        }),
        new webpack.DefinePlugin({
            __WEBSOCKET_PROTOCOL: JSON.stringify("wss://")
        }),
    ],
    output: {
        filename: 'gen/[name].[contenthash].js',
        chunkFilename: 'gen/[id].[chunkhash].js',
    }
})
