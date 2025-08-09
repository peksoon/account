<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
    <div class="bg-white rounded-xl shadow-2xl w-full max-w-4xl max-h-[90vh] overflow-hidden" @click.stop>
      <!-- 헤더 -->
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold">결제수단 관리</h2>
          <button @click="$emit('close')" class="text-white hover:text-gray-200 transition-colors">
            <span class="text-2xl">×</span>
          </button>
        </div>
      </div>

      <!-- 컨텐츠 -->
      <div class="p-6 max-h-[calc(90vh-120px)] overflow-y-auto">
        <!-- 새 결제수단 추가 -->
        <div class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h3 class="text-lg font-semibold mb-4">새 결제수단 추가</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <el-select v-model="newPaymentMethod.parent_id" placeholder="카테고리 선택" size="large" clearable>
              <el-option v-for="parent in parentPaymentMethods" :key="parent.id" :label="parent.name"
                :value="parent.id" />
            </el-select>
            <el-input v-model="newPaymentMethod.name" placeholder="결제수단 이름을 입력하세요" size="large" />
            <el-button type="primary" @click="addPaymentMethod" :loading="loading" size="large" class="w-full">
              추가
            </el-button>
          </div>
        </div>

        <!-- 결제수단 목록 -->
        <div>
          <h3 class="text-lg font-semibold mb-4">결제수단 목록</h3>

          <div v-if="loading" class="text-center py-8">
            <el-icon class="is-loading text-blue-500 text-xl">
              <span>⟳</span>
            </el-icon>
            <p class="mt-2 text-gray-500">데이터를 불러오는 중...</p>
          </div>

          <div v-else-if="paymentMethodStore.paymentMethods.length === 0" class="text-center py-8">
            <p class="text-gray-500">등록된 결제수단이 없습니다.</p>
          </div>

          <div v-else class="space-y-4">
            <!-- 카테고리별 그룹 -->
            <div v-for="category in paymentMethodStore.paymentMethods" :key="category.id"
              class="border border-gray-200 rounded-lg">
              <!-- 카테고리 헤더 -->
              <div class="bg-gray-100 p-4 rounded-t-lg border-b">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                      <CreditCard class="w-5 h-5 text-purple-600" />
                    </div>
                    <div>
                      <p class="font-semibold text-gray-900">{{ category.name }}</p>
                      <p class="text-sm text-gray-500">카테고리 · {{ category.children?.length || 0 }}개 항목</p>
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <el-button size="small" @click="editPaymentMethod(category)">
                      수정
                    </el-button>
                    <el-button size="small" @click="confirmDelete(category)" type="danger"
                      :disabled="category.children && category.children.length > 0">
                      삭제
                    </el-button>
                  </div>
                </div>
              </div>

              <!-- 세부 결제수단 목록 -->
              <div v-if="category.children && category.children.length > 0" class="p-4 space-y-2">
                <div v-for="child in category.children" :key="child.id"
                  class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                  <div class="flex items-center gap-3">
                    <div class="w-6 h-6 bg-white rounded border flex items-center justify-center ml-4">
                      <span class="text-xs text-gray-600">•</span>
                    </div>
                    <div>
                      <p class="font-medium text-gray-800">{{ child.name }}</p>
                      <p class="text-xs text-gray-500">ID: {{ child.id }}</p>
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <el-button size="small" @click="editPaymentMethod(child)">
                      수정
                    </el-button>
                    <el-button size="small" @click="confirmDelete(child)" type="danger">
                      삭제
                    </el-button>
                  </div>
                </div>
              </div>

              <!-- 빈 카테고리 메시지 -->
              <div v-else class="p-4 text-center text-gray-500">
                <p class="text-sm">세부 결제수단이 없습니다.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 편집 다이얼로그 -->
      <el-dialog v-model="showEditDialog" title="결제수단 수정" width="400px">
        <el-form :model="editingMethod" label-position="top">
          <el-form-item label="이름">
            <el-input v-model="editingMethod.name" placeholder="결제수단 이름을 입력하세요" />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="showEditDialog = false">취소</el-button>
            <el-button type="primary" @click="updatePaymentMethod" :loading="loading">
              저장
            </el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { CreditCard } from 'lucide-vue-next';
