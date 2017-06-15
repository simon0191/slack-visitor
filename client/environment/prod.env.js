module.exports = {
  ENV: '"production"',
  NODE_ENV: '"production"',
  DEBUG_MODE: true,
  API_URL: `"${process.env.API_URL || 'http://localhost:8000'}"`,
  API_HOST_NAME: `"${process.env.API_HOST_NAME || 'http://localhost'}"`,
  API_PORT: `${process.env.API_PORT || 8000}`,
}
