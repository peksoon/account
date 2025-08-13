<template>
  <!-- 메인 모달 -->
  <BaseModal 
    title="사용자 관리" 
    theme="blue" 
    size="2xl" 
    @close="$emit('close')"
  >
    <!-- 새 사용자 추가 폼 -->
    <div class="mb-6 p-4 bg-gray-50 rounded-lg">
      <h3 class="text-lg font-semibold mb-4">새 사용자 추가</h3>
      <el-form ref="addFormRef" :model="newUser" :rules="validationRules" label-position="top">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <el-form-item label="사용자 이름" prop="name">
            <el-input 
              v-model="newUser.name" 
              placeholder="사용자 이름을 입력하세요" 
              size="large" 
              maxlength="50"
              show-word-limit
            />
          </el-form-item>
          <el-form-item label="이메일 (선택사항)" prop="email">
            <el-input 
              v-model="newUser.email" 
              placeholder="이메일 주소를 입력하세요" 
              size="large"
              type="email"
            />
          </el-form-item>
        </div>
        <div class="flex justify-end">
          <el-button 
            type="primary" 
            @click="handleAddUser" 
            :loading="loading" 
            size="large"
          >
            추가
          </el-button>
        </div>
      </el-form>
    </div>

    <!-- 사용자 목록 -->
    <div>
      <h3 class="text-lg font-semibold mb-4">사용자 목록</h3>
      
      <!-- 로딩 상태 -->
      <LoadingState v-if="userStore.loading" message="사용자 목록을 불러오는 중..." />
      
      <!-- 빈 상태 -->
      <EmptyState 
        v-else-if="userStore.users.length === 0"
        icon-type="users"
        message="등록된 사용자가 없습니다."
        description="새 사용자를 추가해보세요."
      />
      
      <!-- 사용자 목록 -->
      <div v-else class="space-y-3">
        <UserListItem
          v-for="user in userStore.users" 
          :key="user.id"
          :user="user"
          :loading="loading"
          @edit="handleEditUser"
          @delete="handleDeleteUser"
        />
      </div>
    </div>
  </BaseModal>

  <!-- 수정 모달 -->
  <BaseModal
    v-if="editingUser"
    title="사용자 수정"
    theme="blue"
    size="md"
    @close="handleCancelEdit"
  >
    <el-form ref="editFormRef" :model="editingUser" :rules="validationRules" label-position="top">
      <el-form-item label="사용자 이름" prop="name">
        <el-input 
          v-model="editingUser.name" 
          placeholder="사용자 이름을 입력하세요" 
          size="large"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>
      <el-form-item label="이메일 (선택사항)" prop="email">
        <el-input 
          v-model="editingUser.email" 
          placeholder="이메일 주소를 입력하세요" 
          size="large"
          type="email"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <el-button @click="handleCancelEdit" :disabled="loading">취소</el-button>
        <el-button type="primary" @click="handleUpdateUser" :loading="loading">수정</el-button>
      </div>
    </template>
  </BaseModal>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useUserStore } from '../stores/userStore';
import { deepClone, isEmpty, isValidEmail } from '../utils';

// 컴포넌트들 import
import BaseModal from './common/BaseModal.vue';
import LoadingState from './common/LoadingState.vue';
import EmptyState from './common/EmptyState.vue';
import UserListItem from './common/UserListItem.vue';

/**
 * 사용자 관리 컴포넌트
 * 사용자 목록 조회, 생성, 수정, 삭제 기능을 제공합니다.
 */

// Props & Emits
defineEmits(['close']);

// Store & Refs
const userStore = useUserStore();
const addFormRef = ref(null);
const editFormRef = ref(null);
const loading = ref(false);
const editingUser = ref(null);

// Form Data
const newUser = ref({
  name: '',
  email: ''
});

