<template>
  <div class="modal-backdrop" @click.self="closePopup">
    <div class="modal-content">
      <!-- 헤더 -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <div class="flex items-center">
          <div class="w-12 h-12 rounded-lg flex items-center justify-center mr-4" :class="localEventDetail.type === 'out'
            ? 'bg-gradient-to-r from-red-500 to-red-600'
            : 'bg-gradient-to-r from-green-500 to-green-600'">
            <TrendingDown v-if="localEventDetail.type === 'out'" class="w-7 h-7 text-white" />
            <TrendingUp v-else class="w-7 h-7 text-white" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-gray-900">
              {{ isEditMode ? '데이터 편집' : (localEventDetail.type === 'out' ? '지출 상세' : '수입 상세') }}
            </h3>
            <p class="text-sm text-gray-500">
              {{ isEditMode ? '정보를 수정하세요' : '거래 내역을 확인하세요' }}
            </p>
          </div>
        </div>
        <el-button @click="closePopup" size="large" circle>×</el-button>
      </div>

      <!-- 내용 -->
      <div class="p-6">
        <!-- 편집 모드 -->
        <div v-if="isEditMode">
          <el-form ref="formRef" :model="localEventDetail" :rules="rules" label-position="top">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- 사용자 -->
              <el-form-item label="사용자" prop="user">
                <el-input v-model="localEventDetail.user" placeholder="사용자명을 입력하세요" size="large" />
              </el-form-item>

              <!-- 금액 -->
              <el-form-item label="금액" prop="money">
                <el-input v-model.number="localEventDetail.money" type="number" placeholder="금액을 입력하세요" size="large">
                  <template #suffix>원</template>
                </el-input>
              </el-form-item>
            </div>

            <!-- 카테고리 -->
            <el-form-item label="카테고리" prop="category_name" class="mb-6">
              <el-select v-model="localEventDetail.category_name" placeholder="카테고리를 선택하세요" size="large" class="w-full"
                filterable allow-create :loading="categoryStore.loading">
                <el-option v-for="category in availableCategories" :key="category" :label="category"
                  :value="category" />
              </el-select>
              <div v-if="categoryStore.loading" class="mt-1 text-xs text-gray-500">
                카테고리 목록을 불러오는 중...
              </div>
              <div v-else-if="availableCategories.length === 0" class="mt-1 text-xs text-orange-500">
                {{ localEventDetail.type === 'out' ? '지출' : '수입' }} 카테고리가 없습니다.
              </div>
              <div v-else class="mt-1 text-xs text-gray-500">
                {{ availableCategories.length }}개의 카테고리를 사용할 수 있습니다.
              </div>
            </el-form-item>

            <!-- 지출 전용 필드 -->
            <div v-if="localEventDetail.type === 'out'" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 키워드 -->
                <el-form-item label="키워드">
                  <el-input v-model="localEventDetail.keyword_name" :prefix-icon="Tag" placeholder="키워드를 입력하세요"
                    size="large" />
                </el-form-item>

                <!-- 결제 방법 -->
                <el-form-item label="결제 방법">
                  <el-select v-model="localEventDetail.payment_method_id" placeholder="결제 방법을 선택하세요" size="large"
                    class="w-full">
                    <el-option-group v-for="category in paymentMethodStore.activePaymentMethods" :key="category.id"
                      :label="category.name">
                      <el-option v-for="method in category.children" :key="method.id" :label="method.name"
                        :value="method.id">
                        <div class="flex items-center">
                          <CreditCard class="w-4 h-4 mr-2 text-purple-600" />
                          {{ method.name }}
                        </div>
                      </el-option>
                    </el-option-group>
                  </el-select>
                </el-form-item>
              </div>
            </div>

            <!-- 수입 전용 필드 -->
            <div v-if="localEventDetail.type === 'in'" class="mb-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 키워드 -->
                <el-form-item label="키워드">
                  <el-input v-model="localEventDetail.keyword_name" :prefix-icon="Tag" placeholder="키워드를 입력하세요"
                    size="large" />
                </el-form-item>

                <!-- 입금경로 -->
                <el-form-item label="입금경로">
                  <el-select v-model="localEventDetail.deposit_path" placeholder="입금경로를 선택하세요" size="large"
                    class="w-full" filterable allow-create>
                    <el-option v-for="path in depositPathOptions" :key="path.id" :label="path.name"
                      :value="path.name" />
                  </el-select>
                </el-form-item>
              </div>
            </div>

            <!-- 메모 -->
            <el-form-item label="메모" class="mb-6">
              <el-input v-model="localEventDetail.memo" type="textarea" :rows="3" placeholder="추가 메모를 입력하세요"
                resize="none" />
            </el-form-item>

            <!-- 날짜 -->
            <el-form-item label="날짜" prop="date" class="mb-6">
              <el-date-picker v-model="localEventDetail.date" type="datetime" placeholder="날짜와 시간을 선택하세요"
                format="YYYY-MM-DD HH:mm" value-format="YYYY-MM-DD HH:mm:ss" size="large" class="w-full" />
            </el-form-item>
          </el-form>
        </div>

        <!-- 상세 보기 모드 -->
        <div v-else class="space-y-4">
          <!-- 거래 정보 카드 -->
          <div class="bg-gray-50 rounded-xl p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- 금액 -->
              <div class="text-center md:text-left">
                <p class="text-sm text-gray-500 mb-1">거래 금액</p>
                <p class="text-3xl font-bold"
                  :class="localEventDetail.type === 'out' ? 'text-red-600' : 'text-green-600'">
                  {{ localEventDetail.type === 'out' ? '-' : '+' }}{{ formatMoney(localEventDetail.money) }}원
                </p>
              </div>

              <!-- 사용자 -->
              <div class="text-center md:text-left">
                <p class="text-sm text-gray-500 mb-1">사용자</p>
                <p class="text-xl font-semibold text-gray-900">{{ localEventDetail.user }}</p>
              </div>
            </div>
          </div>

          <!-- 세부 정보 -->
          <div class="space-y-4">
            <!-- 카테고리 -->
            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
                <Folder class="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">카테고리</p>
                <p class="font-semibold text-gray-900">{{ getCategoryName(localEventDetail.category_id) ||
                  localEventDetail.category || '-' }}</p>
              </div>
            </div>

            <!-- 키워드 (지출만) -->
            <div v-if="(localEventDetail.keyword_name || localEventDetail.keyword) && localEventDetail.type === 'out'"
              class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
                <Tag class="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">키워드</p>
                <p class="font-semibold text-gray-900">{{ localEventDetail.keyword_name || localEventDetail.keyword ||
                  '-' }}
                </p>
              </div>
            </div>

            <!-- 결제 방법 (지출만) -->
            <div v-if="localEventDetail.type === 'out'"
              class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center mr-4">
                <CreditCard class="w-5 h-5 text-orange-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">결제 방법</p>
                <p class="font-semibold text-gray-900">{{ getPaymentMethodName(localEventDetail.payment_method_id) ||
                  localEventDetail.payment || '-' }}</p>
              </div>
            </div>

            <!-- 키워드 (수입만) -->
            <div v-if="(localEventDetail.keyword_name || localEventDetail.keyword) && localEventDetail.type === 'in'"
              class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
                <Tag class="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">키워드</p>
                <p class="font-semibold text-gray-900">{{ localEventDetail.keyword_name || localEventDetail.keyword ||
                  '-' }}
                </p>
              </div>
            </div>

            <!-- 입금경로 (수입만) -->
            <div
              v-if="(localEventDetail.deposit_path_name || localEventDetail.deposit_path) && localEventDetail.type === 'in'"
              class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center mr-4">
                <CreditCard class="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">입금경로</p>
                <p class="font-semibold text-gray-900">{{ localEventDetail.deposit_path_name ||
                  localEventDetail.deposit_path ||
                  '-' }}</p>
              </div>
            </div>

            <!-- 메모 -->
            <div class="flex items-start p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center mr-4 flex-shrink-0">
                <FileText class="w-5 h-5 text-gray-600" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm text-gray-500">메모</p>
                <p class="font-semibold text-gray-900 break-words">
                  {{ localEventDetail.memo || '메모가 없습니다.' }}
                </p>
              </div>
            </div>

            <!-- 날짜 -->
            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
              <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center mr-4">
                <Calendar class="w-5 h-5 text-indigo-600" />
              </div>
              <div>
                <p class="text-sm text-gray-500">거래 일시</p>
                <p class="font-semibold text-gray-900">{{ formatDate(localEventDetail.date) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 하단 버튼 -->
      <div class="flex justify-between p-6 border-t border-gray-200 bg-gray-50">
        <div>
          <!-- 편집 모드가 아닐 때만 삭제 버튼 표시 -->
          <el-button v-if="!isEditMode" type="danger" @click="confirmDelete" size="large">
            삭제
          </el-button>
        </div>

        <div class="flex space-x-3">
          <!-- 편집 모드 -->
          <template v-if="isEditMode">
            <el-button @click="cancelEdit" size="large">
              취소
            </el-button>
            <el-button type="primary" @click="updateAccount" size="large" :loading="updating" :disabled="!isFormValid">
              <Save class="w-4 h-4 mr-2" />
              저장
            </el-button>
          </template>

          <!-- 보기 모드 -->
          <template v-else>
            <el-button @click="closePopup" size="large">
              닫기
            </el-button>
            <el-button type="primary" @click="openEditMode" size="large">
              편집
            </el-button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed, onMounted } from 'vue';
import {
  Tag,
  CreditCard,
  Save,
  TrendingUp,
  TrendingDown,
  Calendar,
  Folder,
  FileText
} from 'lucide-vue-next';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useCategoryStore } from '../stores/categoryStore';
import { usePaymentMethodStore } from '../stores/paymentMethodStore';
import { useDepositPathStore } from '../stores/depositPathStore';

