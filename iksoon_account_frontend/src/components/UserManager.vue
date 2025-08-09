<template>
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
        <div class="bg-white rounded-xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden" @click.stop>
            <!-- 헤더 -->
            <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-6">
                <div class="flex justify-between items-center">
                    <h2 class="text-xl font-bold">사용자 관리</h2>
                    <button @click="$emit('close')" class="text-white hover:text-gray-200 transition-colors">
                        <span class="text-2xl">×</span>
                    </button>
                </div>
            </div>

            <!-- 컨텐츠 -->
            <div class="p-6 max-h-[calc(90vh-120px)] overflow-y-auto">
                <!-- 새 사용자 추가 -->
                <div class="mb-6 p-4 bg-gray-50 rounded-lg">
                    <h3 class="text-lg font-semibold mb-4">새 사용자 추가</h3>
                    <el-form ref="addFormRef" :model="newUser" :rules="rules" label-position="top">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <el-form-item label="사용자 이름" prop="name">
                                <el-input v-model="newUser.name" placeholder="사용자 이름을 입력하세요" size="large" />
                            </el-form-item>
                            <el-form-item label="이메일 (선택사항)" prop="email">
                                <el-input v-model="newUser.email" placeholder="이메일 주소를 입력하세요" size="large" />
                            </el-form-item>
                        </div>
                        <div class="flex justify-end">
                            <el-button type="primary" @click="addUser" :loading="loading" size="large">
                                추가
                            </el-button>
                        </div>
                    </el-form>
                </div>

                <!-- 사용자 목록 -->
                <div>
                    <h3 class="text-lg font-semibold mb-4">사용자 목록</h3>
                    <div v-if="userStore.loading" class="text-center py-8">
                        <div class="inline-flex items-center">
                            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
                            <span class="ml-3 text-gray-600">로딩 중...</span>
                        </div>
                    </div>
                    <div v-else-if="userStore.users.length === 0" class="text-center py-8 text-gray-500">
                        등록된 사용자가 없습니다.
                    </div>
                    <div v-else class="space-y-3">
                        <div v-for="user in userStore.users" :key="user.id"
                            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                            <div class="flex-1">
                                <div class="flex items-center">
                                    <span class="font-semibold text-gray-900">{{ user.name }}</span>
                                    <span v-if="user.email" class="ml-2 text-sm text-gray-500">({{ user.email }})</span>
                                </div>
                                <div class="text-xs text-gray-400 mt-1">
                                    ID: {{ user.id }} | 생성일: {{ formatDate(user.created_at) }}
                                </div>
                            </div>
                            <div class="flex space-x-2">
                                <el-button size="small" @click="editUser(user)" type="primary" :disabled="loading">
                                    수정
                                </el-button>
                                <el-button size="small" @click="checkAndDeleteUser(user)" type="danger"
                                    :disabled="loading">
                                    삭제
                                </el-button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 수정 모달 -->
        <div v-if="editingUser" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-60">
            <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6" @click.stop>
                <h3 class="text-lg font-semibold mb-4">사용자 수정</h3>
                <el-form ref="editFormRef" :model="editingUser" :rules="rules" label-position="top">
                    <el-form-item label="사용자 이름" prop="name">
                        <el-input v-model="editingUser.name" placeholder="사용자 이름을 입력하세요" size="large" />
                    </el-form-item>
                    <el-form-item label="이메일 (선택사항)" prop="email">
                        <el-input v-model="editingUser.email" placeholder="이메일 주소를 입력하세요" size="large" />
                    </el-form-item>
                    <div class="flex justify-end space-x-3">
                        <el-button @click="cancelEdit" :disabled="loading">취소</el-button>
                        <el-button type="primary" @click="updateUser" :loading="loading">수정</el-button>
                    </div>
                </el-form>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useUserStore } from '../stores/userStore';

