const base = require('./base.config.js')
const merge = require('webpack-merge')
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = merge(base, {
    mode: 'development',
    devtool: "source-map",
    plugins: [
        new BundleAnalyzerPlugin({
            analyzerMode: 'static'
        }),
        new MiniCssExtractPlugin({
            filename: 'gen/[name].[contenthash].css',
            chunkFilename: 'gen/chunk-[name].[chunkhash].css',
        }),
    ],
    output: {
        filename: 'gen/[name].[contenthash].js',
        chunkFilename: 'gen/chunk-[name].[chunkhash].js',
    }
})
