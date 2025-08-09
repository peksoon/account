<template>
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
        <div class="bg-white rounded-xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden" @click.stop>
            <!-- 헤더 -->
            <div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-6">
                <div class="flex justify-between items-center">
                    <h2 class="text-xl font-bold">입금경로 관리</h2>
                    <button @click="$emit('close')" class="text-white hover:text-gray-200 transition-colors">
                        <span class="text-2xl">×</span>
                    </button>
                </div>
            </div>

            <!-- 컨텐츠 -->
            <div class="p-6 max-h-[calc(90vh-120px)] overflow-y-auto">
                <!-- 새 입금경로 추가 -->
                <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                    <h3 class="text-lg font-semibold mb-4">새 입금경로 추가</h3>
                    <el-form ref="addFormRef" :model="newDepositPath" :rules="rules" label-position="top">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <el-form-item label="입금경로 이름" prop="name">
                                <el-input v-model="newDepositPath.name" placeholder="입금경로 이름을 입력하세요" size="large" />
                            </el-form-item>
                            <div class="flex items-end">
                                <el-button type="primary" @click="addDepositPath" :loading="loading" size="large"
                                    class="w-full">
                                    추가
                                </el-button>
                            </div>
                        </div>
                    </el-form>
                </div>

                <!-- 입금경로 목록 -->
                <div>
                    <h3 class="text-lg font-semibold mb-4">입금경로 목록</h3>

                    <div v-if="loading" class="text-center py-8">
                        <el-icon class="is-loading text-blue-500 text-xl">
                            <span>⟳</span>
                        </el-icon>
                        <p class="mt-2 text-gray-500">데이터를 불러오는 중...</p>
                    </div>

                    <div v-else-if="depositPathStore.depositPaths.length === 0" class="text-center py-8">
                        <p class="text-gray-500">등록된 입금경로가 없습니다.</p>
                    </div>

                    <div v-else class="space-y-3">
                        <div v-for="path in depositPathStore.depositPaths" :key="path.id"
                            class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow">
                            <div class="flex items-center space-x-3">
                                <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                                    <ArrowDownToLine class="w-5 h-5 text-green-600" />
                                </div>
                                <div>
                                    <h4 class="font-semibold text-gray-900">
                                        {{ editingPath && editingPath.id === path.id ? '' : path.name }}
                                    </h4>
                                    <p class="text-sm text-gray-500">ID: {{ path.id }}</p>
                                </div>
                            </div>

                            <!-- 편집 모드 -->
                            <div v-if="editingPath && editingPath.id === path.id" class="flex items-center space-x-2">
                                <el-input v-model="editingPath.name" size="default" placeholder="입금경로 이름"
                                    style="width: 200px" />
                                <el-button size="small" type="primary" @click="updateDepositPath" :loading="loading">
                                    저장
                                </el-button>
                                <el-button size="small" @click="cancelEdit">
                                    취소
                                </el-button>
                            </div>

                            <!-- 일반 모드 -->
                            <div v-else class="flex items-center space-x-2">
                                <el-button size="small" @click="editDepositPath(path)">
                                    수정
                                </el-button>
                                <el-button size="small" @click="checkAndDeleteDepositPath(path)" type="danger"
                                    :disabled="loading">
                                    삭제
                                </el-button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ArrowDownToLine } from 'lucide-vue-next';
import { useDepositPathStore } from '../stores/depositPathStore';

