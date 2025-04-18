const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    createProxyMiddleware([
      '/sign_in',
      '/sign_up',
      '/get_tasks',
      '/logout',
      '/create_task',
      '/update_priority',
      '/update_task'
    ],
    {
      target: 'http://158.160.24.141:80',
      changeOrigin: true,
    })
  );
};
