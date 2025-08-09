<template>
    <div class="keyword-autocomplete">
        <el-autocomplete v-model="inputValue" :fetch-suggestions="fetchSuggestions" :placeholder="placeholder"
            :disabled="disabled" :size="size" class="w-full" clearable @select="handleSelect" @input="handleInput"
            @clear="handleClear">
            <template #default="{ item }">
                <div class="flex items-center justify-between w-full">
                    <div class="flex items-center">
                        <Tag class="w-4 h-4 text-gray-500 mr-2" />
                        <span class="text-gray-900">{{ item.name }}</span>
                    </div>
                    <div class="flex items-center text-xs text-gray-500">
                        <span>{{ item.usage_count }}회</span>
                    </div>
                </div>
            </template>

            <template #prefix>
                <Tag class="w-4 h-4 text-gray-500" />
            </template>
        </el-autocomplete>

        <!-- 키워드 관리 버튼 -->
        <div v-if="showManageButton" class="mt-2">
            <el-button size="small" text @click="openKeywordManager">
                키워드 관리
            </el-button>
        </div>
    </div>
</template>

<script>
import { ref, watch } from 'vue';
import { Tag } from 'lucide-vue-next';
import { useKeywordStore } from '../stores/keywordStore';
// import { ElMessage } from 'element-plus';

export default {
    name: 'KeywordAutocomplete',
    components: {
        Tag
    },
    props: {
        modelValue: {
            type: String,
            default: ''
        },
        categoryId: {
            type: [Number, null],
            default: null
        },
        placeholder: {
            type: String,
            default: '키워드를 입력하세요'
        },
        disabled: {
            type: Boolean,
            default: false
        },
        size: {
            type: String,
            default: 'default'
        },
        showManageButton: {
            type: Boolean,
            default: false
        }
    },
    emits: ['update:modelValue', 'select', 'input', 'clear'],
    setup(props, { emit }) {
        const keywordStore = useKeywordStore();

        const inputValue = ref(props.modelValue);
        const suggestions = ref([]);
        const loading = ref(false);

        // 입력값 변경 감지
        watch(() => props.modelValue, (newValue) => {
            inputValue.value = newValue;
        });

        // 카테고리 변경 시 입력값 초기화
        watch(() => props.categoryId, () => {
            inputValue.value = '';
            emit('update:modelValue', '');
            suggestions.value = [];
        });

        // 자동완성 제안 가져오기
        const fetchSuggestions = async (queryString, callback) => {
            if (!props.categoryId) {
                callback([]);
                return;
            }

            try {
                loading.value = true;
                const results = await keywordStore.getKeywordSuggestions(
                    props.categoryId,
                    queryString,
                    10
                );

                // 결과를 autocomplete 형식에 맞게 변환
                const formattedResults = results.map(item => ({
                    value: item.name,
                    name: item.name,
                    id: item.id,
                    usage_count: item.usage_count
                }));

                callback(formattedResults);
                suggestions.value = formattedResults;

            } catch (error) {
                console.error('키워드 제안 조회 오류:', error);
                callback([]);
            } finally {
                loading.value = false;
            }
        };

        // 제안 선택 핸들러
        const handleSelect = (item) => {
            inputValue.value = item.value;
            emit('update:modelValue', item.value);
            emit('select', {
                id: item.id,
                name: item.name,
                usage_count: item.usage_count
            });
        };

        // 입력 핸들러
        const handleInput = (value) => {
            inputValue.value = value;
            emit('update:modelValue', value);
            emit('input', value);
        };

        // 클리어 핸들러
        const handleClear = () => {
            inputValue.value = '';
            emit('update:modelValue', '');
            emit('clear');
        };

        // 키워드 관리자 열기
        const openKeywordManager = () => {
            // 부모 컴포넌트에서 키워드 관리 모달을 열도록 이벤트 발생
            emit('open-keyword-manager', props.categoryId);
        };

        return {
            inputValue,
            suggestions,
            loading,

            fetchSuggestions,
            handleSelect,
            handleInput,
            handleClear,
            openKeywordManager,

            // 아이콘들
            Tag
        };
    }
};
</script>

<style scoped>
/* Autocomplete 스타일 커스텀 */
:deep(.el-autocomplete) {
    @apply w-full;
}

:deep(.el-input__wrapper) {
    @apply transition-all duration-200;
}

:deep(.el-input__wrapper:hover) {
    @apply border-gray-300;
}

:deep(.el-input__wrapper.is-focus) {
    @apply border-primary-500 ring-2 ring-primary-100;
}

/* 자동완성 드롭다운 스타일 */
:deep(.el-autocomplete-suggestion) {
    @apply border border-gray-200 rounded-lg shadow-lg;
}

:deep(.el-autocomplete-suggestion__list) {
    @apply py-2;
}

:deep(.el-autocomplete-suggestion__item) {
    @apply px-4 py-3 hover:bg-gray-50 transition-colors duration-200;
}

:deep(.el-autocomplete-suggestion__item.highlighted) {
    @apply bg-primary-50 text-primary-700;
}

/* 키워드 관리 버튼 */
:deep(.el-button--text) {
    @apply text-gray-600 hover:text-primary-600 transition-colors duration-200;
}

/* 로딩 상태 */
.keyword-autocomplete :deep(.el-input__suffix) {
    @apply flex items-center;
}
</style>
