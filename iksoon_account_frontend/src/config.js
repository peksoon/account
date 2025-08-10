// API 설정 파일 - 환경변수 기반
export const getApiBaseUrl = () => {
  // 환경변수에서 API URL 가져오기 (환경별 .env 파일에서 설정)
  const apiUrl = process.env.VUE_APP_BACKEND_API_BASE_URL;
  
  if (process.env.NODE_ENV === 'production') {
    console.log('🚀 Production mode: Using', apiUrl);
  } else {
    console.log('🛠️ Development mode: Using', apiUrl);
  }
  
  return apiUrl || 'http://localhost:8080'; // 기본값
};

export default { getApiBaseUrl };