export default {
  components: {
    Tag,
    CreditCard,
    Save,
    TrendingUp,
    TrendingDown,
    Calendar,
    Folder,
    FileText
  },
  props: {
    eventDetail: {
      type: Object,
      required: true
    }
  },
  emits: ['update', 'delete', 'close', 'edit'],
  setup(props, { emit }) {
    const formRef = ref(null);
    const updating = ref(false);
    const isEditMode = ref(false);

    // Stores
    const categoryStore = useCategoryStore();
    const paymentMethodStore = usePaymentMethodStore();
    const depositPathStore = useDepositPathStore();

    // 로컬 이벤트 상세 데이터
    const localEventDetail = ref({ ...props.eventDetail });

    // 폼 검증 규칙
    const rules = {
      user: [
        { required: true, message: '사용자명을 입력해주세요', trigger: 'blur' }
      ],
      money: [
        { required: true, message: '금액을 입력해주세요', trigger: 'blur' },
        { type: 'number', min: 1, message: '금액은 1원 이상이어야 합니다', trigger: 'blur' }
      ],
      category_name: [
        { required: true, message: '카테고리를 선택해주세요', trigger: 'change' }
      ],
      date: [
        { required: true, message: '날짜를 선택해주세요', trigger: 'change' }
      ]
    };

    // 사용 가능한 카테고리 목록 (백엔드에서 가져오기)
    const availableCategories = computed(() => {
      if (!categoryStore.categories || categoryStore.categories.length === 0) {
        return [];
      }

      const filtered = categoryStore.categories
        .filter(category => category.type === localEventDetail.value.type)
        .map(category => category.name);

      return filtered;
    });

    // 입금경로 옵션 목록
    const depositPathOptions = computed(() => {
      return depositPathStore.activeDepositPaths.map(path => ({
        id: path.id,
        name: path.name,
        value: path.name
      }));
    });

    // 폼 유효성 검사
    const isFormValid = computed(() => {
      return localEventDetail.value.user &&
        localEventDetail.value.money > 0 &&
        (localEventDetail.value.category_name || localEventDetail.value.category) &&
        localEventDetail.value.date;
    });

    // props 변경 감지
    watch(() => props.eventDetail, (newVal) => {
      localEventDetail.value = { ...newVal };

      // 카테고리 이름이 없고 ID가 있는 경우 이름 설정
      if (!localEventDetail.value.category_name) {
        if (localEventDetail.value.category_id) {
          const categoryName = getCategoryName(localEventDetail.value.category_id);
          if (categoryName) {
            localEventDetail.value.category_name = categoryName;
          }
        } else if (localEventDetail.value.category) {
          localEventDetail.value.category_name = localEventDetail.value.category;
        }
      }
    }, { deep: true });

    // 계정 업데이트
    const updateAccount = async () => {
      if (!formRef.value) return;

      try {
        // 폼 검증
        await formRef.value.validate();

        updating.value = true;

        // 데이터 정리
        const accountData = { ...localEventDetail.value };

        // category_name이 있으면 category 필드에도 설정
        if (accountData.category_name) {
          accountData.category = accountData.category_name;
        }

        console.log('DetailPopup에서 업데이트 이벤트 발생:', accountData);

        // 부모 컴포넌트에 업데이트 이벤트 전송
        emit('update', accountData);

        // 약간의 지연 후 편집 모드 종료 (비동기 업데이트 완료 대기)
        await new Promise(resolve => setTimeout(resolve, 100));

        // 컴포넌트가 여전히 마운트된 상태인지 확인
        if (isEditMode.value !== undefined) {
          isEditMode.value = false;
          ElMessage.success('데이터가 성공적으로 업데이트되었습니다.');
        }

      } catch (error) {
        console.error('Form validation failed:', error);
        if (isEditMode.value !== undefined) {
          ElMessage.error('입력 정보를 확인해주세요.');
        }
      } finally {
        if (updating.value !== undefined) {
          updating.value = false;
        }
      }
    };

    // 삭제 확인
    const confirmDelete = async () => {
      try {
        await ElMessageBox.confirm(
          '정말로 이 항목을 삭제하시겠습니까? 삭제된 데이터는 복구할 수 없습니다.',
          '삭제 확인',
          {
            confirmButtonText: '삭제',
            cancelButtonText: '취소',
            type: 'warning',
            confirmButtonClass: 'el-button--danger'
          }
        );

        emit('delete', localEventDetail.value);
        ElMessage.success('데이터가 삭제되었습니다.');

      } catch {
        // 사용자가 취소함
      }
    };

    // 팝업 닫기
    const closePopup = () => {
      emit('close');
    };

    // 편집 모드 열기
    const openEditMode = async () => {
      try {
        // 항상 최신 카테고리 데이터를 로드 (캐시 무시)
        await categoryStore.fetchCategories();

        isEditMode.value = true;
        // 편집 모드 시작 시 현재 데이터를 다시 로드하여 최신 상태 반영
        localEventDetail.value = { ...props.eventDetail };

        // 카테고리 이름 설정 (ID가 있는 경우 이름으로 변환)
        if (localEventDetail.value.category_id && !localEventDetail.value.category_name) {
          const categoryName = getCategoryName(localEventDetail.value.category_id);
          if (categoryName) {
            localEventDetail.value.category_name = categoryName;
          }
        }
        // category 필드가 있고 category_name이 없는 경우
        else if (localEventDetail.value.category && !localEventDetail.value.category_name) {
          localEventDetail.value.category_name = localEventDetail.value.category;
        }
      } catch (error) {
        console.error('편집 모드 시작 중 오류:', error);
        ElMessage.error('편집 모드를 시작하는데 실패했습니다.');
      }
    };

    // 편집 취소
    const cancelEdit = () => {
      isEditMode.value = false;
      // 원본 데이터로 되돌리기
      localEventDetail.value = { ...props.eventDetail };
    };

    // 금액 포맷팅
    const formatMoney = (amount) => {
      if (!amount) return '0';
      return new Intl.NumberFormat('ko-KR').format(amount);
    };

    // 날짜 포맷팅 (시간 제거)
    const formatDate = (dateString) => {
      if (!dateString) return '';

      try {
        const date = new Date(dateString);
        const options = {
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          timeZone: 'Asia/Seoul'
        };

        return date.toLocaleDateString('ko-KR', options);
      } catch {
        return dateString;
      }
    };

    // 카테고리 이름 가져오기
    const getCategoryName = (categoryId) => {
      if (!categoryId) return '';
      const category = categoryStore.categories.find(cat => cat.id === categoryId);
      return category ? category.name : '';
    };

    // 결제수단 이름 가져오기
    const getPaymentMethodName = (paymentMethodId) => {
      if (!paymentMethodId) return '';
      const paymentMethod = paymentMethodStore.flatPaymentMethods.find(pm => pm.id === paymentMethodId);
      return paymentMethod ? paymentMethod.name : '';
    };

    // 컴포넌트 마운트 시 데이터 로드
    onMounted(async () => {
      try {
        // 모든 필요한 데이터를 병렬로 로드
        await Promise.all([
          categoryStore.fetchCategories(),
          paymentMethodStore.fetchPaymentMethods(),
          depositPathStore.fetchDepositPaths()
        ]);

        // 초기 카테고리 이름 설정
        if (!localEventDetail.value.category_name) {
          if (localEventDetail.value.category_id) {
            const categoryName = getCategoryName(localEventDetail.value.category_id);
            if (categoryName) {
              localEventDetail.value.category_name = categoryName;
            }
          } else if (localEventDetail.value.category) {
            localEventDetail.value.category_name = localEventDetail.value.category;
          }
        }
      } catch (error) {
        console.error('데이터 로드 오류:', error);
        ElMessage.error('데이터를 불러오는데 실패했습니다.');
      }
    });

    return {
      formRef,
      updating,
      isEditMode,
      localEventDetail,
      rules,
      availableCategories,
      depositPathOptions,
      isFormValid,
      updateAccount,
      confirmDelete,
      closePopup,
      openEditMode,
      cancelEdit,
      formatMoney,
      formatDate,
      getCategoryName,
      getPaymentMethodName,

      // Stores
      paymentMethodStore,
      categoryStore,

      // 아이콘들
      Tag,
      CreditCard,
      Save,
      TrendingUp,
      TrendingDown,
      Calendar,
      Folder,
      FileText
    };
  }
};
</script>

