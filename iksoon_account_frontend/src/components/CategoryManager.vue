<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
    <div class="bg-white rounded-xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden" @click.stop>
      <!-- 헤더 -->
      <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold">카테고리 관리</h2>
          <button @click="$emit('close')" class="text-white hover:text-gray-200 transition-colors">
            <span class="text-2xl">×</span>
          </button>
        </div>
      </div>

      <!-- 컨텐츠 -->
      <div class="p-6 max-h-[calc(90vh-120px)] overflow-y-auto">
        <!-- 새 카테고리 추가 -->
        <div class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h3 class="text-lg font-semibold mb-4">새 카테고리 추가</h3>
          <el-form ref="addFormRef" :model="newCategory" :rules="rules" label-position="top">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <el-form-item label="카테고리 이름" prop="name">
                <el-input v-model="newCategory.name" placeholder="카테고리 이름을 입력하세요" size="large" />
              </el-form-item>
              <el-form-item label="카테고리 타입" prop="type">
                <el-select v-model="newCategory.type" placeholder="타입을 선택하세요" size="large" class="w-full">
                  <el-option label="지출" value="out" />
                  <el-option label="수입" value="in" />
                </el-select>
              </el-form-item>
            </div>
            <div class="flex justify-end">
              <el-button type="primary" @click="addCategory" :loading="loading" size="large">
                추가
              </el-button>
            </div>
          </el-form>
        </div>

        <!-- 카테고리 목록 -->
        <div>
          <h3 class="text-lg font-semibold mb-4">카테고리 목록</h3>

          <!-- 탭 -->
          <el-tabs v-model="activeTab" class="mb-4">
            <el-tab-pane label="지출 카테고리" name="out">
              <CategoryList :categories="outCategories" :loading="categoryStore.loading" @edit="editCategory"
                @delete="checkAndDeleteCategory" type="out" />
            </el-tab-pane>
            <el-tab-pane label="수입 카테고리" name="in">
              <CategoryList :categories="inCategories" :loading="categoryStore.loading" @edit="editCategory"
                @delete="checkAndDeleteCategory" type="in" />
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <!-- 수정 모달 -->
    <div v-if="editingCategory" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-60">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6" @click.stop>
        <h3 class="text-lg font-semibold mb-4">카테고리 수정</h3>
        <el-form ref="editFormRef" :model="editingCategory" :rules="rules" label-position="top">
          <el-form-item label="카테고리 이름" prop="name">
            <el-input v-model="editingCategory.name" placeholder="카테고리 이름을 입력하세요" size="large" />
          </el-form-item>
          <el-form-item label="카테고리 타입" prop="type">
            <el-select v-model="editingCategory.type" placeholder="타입을 선택하세요" size="large" class="w-full">
              <el-option label="지출" value="out" />
              <el-option label="수입" value="in" />
            </el-select>
          </el-form-item>
          <div class="flex justify-end space-x-3">
            <el-button @click="cancelEdit" :disabled="loading">취소</el-button>
            <el-button type="primary" @click="updateCategory" :loading="loading">수정</el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useCategoryStore } from '../stores/categoryStore';
import CategoryList from './CategoryList.vue';

