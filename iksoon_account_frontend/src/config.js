// API 설정 파일 - 환경별 구분
export const API_CONFIG = {
  DEV_URL: 'http://localhost:8080',      // 개발: 로컬 백엔드
  PROD_URL: 'http://133.186.153.179:8080' // 운영: 운영 서버 백엔드
};

// 현재 환경에 맞는 API URL 반환
export const getApiBaseUrl = () => {
  // NODE_ENV로 환경 구분
  if (process.env.NODE_ENV === 'production') {
    console.log('🚀 Production mode: Using', API_CONFIG.PROD_URL);
    return API_CONFIG.PROD_URL;
  } else {
    console.log('🛠️ Development mode: Using', API_CONFIG.DEV_URL);
    return API_CONFIG.DEV_URL;
  }
};

export default API_CONFIG;
