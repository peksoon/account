<template>
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
        <div class="bg-white rounded-xl shadow-2xl w-full max-w-4xl max-h-[90vh] overflow-hidden" @click.stop>
            <!-- 헤더 -->
            <div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-6">
                <div class="flex justify-between items-center">
                    <h2 class="text-xl font-bold">키워드 관리</h2>
                    <button @click="$emit('close')" class="text-white hover:text-gray-200 transition-colors">
                        <span class="text-2xl">×</span>
                    </button>
                </div>
            </div>

            <!-- 컨텐츠 -->
            <div class="p-6 max-h-[calc(90vh-120px)] overflow-y-auto">
                <!-- 카테고리 선택 -->
                <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                    <h3 class="text-lg font-semibold mb-4">카테고리 선택</h3>
                    <el-select v-model="selectedCategoryId" placeholder="카테고리를 선택하세요" size="large" class="w-full"
                        @change="handleCategoryChange">
                        <el-option label="지출 카테고리" disabled />
                        <el-option v-for="category in outCategories" :key="category.id" :label="category.name"
                            :value="category.id" />
                        <el-option label="수입 카테고리" disabled />
                        <el-option v-for="category in inCategories" :key="category.id" :label="category.name"
                            :value="category.id" />
                    </el-select>
                </div>

                <!-- 새 키워드 추가 -->
                <div v-if="selectedCategoryId" class="mb-6 p-4 bg-gray-50 rounded-lg">
                    <h3 class="text-lg font-semibold mb-4">새 키워드 추가</h3>
                    <el-form ref="addFormRef" :model="newKeyword" :rules="rules" label-position="top">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <el-form-item label="키워드 이름" prop="name">
                                <el-input v-model="newKeyword.name" placeholder="키워드 이름을 입력하세요" size="large" />
                            </el-form-item>
                            <el-form-item>
                                <div style="height: 32px;"></div>
                                <el-button type="primary" @click="addKeyword" :loading="loading" size="large">
                                    추가
                                </el-button>
                            </el-form-item>
                        </div>
                    </el-form>
                </div>

                <!-- 키워드 목록 -->
                <div v-if="selectedCategoryId">
                    <h3 class="text-lg font-semibold mb-4">키워드 목록</h3>

                    <div v-if="keywordStore.loading" class="text-center py-8">
                        <div class="inline-flex items-center">
                            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-purple-500"></div>
                            <span class="ml-3 text-gray-600">로딩 중...</span>
                        </div>
                    </div>
                    <div v-else-if="currentKeywords.length === 0" class="text-center py-8 text-gray-500">
                        선택한 카테고리에 등록된 키워드가 없습니다.
                    </div>
                    <div v-else class="space-y-3">
                        <div v-for="keyword in sortedKeywords" :key="keyword.id"
                            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                            <div class="flex-1">
                                <div class="flex items-center">
                                    <span class="font-semibold text-gray-900">{{ keyword.name }}</span>
                                    <span class="ml-2 px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-600">
                                        사용횟수: {{ keyword.use_count || 0 }}
                                    </span>
                                </div>
                                <div class="text-xs text-gray-400 mt-1">
                                    ID: {{ keyword.id }} | 생성일: {{ formatDate(keyword.created_at) }}
                                </div>
                            </div>
                            <div class="flex space-x-2">
                                <el-button size="small" @click="editKeyword(keyword)" type="primary"
                                    :disabled="loading">
                                    수정
                                </el-button>
                                <el-button size="small" @click="checkAndDeleteKeyword(keyword)" type="danger"
                                    :disabled="loading">
                                    삭제
                                </el-button>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 카테고리 미선택 안내 -->
                <div v-else class="text-center py-12 text-gray-500">
                    <div class="w-16 h-16 mx-auto mb-4 text-gray-300">
                        <svg fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd"
                                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                                clip-rule="evenodd" />
                        </svg>
                    </div>
                    <p class="text-lg font-medium">카테고리를 선택해주세요</p>
                    <p class="text-sm">선택한 카테고리의 키워드를 관리할 수 있습니다.</p>
                </div>
            </div>
        </div>

        <!-- 수정 모달 -->
        <div v-if="editingKeyword" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-60">
            <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6" @click.stop>
                <h3 class="text-lg font-semibold mb-4">키워드 수정</h3>
                <el-form ref="editFormRef" :model="editingKeyword" :rules="rules" label-position="top">
                    <el-form-item label="키워드 이름" prop="name">
                        <el-input v-model="editingKeyword.name" placeholder="키워드 이름을 입력하세요" size="large" />
                    </el-form-item>
                    <div class="flex justify-end space-x-3">
                        <el-button @click="cancelEdit" :disabled="loading">취소</el-button>
                        <el-button type="primary" @click="updateKeyword" :loading="loading">수정</el-button>
                    </div>
                </el-form>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useKeywordStore } from '../stores/keywordStore';