export default {
    emits: ['close'],
    setup() {
        const userStore = useUserStore();
        const addFormRef = ref(null);
        const editFormRef = ref(null);
        const loading = ref(false);
        const editingUser = ref(null);

        const newUser = ref({
            name: '',
            email: ''
        });

        const rules = {
            name: [
                { required: true, message: '사용자 이름을 입력해주세요', trigger: 'blur' },
                { min: 1, max: 100, message: '사용자 이름은 1-100자 사이여야 합니다', trigger: 'blur' }
            ],
            email: [
                { type: 'email', message: '올바른 이메일 형식을 입력해주세요', trigger: 'blur' }
            ]
        };

        // 사용자 추가
        const addUser = async () => {
            try {
                await addFormRef.value.validate();
                loading.value = true;

                await userStore.createUser({
                    name: newUser.value.name,
                    email: newUser.value.email || ''
                });

                ElMessage.success('사용자가 성공적으로 추가되었습니다');

                // 폼 초기화
                newUser.value = { name: '', email: '' };
                addFormRef.value.resetFields();
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('사용자 추가 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        // 사용자 수정
        const editUser = (user) => {
            editingUser.value = { ...user };
        };

        const updateUser = async () => {
            try {
                await editFormRef.value.validate();
                loading.value = true;

                await userStore.updateUser({
                    id: editingUser.value.id,
                    name: editingUser.value.name,
                    email: editingUser.value.email || ''
                });

                ElMessage.success('사용자가 성공적으로 수정되었습니다');
                editingUser.value = null;
            } catch (error) {
                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('사용자 수정 중 오류가 발생했습니다');
                }
            } finally {
                loading.value = false;
            }
        };

        const cancelEdit = () => {
            editingUser.value = null;
        };

        // 사용자 삭제
        const checkAndDeleteUser = async (user) => {
            try {
                loading.value = true;

                // 사용 여부 확인
                const inUse = await userStore.checkUserUsage(user.id);

                if (inUse) {
                    const result = await ElMessageBox.confirm(
                        `'${user.name}' 사용자는 현재 사용 중입니다. 정말 삭제하시겠습니까?`,
                        '사용자 삭제 확인',
                        {
                            confirmButtonText: '삭제',
                            cancelButtonText: '취소',
                            type: 'warning'
                        }
                    );

                    if (result === 'confirm') {
                        await userStore.forceDeleteUser(user.id);
                        ElMessage.success('사용자가 삭제되었습니다');
                    }
                } else {
                    const result = await ElMessageBox.confirm(
                        `'${user.name}' 사용자를 삭제하시겠습니까?`,
                        '사용자 삭제 확인',
                        {
                            confirmButtonText: '삭제',
                            cancelButtonText: '취소',
                            type: 'warning'
                        }
                    );

                    if (result === 'confirm') {
                        await userStore.deleteUser(user.id);
                        ElMessage.success('사용자가 삭제되었습니다');
                    }
                }
            } catch (error) {
                if (error === 'cancel') {
                    // 사용자가 취소한 경우
                    return;
                }

                if (error.response?.data?.message) {
                    ElMessage.error(error.response.data.message);
                } else {
                    ElMessage.error('사용자 삭제 중 오류가 발생했습니다');
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

        // 컴포넌트 마운트 시 사용자 목록 로드
        onMounted(async () => {
            try {
                await userStore.fetchUsers();
            } catch (error) {
                ElMessage.error('사용자 목록을 불러올 수 없습니다');
            }
        });

        return {
            userStore,
            addFormRef,
            editFormRef,
            loading,
            editingUser,
            newUser,
            rules,
            addUser,
            editUser,
            updateUser,
            cancelEdit,
            checkAndDeleteUser,
            formatDate
        };
    }
};
</script>

<style scoped>
.z-60 {
    z-index: 60;
}

/* 아이콘 색상 수정 */
.lucide-user {
    color: #6366f1 !important;
}
</style>
