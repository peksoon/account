/**
 * 날짜를 YYYY-MM-DD 형식의 문자열로 변환
 * @param {Date} date - 변환할 날짜 객체
 * @returns {string} YYYY-MM-DD 형식의 문자열
 */
export function formatDateToString(date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

/**
 * 날짜를 로컬 시간대의 ISO 문자열로 변환
 * @param {Date} date - 변환할 날짜 객체
 * @returns {string} 로컬 시간대의 ISO 문자열 (YYYY-MM-DD HH:mm:ss 형식)
 */
export function formatDateToLocalISOString(date) {
  const offset = date.getTimezoneOffset() * 60000;
  const DateOffset = new Date(date.getTime() - offset);
  return DateOffset.toISOString().slice(0, 19).replace('T', ' ');
}

/**
 * 문자열을 Date 객체로 변환
 * @param {string} dateString - 변환할 날짜 문자열
 * @returns {Date} Date 객체
 */
export function parseDate(dateString) {
  return new Date(dateString);
}

/**
 * 숫자를 통화 형식으로 포맷
 * @param {number} amount - 포맷할 숫자
 * @param {string} currency - 통화 코드 (기본값: 'KRW')
 * @returns {string} 포맷된 통화 문자열
 */
export function formatCurrency(amount, currency = 'KRW') {
  return new Intl.NumberFormat('ko-KR', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 0
  }).format(amount);
}

/**
 * 숫자를 천 단위 구분자로 포맷
 * @param {number} number - 포맷할 숫자
 * @returns {string} 포맷된 숫자 문자열
 */
export function formatNumber(number) {
  return new Intl.NumberFormat('ko-KR').format(number);
}

/**
 * 객체의 깊은 복사를 수행
 * @param {any} obj - 복사할 객체
 * @returns {any} 복사된 객체
 */
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }
  
  if (obj instanceof Date) {
    return new Date(obj.getTime());
  }
  
  if (obj instanceof Array) {
    return obj.map(item => deepClone(item));
  }
  
  if (typeof obj === 'object') {
    const clonedObj = {};
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        clonedObj[key] = deepClone(obj[key]);
      }
    }
    return clonedObj;
  }
}

/**
 * 디바운스 함수
 * @param {Function} func - 실행할 함수
 * @param {number} delay - 지연 시간 (밀리초)
 * @returns {Function} 디바운스된 함수
 */
export function debounce(func, delay) {
  let timeoutId;
  return function (...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func.apply(this, args), delay);
  };
}

/**
 * 스로틀 함수
 * @param {Function} func - 실행할 함수
 * @param {number} delay - 지연 시간 (밀리초)
 * @returns {Function} 스로틀된 함수
 */
export function throttle(func, delay) {
  let lastCall = 0;
  return function (...args) {
    const now = Date.now();
    if (now - lastCall >= delay) {
      lastCall = now;
      return func.apply(this, args);
    }
  };
}

/**
 * 배열에서 고유한 값들만 추출
 * @param {Array} array - 대상 배열
 * @param {string} key - 객체 배열인 경우 비교할 키
 * @returns {Array} 중복이 제거된 배열
 */
export function uniqueArray(array, key = null) {
  if (key) {
    const seen = new Set();
    return array.filter(item => {
      const value = item[key];
      if (seen.has(value)) {
        return false;
      }
      seen.add(value);
      return true;
    });
  }
  return [...new Set(array)];
}

/**
 * 문자열이 비어있는지 확인
 * @param {string} str - 확인할 문자열
 * @returns {boolean} 비어있으면 true
 */
export function isEmpty(str) {
  return !str || str.trim().length === 0;
}

/**
 * 이메일 형식 유효성 검사
 * @param {string} email - 검사할 이메일
 * @returns {boolean} 유효하면 true
 */
export function isValidEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * 로컬 스토리지에 데이터 저장
 * @param {string} key - 저장할 키
 * @param {any} value - 저장할 값
 */
export function setLocalStorage(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.error('로컬 스토리지 저장 오류:', error);
  }
}

/**
 * 로컬 스토리지에서 데이터 조회
 * @param {string} key - 조회할 키
 * @param {any} defaultValue - 기본값
 * @returns {any} 저장된 값 또는 기본값
 */
export function getLocalStorage(key, defaultValue = null) {
  try {
    const item = localStorage.getItem(key);
    return item ? JSON.parse(item) : defaultValue;
  } catch (error) {
    console.error('로컬 스토리지 조회 오류:', error);
    return defaultValue;
  }
}

/**
 * 로컬 스토리지에서 데이터 삭제
 * @param {string} key - 삭제할 키
 */
export function removeLocalStorage(key) {
  try {
    localStorage.removeItem(key);
  } catch (error) {
    console.error('로컬 스토리지 삭제 오류:', error);
  }
}