// API 설정 파일 - 프록시 기반
export const getApiBaseUrl = () => {
  // 프로덕션 환경에서는 Nginx 프록시를 통해 /api로 접근
  if (process.env.NODE_ENV === 'production') {
    console.log('🚀 Production mode: Using Nginx proxy /api');
    return '/api'; // Nginx 프록시를 통한 상대경로
  } else {
    // 개발 환경에서는 직접 localhost 백엔드 접근
    console.log('🛠️ Development mode: Using localhost backend');
    return 'http://localhost:8080'; // 개발용 직접 접근
  }
};

export default { getApiBaseUrl };