import { useCategoryStore } from '../stores/categoryStore';

export default {
    props: {
        categoryId: {
            type: Number,
            default: null
        }
    },
    emits: ['close'],
    setup(props) {
        const keywordStore = useKeywordStore();
        const categoryStore = useCategoryStore();
        const addFormRef = ref(null);
        const editFormRef = ref(null);
        const loading = ref(false);
        const editingKeyword = ref(null);
        const selectedCategoryId = ref(props.categoryId);

        const newKeyword = ref({
            name: ''
        });

        const rules = {
            name: [
                { required: true, message: '키워드 이름을 입력해주세요', trigger: 'blur' },
                { min: 1, max: 100, message: '키워드 이름은 1-100자 사이여야 합니다', trigger: 'blur' }
            ]
        };

        // 카테고리 목록
        const outCategories = computed(() => categoryStore.outCategories);
        const inCategories = computed(() => categoryStore.inCategories);

        // 현재 선택된 카테고리의 키워드들
        const currentKeywords = computed(() => {
            if (!selectedCategoryId.value) return [];
            return keywordStore.keywords.filter(k => k.category_id === selectedCategoryId.value);
        });

        // 사용횟수순으로 정렬된 키워드
        const sortedKeywords = computed(() => {
            return [...currentKeywords.value].sort((a, b) => (b.use_count || 0) - (a.use_count || 0));
        });

        // 카테고리 변경 핸들러
        const handleCategoryChange = async () => {
            if (selectedCategoryId.value) {
                await keywordStore.fetchKeywords(selectedCategoryId.value);
            }
        };

        // 키워드 추가
        const addKeyword = async () => {
            try {
                await addFormRef.value.validate();
                loading.value = true;

                await keywordStore.createKeyword({
                    category_id: selectedCategoryId.value,
                    name: newKeyword.value.name
                });

                ElMessage.success('키워드가 성공적으로 추가되었습니다');

                // 폼 초기화
                newKeyword.value = { name: '' };
                addFormRef.value.resetFields();
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('키워드 추가 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 키워드 수정
        const editKeyword = (keyword) => {
            editingKeyword.value = { ...keyword };
        };

        const updateKeyword = async () => {
            try {
                await editFormRef.value.validate();
                loading.value = true;

                await keywordStore.updateKeyword(editingKeyword.value.id, {
                    name: editingKeyword.value.name
                });

                ElMessage.success('키워드가 성공적으로 수정되었습니다');
                editingKeyword.value = null;
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('키워드 수정 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        const cancelEdit = () => {
            editingKeyword.value = null;
        };

        // 키워드 삭제
        const checkAndDeleteKeyword = async (keyword) => {
            try {
                const result = await ElMessageBox.confirm(
                    `'${keyword.name}' 키워드를 삭제하시겠습니까?`,
                    '키워드 삭제 확인',
                    {
                        confirmButtonText: '삭제',
                        cancelButtonText: '취소',
                        type: 'warning'
                    }
                );

                if (result === 'confirm') {
                    loading.value = true;
                    await keywordStore.deleteKeyword(keyword.id);
                    ElMessage.success('키워드가 삭제되었습니다');
                }
            } catch (error) {
                if (error === 'cancel') {
                    return;
                }

                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('키워드 삭제 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 날짜 포맷팅
        const formatDate = (dateString) => {
            if (!dateString) return '';
            return new Date(dateString).toLocaleString('ko-KR', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit'
            });
        };

        // props.categoryId 변경 감지
        watch(() => props.categoryId, (newId) => {
            selectedCategoryId.value = newId;
            if (newId) {
                handleCategoryChange();
            }
        });

        // 컴포넌트 마운트 시 데이터 로드
        onMounted(async () => {
            try {
                await categoryStore.fetchCategories();
                if (selectedCategoryId.value) {
                    await keywordStore.fetchKeywords(selectedCategoryId.value);
                }
            } catch (error) {
                ElMessage.error('데이터를 불러올 수 없습니다');
            }
        });

        return {
            keywordStore,
            addFormRef,
            editFormRef,
            loading,
            editingKeyword,
            selectedCategoryId,
            newKeyword,
            rules,
            outCategories,
            inCategories,
            currentKeywords,
            sortedKeywords,
            handleCategoryChange,
            addKeyword,
            editKeyword,
            updateKeyword,
            cancelEdit,
            checkAndDeleteKeyword,
            formatDate
        };
    }
};
</script>

<style scoped>
.z-60 {
    z-index: 60;
}
</style>