// Validation Rules
const validationRules = computed(() => ({
  name: [
    { required: true, message: '사용자 이름을 입력해주세요', trigger: 'blur' },
    { min: 1, max: 50, message: '사용자 이름은 1-50자 사이여야 합니다', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (isEmpty(value.trim())) {
          callback(new Error('사용자 이름은 공백만으로 구성될 수 없습니다'));
        } else {
          callback();
        }
      }, 
      trigger: 'blur' 
    }
  ],
  email: [
    { 
      validator: (rule, value, callback) => {
        if (value && !isValidEmail(value)) {
          callback(new Error('올바른 이메일 형식을 입력해주세요'));
        } else {
          callback();
        }
      }, 
      trigger: 'blur' 
    }
  ]
}));

/**
 * 새 사용자 추가 처리
 */
const handleAddUser = async () => {
  try {
    await addFormRef.value.validate();
    loading.value = true;

    await userStore.createUser({
      name: newUser.value.name.trim(),
      email: newUser.value.email?.trim() || ''
    });

    // 폼 초기화
    resetAddForm();
  } catch (error) {
    console.error('사용자 추가 오류:', error);
  } finally {
    loading.value = false;
  }
};

/**
 * 사용자 수정 시작
 * @param {Object} user - 수정할 사용자 객체
 */
const handleEditUser = (user) => {
  editingUser.value = deepClone(user);
};

/**
 * 사용자 수정 처리
 */
const handleUpdateUser = async () => {
  try {
    await editFormRef.value.validate();
    loading.value = true;

    await userStore.updateUser({
      id: editingUser.value.id,
      name: editingUser.value.name.trim(),
      email: editingUser.value.email?.trim() || ''
    });

    editingUser.value = null;
  } catch (error) {
    console.error('사용자 수정 오류:', error);
  } finally {
    loading.value = false;
  }
};

/**
 * 사용자 수정 취소
 */
const handleCancelEdit = () => {
  editingUser.value = null;
};

/**
 * 사용자 삭제 처리
 * @param {Object} user - 삭제할 사용자 객체
 */
const handleDeleteUser = async (user) => {
  try {
    loading.value = true;

    // 사용 여부 확인
    const inUse = await userStore.checkUserUsage(user.id);

    const confirmOptions = {
      confirmButtonText: '삭제',
      cancelButtonText: '취소',
      type: 'warning'
    };

    let confirmed = false;

    if (inUse) {
      const result = await ElMessageBox.confirm(
        `'${user.name}' 사용자는 현재 사용 중입니다.\n관련된 모든 데이터가 함께 삭제됩니다.\n정말 삭제하시겠습니까?`,
        '사용자 강제 삭제 확인',
        { ...confirmOptions, type: 'error' }
      );
      confirmed = result === 'confirm';

      if (confirmed) {
        await userStore.forceDeleteUser(user.id);
      }
    } else {
      const result = await ElMessageBox.confirm(
        `'${user.name}' 사용자를 삭제하시겠습니까?`,
        '사용자 삭제 확인',
        confirmOptions
      );
      confirmed = result === 'confirm';

      if (confirmed) {
        await userStore.deleteUser(user.id);
      }
    }
  } catch (error) {
    if (error === 'cancel') {
      // 사용자가 취소한 경우
      return;
    }
    console.error('사용자 삭제 오류:', error);
  } finally {
    loading.value = false;
  }
};

/**
 * 추가 폼 초기화
 */
const resetAddForm = () => {
  newUser.value = { name: '', email: '' };
  if (addFormRef.value) {
    addFormRef.value.resetFields();
  }
};

/**
 * 컴포넌트 마운트 시 사용자 목록 로드
 */
onMounted(async () => {
  try {
    await userStore.fetchUsers();
  } catch (error) {
    console.error('사용자 목록 로드 오류:', error);
    ElMessage.error('사용자 목록을 불러올 수 없습니다');
  }
});
</script>

<style scoped>
/* 사용자 관리 컴포넌트 스타일 */
</style>