import { usePaymentMethodStore } from '@/stores/paymentMethodStore';

export default {
  name: 'PaymentMethodManager',
  components: {
    CreditCard
  },
  emits: ['close'],
  setup() {
    const paymentMethodStore = usePaymentMethodStore();
    const loading = ref(false);
    const showEditDialog = ref(false);

    const newPaymentMethod = ref({
      name: '',
      parent_id: null
    });

    const editingMethod = ref({
      id: null,
      name: ''
    });

    // 부모 결제수단 (카테고리) 목록
    const parentPaymentMethods = computed(() => {
      return paymentMethodStore.parentPaymentMethods || [];
    });

    // 결제수단 목록 로드
    const loadPaymentMethods = async () => {
      try {
        loading.value = true;
        await paymentMethodStore.fetchPaymentMethods();
      } catch (error) {
        console.error('결제수단 목록 로드 오류:', error);
        ElMessage.error('결제수단 목록을 불러오는데 실패했습니다.');
      } finally {
        loading.value = false;
      }
    };

    // 새 결제수단 추가
    const addPaymentMethod = async () => {
      if (!newPaymentMethod.value.name.trim()) {
        ElMessage.warning('결제수단 이름을 입력해주세요.');
        return;
      }

      try {
        loading.value = true;
        await paymentMethodStore.createPaymentMethod({
          name: newPaymentMethod.value.name.trim(),
          parent_id: newPaymentMethod.value.parent_id
        });

        // 폼 초기화
        newPaymentMethod.value = {
          name: '',
          parent_id: null
        };

        ElMessage.success('결제수단이 추가되었습니다.');
      } catch (error) {
        console.error('결제수단 추가 오류:', error);
        ElMessage.error('결제수단 추가에 실패했습니다.');
      } finally {
        loading.value = false;
      }
    };

    // 결제수단 수정
    const editPaymentMethod = (method) => {
      editingMethod.value = {
        id: method.id,
        name: method.name
      };
      showEditDialog.value = true;
    };

    // 결제수단 업데이트
    const updatePaymentMethod = async () => {
      if (!editingMethod.value.name.trim()) {
        ElMessage.warning('결제수단 이름을 입력해주세요.');
        return;
      }

      try {
        loading.value = true;
        await paymentMethodStore.updatePaymentMethod(editingMethod.value.id, {
          name: editingMethod.value.name.trim()
        });

        showEditDialog.value = false;
        ElMessage.success('결제수단이 수정되었습니다.');
      } catch (error) {
        console.error('결제수단 수정 오류:', error);
        ElMessage.error('결제수단 수정에 실패했습니다.');
      } finally {
        loading.value = false;
      }
    };

    // 결제수단 삭제 확인
    const confirmDelete = async (method) => {
      try {
        await ElMessageBox.confirm(
          `'${method.name}' 결제수단을 삭제하시겠습니까?`,
          '삭제 확인',
          {
            confirmButtonText: '삭제',
            cancelButtonText: '취소',
            type: 'warning'
          }
        );

        await deletePaymentMethod(method.id);
      } catch (error) {
        // 사용자가 취소한 경우
        if (error === 'cancel') return;
        console.error('삭제 확인 오류:', error);
      }
    };

    // 결제수단 삭제
    const deletePaymentMethod = async (id) => {
      try {
        loading.value = true;
        await paymentMethodStore.deletePaymentMethod(id);
        ElMessage.success('결제수단이 삭제되었습니다.');
      } catch (error) {
        console.error('결제수단 삭제 오류:', error);
        ElMessage.error('결제수단 삭제에 실패했습니다.');
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      loadPaymentMethods();
    });

    return {
      paymentMethodStore,
      loading,
      showEditDialog,
      newPaymentMethod,
      editingMethod,
      parentPaymentMethods,
      addPaymentMethod,
      editPaymentMethod,
      updatePaymentMethod,
      confirmDelete
    };
  }
};
</script>

<style scoped>
/* 아이콘 색상 수정 */
.lucide-credit-card {
  color: #8b5cf6 !important;
}
</style>