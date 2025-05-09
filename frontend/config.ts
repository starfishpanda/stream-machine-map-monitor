export const getWebSocketURL = () => {
  const isProduction = import.meta.env.MODE === 'production';
  const host = isProduction ? window.location.hostname : 'localhost'
  const port = isProduction ? (window.location.port || '80') : '3001'

  return `ws://${host}:${port}/machine`;
};

export const getMapsApiKey = () => {
  return import.meta.env.VITE_GOOGLE_MAPS_API_KEY || '';
}