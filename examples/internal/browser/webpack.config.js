const path = require('path');

module.exports = {
  mode: 'development',
  resolve: {
    alias: {
      'swagger-client$': path.resolve(__dirname, 'node_modules/swagger-client/dist/swagger-client.browser.js')
    },
    fallback: {
      "http": false,
      "https": false,
      "stream": false,
      "url": false,
      "buffer": false,
      "util": false
    }
  }
};