<style scoped>
/* 커스텀 입력 필드 스타일 */
:deep(.el-form-item__label) {
  @apply font-semibold text-gray-700 text-sm;
}

:deep(.el-input__wrapper) {
  @apply border-gray-200 rounded-lg;
  transition: all 0.2s ease;
}

:deep(.el-input__wrapper:hover) {
  @apply border-gray-300;
}

:deep(.el-input__wrapper.is-focus) {
  @apply border-blue-500 ring-2 ring-blue-100;
}

:deep(.el-select .el-input__wrapper) {
  @apply border-gray-200 rounded-lg;
}

:deep(.el-textarea__inner) {
  @apply border-gray-200 rounded-lg;
  resize: none;
}

:deep(.el-date-editor.el-input) {
  @apply w-full;
}

/* 버튼 커스텀 스타일 */
:deep(.el-button--primary) {
  @apply bg-gradient-to-r from-blue-500 to-blue-600 border-blue-500;
}

:deep(.el-button--primary:hover) {
  @apply from-blue-600 to-blue-700 border-blue-600;
}

:deep(.el-button--danger) {
  @apply bg-gradient-to-r from-red-500 to-red-600 border-red-500;
}

:deep(.el-button--danger:hover) {
  @apply from-red-600 to-red-700 border-red-600;
}

