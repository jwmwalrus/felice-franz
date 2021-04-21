const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');

const devmode = process.env.DEV_MODE === '1';

module.exports = {
    name: 'index',
    entry: {
        index: './js/index.js',
    },
    output: {
        path: path.join(__dirname, '..', 'public'),
        filename: 'js/[name].js',
        publicPath: path.join('..', 'public'),
    },
    target: 'es2021',
    cache: false,
    stats: 'minimal',
    parallelism: 2,
    devtool: devmode ? 'eval-source-map' : 'source-map',
    performance: {
        maxEntrypointSize: 524288,
        maxAssetSize: 524288,
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            // sourceMap: true,
                        },
                    },
                    {
                        loader: 'css-loader',
                        options: {
                            sourceMap: true,
                        },
                    },
                ],
            },
            {
                test: /\.(png|jpg|jpeg|gif|ico|svg)$/,
                use: {
                    loader: 'file-loader',
                    options: {
                        name: 'img/[name].[ext]',
                    },
                },
            },
            {
                test: /\.(woff|woff2|eot|ttf|otf)$/,
                use: {
                    loader: 'file-loader',
                    options: {
                        name: 'fonts/[name].[ext]',
                    },
                },
            },
        ],
    },
    resolve: {
        alias: {},
        modules: [
            path.resolve(__dirname, 'node_modules'),
        ],
        // unsafeCache: false
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'css/[name].css',
            // chunkFilename: "css/[id].css",
        }),
        new CopyPlugin({
            patterns: [
                {
                    from: 'node_modules/bootstrap/dist/css/bootstrap.min.css',
                    to: 'tp/bootstrap/css/',
                    toType: 'dir',
                },
                {
                    from: 'node_modules/bootstrap/dist/js/bootstrap.bundle.min.js',
                    to: 'tp/bootstrap/js/',
                    toType: 'dir',
                },
                {
                    from: 'node_modules/jquery/dist/jquery.min.js',
                    to: 'tp/jquery',
                    toType: 'dir',
                },
                {
                    from: 'node_modules/simple-line-icons/dist/styles/simple-line-icons.css',
                    to: 'tp/simple-line-icons/css/',
                    toType: 'dir',
                },
                {
                    from: 'node_modules/simple-line-icons/dist/fonts',
                    to: 'tp/simple-line-icons/fonts',
                    toType: 'dir',
                },
            ],
        }),
    ],
};
