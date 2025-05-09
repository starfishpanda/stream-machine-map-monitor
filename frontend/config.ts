export const getWebSocketURL = () => {
  const isProduction = import.meta.env.MODE === 'production';
  
  if (isProduction) {
    // In production inside Docker, use the host name without specifying port
    const url = `ws://${window.location.host}/ws/machine`;
    // console.log('Production WebSocket URL:', url);
    return url;
  } else {
    // In development, connect directly to the WebSocket server
    console.log('Development WebSocket URL: ws://localhost:3001/machine');
    return 'ws://localhost:3001/machine';
  }
};

export const getMapsApiKey = () => {
  return import.meta.env.VITE_GOOGLE_MAPS_API_KEY || '';
}