/* 모달 애니메이션 */
.modal-backdrop {
  animation: fadeIn 0.3s ease-out;
}

.modal-content {
  animation: slideIn 0.3s ease-out;
  max-width: 600px;
  width: 95%;
  max-height: 90vh;
  overflow-y: auto;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 정보 카드 호버 효과 */
.bg-white.border:hover {
  @apply border-gray-300 shadow-sm;
  transition: all 0.2s ease;
}

/* 금액 표시 애니메이션 */
.text-3xl.font-bold {
  animation: numberPop 0.4s ease-out;
}

@keyframes numberPop {
  0% {
    transform: scale(0.8);
    opacity: 0;
  }

  50% {
    transform: scale(1.05);
  }

  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 아이콘 색상 애니메이션 */
.w-5.h-5 {
  transition: all 0.2s ease;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
  .modal-content {
    max-width: none;
    width: 100%;
    margin-left: 1rem;
    margin-right: 1rem;
    max-height: 95vh;
  }

  .grid.grid-cols-1.md\\:grid-cols-2 {
    grid-template-columns: repeat(1, minmax(0, 1fr));
  }

  .text-3xl {
    font-size: 1.5rem;
    line-height: 2rem;
  }

  .text-xl {
    font-size: 1.125rem;
    line-height: 1.75rem;
  }

  .p-6 {
    padding: 1rem;
  }

  .space-y-4>*+* {
    margin-top: 1rem;
  }
}

/* 스크롤바 커스텀 */
.modal-content::-webkit-scrollbar {
  width: 6px;
}

.modal-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

.modal-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 10px;
}

.modal-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 반짝이는 효과 */
@keyframes shimmer {
  0% {
    background-position: -468px 0;
  }

  100% {
    background-position: 468px 0;
  }
}

.bg-gray-50 {
  position: relative;
  overflow: hidden;
}

.bg-gray-50::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
      transparent,
      rgba(255, 255, 255, 0.4),
      transparent);
  animation: shimmer 1.5s ease-in-out infinite;
}

/* 로딩 상태 스타일 */
:deep(.el-button.is-loading) {
  @apply pointer-events-none;
}

/* 부드러운 전환 효과 */
.space-y-4>* {
  transition: all 0.2s ease;
}

.space-y-4>*:hover {
  transform: translateY(-1px);
}
</style>