export default {
    components: {
        ArrowDownToLine
    },
    emits: ['close'],
    setup() {
        const depositPathStore = useDepositPathStore();
        const addFormRef = ref(null);
        const loading = ref(false);
        const editingPath = ref(null);

        // 새 입금경로 폼
        const newDepositPath = ref({
            name: ''
        });

        // 폼 검증 규칙
        const rules = {
            name: [
                { required: true, message: '입금경로 이름을 입력해주세요', trigger: 'blur' }
            ]
        };

        // 입금경로 추가
        const addDepositPath = async () => {
            try {
                await addFormRef.value.validate();
                loading.value = true;

                await depositPathStore.createDepositPath({
                    name: newDepositPath.value.name
                });

                ElMessage.success('입금경로가 성공적으로 추가되었습니다');
                newDepositPath.value.name = '';
                addFormRef.value.resetFields();
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('입금경로 추가 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 입금경로 편집 시작
        const editDepositPath = (path) => {
            editingPath.value = { ...path };
        };

        // 입금경로 수정
        const updateDepositPath = async () => {
            try {
                if (!editingPath.value.name.trim()) {
                    ElMessage.error('입금경로 이름을 입력해주세요');
                    return;
                }

                loading.value = true;

                await depositPathStore.updateDepositPath(editingPath.value.id, {
                    name: editingPath.value.name
                });

                ElMessage.success('입금경로가 성공적으로 수정되었습니다');
                editingPath.value = null;
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('입금경로 수정 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 편집 취소
        const cancelEdit = () => {
            editingPath.value = null;
        };

        // 입금경로 삭제 확인 및 삭제
        const checkAndDeleteDepositPath = async (path) => {
            try {
                const result = await ElMessageBox.confirm(
                    `'${path.name}' 입금경로를 삭제하시겠습니까?`,
                    '입금경로 삭제 확인',
                    {
                        confirmButtonText: '삭제',
                        cancelButtonText: '취소',
                        type: 'warning'
                    }
                );

                if (result === 'confirm') {
                    loading.value = true;

                    try {
                        await depositPathStore.deleteDepositPath(path.id);
                        ElMessage.success('입금경로가 삭제되었습니다');
                    } catch (error) {
                        if (error.response?.status === 409) {
                            // 사용 중인 입금경로인 경우 강제 삭제 확인
                            const forceResult = await ElMessageBox.confirm(
                                `'${path.name}' 입금경로는 현재 사용 중입니다. 정말 삭제하시겠습니까?`,
                                '입금경로 강제 삭제 확인',
                                {
                                    confirmButtonText: '삭제',
                                    cancelButtonText: '취소',
                                    type: 'warning'
                                }
                            );

                            if (forceResult === 'confirm') {
                                await depositPathStore.forceDeleteDepositPath(path.id);
                                ElMessage.success('입금경로가 삭제되었습니다');
                            }
                        } else {
                            throw error;
                        }
                    }
                }
            } catch (error) {
                if (error === 'cancel') {
                    return; // 사용자가 취소한 경우
                }

                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('입금경로 삭제 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 컴포넌트 마운트 시 데이터 로드
        onMounted(async () => {
            try {
                await depositPathStore.fetchDepositPaths();
            } catch (error) {
                console.error('입금경로 로드 오류:', error);
                ElMessage.error('입금경로 목록을 불러올 수 없습니다');
            }
        });

        return {
            depositPathStore,
            addFormRef,
            loading,
            editingPath,
            newDepositPath,
            rules,
            addDepositPath,
            editDepositPath,
            updateDepositPath,
            cancelEdit,
            checkAndDeleteDepositPath
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
    @apply border-green-500 ring-2 ring-green-100;
}

:deep(.el-select .el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
}

/* 버튼 커스텀 스타일 */
:deep(.el-button--primary) {
    @apply bg-gradient-to-r from-green-500 to-green-600 border-green-500;
}

:deep(.el-button--primary:hover) {
    @apply from-green-600 to-green-700 border-green-600;
}

:deep(.el-button--danger) {
    @apply bg-gradient-to-r from-red-500 to-red-600 border-red-500;
}

:deep(.el-button--danger:hover) {
    @apply from-red-600 to-red-700 border-red-600;
}

/* 아이콘 색상 수정 */
.lucide-arrow-down-to-line {
    color: #10b981 !important;
}
</style>