export default {
  components: {
    CategoryList
  },
  emits: ['close'],
  setup() {
    const categoryStore = useCategoryStore();
    const addFormRef = ref(null);
    const editFormRef = ref(null);
    const loading = ref(false);
    const editingCategory = ref(null);
    const activeTab = ref('out');

    const newCategory = ref({
      name: '',
      type: 'out'
    });

    const rules = {
      name: [
        { required: true, message: '카테고리 이름을 입력해주세요', trigger: 'blur' },
        { min: 1, max: 255, message: '카테고리 이름은 1-255자 사이여야 합니다', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '카테고리 타입을 선택해주세요', trigger: 'change' }
      ]
    };

    // 필터링된 카테고리 목록
    const outCategories = computed(() =>
      categoryStore.categories.filter(cat => cat.type === 'out')
    );

    const inCategories = computed(() =>
      categoryStore.categories.filter(cat => cat.type === 'in')
    );

    // 카테고리 추가
    const addCategory = async () => {
      try {
        await addFormRef.value.validate();
        loading.value = true;

        await categoryStore.createCategory({
          name: newCategory.value.name,
          type: newCategory.value.type
        });

        ElMessage.success('카테고리가 성공적으로 추가되었습니다');

        // 폼 초기화
        newCategory.value = { name: '', type: 'out' };
        addFormRef.value.resetFields();
      } catch (error) {
        if (error.response?.data?.message) {
          ElMessage.error(error.response.data.message);
        } else {
          ElMessage.error('카테고리 추가 중 오류가 발생했습니다');
        }
      } finally {
        loading.value = false;
      }
    };

    // 카테고리 수정
    const editCategory = (category) => {
      editingCategory.value = { ...category };
    };

    const updateCategory = async () => {
      try {
        await editFormRef.value.validate();
        loading.value = true;

        await categoryStore.updateCategory({
          id: editingCategory.value.id,
          name: editingCategory.value.name,
          type: editingCategory.value.type
        });

        ElMessage.success('카테고리가 성공적으로 수정되었습니다');
        editingCategory.value = null;
      } catch (error) {
        if (error.response?.data?.message) {
          ElMessage.error(error.response.data.message);
        } else {
          ElMessage.error('카테고리 수정 중 오류가 발생했습니다');
        }
      } finally {
        loading.value = false;
      }
    };

    const cancelEdit = () => {
      editingCategory.value = null;
    };

    // 카테고리 삭제
    const checkAndDeleteCategory = async (category) => {
      try {
        const result = await ElMessageBox.confirm(
          `'${category.name}' 카테고리를 삭제하시겠습니까?`,
          '카테고리 삭제 확인',
          {
            confirmButtonText: '삭제',
            cancelButtonText: '취소',
            type: 'warning'
          }
        );

        if (result === 'confirm') {
          loading.value = true;
          await categoryStore.deleteCategory(category.id);
          ElMessage.success('카테고리가 삭제되었습니다');
        }
      } catch (error) {
        if (error === 'cancel') {
          return;
        }



        // 카테고리가 사용 중인 경우 강제 삭제 옵션 제공
        if (error.response?.data?.message &&
          error.response.data.message.includes('사용하는 데이터가 존재')) {
          try {
            const forceResult = await ElMessageBox.confirm(
              `'${category.name}' 카테고리가 사용 중입니다.\n강제로 삭제하시겠습니까? (기존 가계부 데이터는 유지됩니다)`,
              '강제 삭제 확인',
              {
                confirmButtonText: '강제 삭제',
                cancelButtonText: '취소',
                type: 'error'
              }
            );

            if (forceResult === 'confirm') {
              loading.value = true;
              await categoryStore.forceDeleteCategory(category.id);
              ElMessage.success('카테고리가 강제 삭제되었습니다');
            }
          } catch (forceError) {
            if (forceError === 'cancel') {
              return;
            }
            console.error('카테고리 강제 삭제 오류:', forceError);
            ElMessage.error('카테고리 강제 삭제 중 오류가 발생했습니다');
          }
        } else {
          if (error.response?.data?.message) {
            ElMessage.error(error.response.data.message);
          } else {
            ElMessage.error('카테고리 삭제 중 오류가 발생했습니다');
          }
        }
      } finally {
        loading.value = false;
      }
    };

    // 컴포넌트 마운트 시 카테고리 목록 로드
    onMounted(async () => {
      try {
        await categoryStore.fetchCategories();
      } catch (error) {
        console.error('카테고리 로드 오류:', error);
        ElMessage.error('카테고리 목록을 불러올 수 없습니다');
      }
    });

    return {
      categoryStore,
      addFormRef,
      editFormRef,
      loading,
      editingCategory,
      activeTab,
      newCategory,
      rules,
      outCategories,
      inCategories,
      addCategory,
      editCategory,
      updateCategory,
      cancelEdit,
      checkAndDeleteCategory
    };
  }
};
</script>

<style scoped>
.z-60 {
  z-index: 60;
}

/* 아이콘 색상 수정 */
.lucide-folder {
  color: #f59e0b !important;
}

.lucide-trending-up {
  color: #10b981 !important;
}

.lucide-trending-down {
  color: #ef4444 !important;
}
</style>