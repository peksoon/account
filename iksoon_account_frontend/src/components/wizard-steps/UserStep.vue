<template>
    <div class="user-step">
        <div class="step-header">
            <h2 class="step-title">사용자를 선택하세요</h2>
            <p class="step-description">{{ modelValue.type === 'out' ? '지출' : '수입' }}한 사용자를 선택해주세요</p>
        </div>

        <div class="user-selection-container">
            <!-- 사용자 검색/선택 -->
            <div class="user-input-wrapper">
                <el-select
                    ref="userSelectRef"
                    v-model="selectedUser"
                    placeholder="사용자를 선택하거나 입력하세요"
                    size="large"
                    class="user-select"
                    :class="{ 'error': hasError }"
                    filterable
                    allow-create
                    default-first-option
                    @change="handleUserChange"
                    @visible-change="handleDropdownToggle"
                >
                    <el-option
                        v-for="user in userOptions"
                        :key="user.id"
                        :label="user.label"
                        :value="user.value"
                        class="user-option"
                    >
                        <div class="flex items-center justify-between w-full">
                            <div class="flex items-center">
                                <User class="w-4 h-4 text-gray-500 mr-2" />
                                <span>{{ user.label }}</span>
                            </div>
                            <div v-if="user.isRecent" class="text-xs text-blue-600 bg-blue-100 px-2 py-1 rounded">
                                최근
                            </div>
                        </div>
                    </el-option>
                </el-select>
                
                <div v-if="hasError" class="error-message">
                    {{ errorMessage }}
                </div>
            </div>

            <!-- 빠른 선택 (최근 사용자) -->
            <div class="quick-users" v-if="recentUsers.length > 0 && !selectedUser">
                <p class="quick-users-label">최근 사용자</p>
                <div class="quick-users-grid">
                    <button
                        v-for="user in recentUsers"
                        :key="user.value"
                        @click="selectUser(user.value)"
                        class="quick-user-btn"
                    >
                        <User class="w-4 h-4 mr-2" />
                        {{ user.label }}
                    </button>
                </div>
            </div>

            <!-- 사용자 관리 링크 -->
            <div class="user-management">
                <el-button 
                    text 
                    size="small" 
                    @click="openUserManager"
                    class="user-manage-btn"
                >
                    <Settings class="w-4 h-4 mr-1" />
                    사용자 관리
                </el-button>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import { User, Settings } from 'lucide-vue-next';
import { useUserStore } from '../../stores/userStore';

export default {
    name: 'UserStep',
    components: {
        User,
        Settings
    },
    props: {
        modelValue: {
            type: Object,
            required: true
        },
        errors: {
            type: Object,
            default: () => ({})
        }
    },
    emits: ['update:modelValue', 'next', 'auto-advance', 'validate', 'open-user-manager'],
    setup(props, { emit }) {
        const userStore = useUserStore();
        const userSelectRef = ref(null);
        const selectedUser = ref('');

        // 에러 상태
        const hasError = computed(() => !!props.errors.user);
        const errorMessage = computed(() => props.errors.user || '');

        // 사용자 옵션
        const userOptions = computed(() => {
            return userStore.getUserOptions();
        });

        // 최근 사용자 (최대 4명)
        const recentUsers = computed(() => {
            return userOptions.value
                .filter(user => user.isRecent)
                .slice(0, 4);
        });

        // 이벤트 핸들러들
        const handleUserChange = (value) => {
            if (value) {
                selectedUser.value = value;
                updateModelValue(value);
                emit('validate', 'user', true, '');
                
                // 선택 즉시 자동 진행
                setTimeout(() => {
                    emit('auto-advance', 100);
                }, 50);
            } else {
                emit('validate', 'user', false, '사용자를 선택해주세요');
            }
        };

        const handleDropdownToggle = (visible) => {
            if (!visible && selectedUser.value) {
                // 드롭다운이 닫힐 때 값이 있으면 자동 진행
                setTimeout(() => {
                    emit('auto-advance', 200);
                }, 100);
            }
        };

        const selectUser = (value) => {
            selectedUser.value = value;
            handleUserChange(value);
        };

        const updateModelValue = (value) => {
            const updated = { ...props.modelValue, user: value };
            emit('update:modelValue', updated);
        };

        const openUserManager = () => {
            emit('open-user-manager');
        };

        // 초기값 설정
        watch(() => props.modelValue.user, (newValue) => {
            if (newValue !== selectedUser.value) {
                selectedUser.value = newValue || '';
            }
        }, { immediate: true });

        // 컴포넌트 마운트 시 포커스
        onMounted(() => {
            nextTick(() => {
                if (userSelectRef.value) {
                    userSelectRef.value.focus();
                }
            });
        });

        return {
            userSelectRef,
            selectedUser,
            hasError,
            errorMessage,
            userOptions,
            recentUsers,
            handleUserChange,
            handleDropdownToggle,
            selectUser,
            openUserManager
        };
    }
}
</script>

<style scoped>
.user-step {
    @apply max-w-md mx-auto;
}

.step-header {
    @apply text-center mb-8;
}

.step-title {
    @apply text-2xl font-bold text-gray-900 mb-2;
}

.step-description {
    @apply text-gray-600;
}

.user-selection-container {
    @apply space-y-6;
}

.user-input-wrapper {
    @apply relative;
}

.user-select {
    @apply w-full;
}

:deep(.user-select .el-input__inner) {
    @apply h-14 text-lg border-2;
}

:deep(.user-select.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.user-select .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

.user-option {
    @apply py-2;
}

.error-message {
    @apply mt-2 text-sm text-red-600;
}

.quick-users {
    @apply text-center;
}

.quick-users-label {
    @apply text-sm text-gray-600 mb-3;
}

.quick-users-grid {
    @apply grid grid-cols-2 gap-2;
}

.quick-user-btn {
    @apply flex items-center justify-center px-4 py-3 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-colors duration-200;
}

.quick-user-btn:hover {
    @apply bg-gray-50 border-gray-400;
}

.quick-user-btn:active {
    @apply bg-gray-100;
}

.user-management {
    @apply text-center;
}

.user-manage-btn {
    @apply text-blue-600;
}

.user-manage-btn:hover {
    @apply text-blue-700;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .step-title {
        @apply text-xl;
    }
    
    .step-description {
        @apply text-sm;
    }
    
    :deep(.user-select .el-input__inner) {
        @apply h-12 text-base;
        font-size: 16px; /* iOS zoom 방지 */
    }
    
    .quick-users-grid {
        @apply grid-cols-1 gap-3;
    }
    
    .quick-user-btn {
        @apply px-4 py-4 text-sm;
        min-height: 44px; /* iOS 권장 터치 영역 */
    }
}

/* 접근성 */
.user-select:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.quick-user-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.user-manage-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* Element Plus 드롭다운 최적화 */
:deep(.el-select-dropdown) {
    @apply max-h-60;
}

:deep(.el-select-dropdown__item) {
    @apply py-3;
}

/* 애니메이션 */
.quick-users {
    animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>