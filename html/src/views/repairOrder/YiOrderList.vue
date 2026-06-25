<template>
  <div class="yi-order-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('menu.yiOrder') }}</span>
        </div>
      </template>

      <SearchForm
        :model="searchForm"
        :fields="searchFields"
        :initial-values="initialSearchForm"
        i18n-prefix="yiOrder"
        @search="handleSearch"
        @reset="handleReset"
      >
        <template #extra-buttons>
          <el-button 
            type="primary" 
            @click="handleImport"
          >
            <el-icon><Upload /></el-icon>
            {{ $t('common.import') }}
          </el-button>
          <el-button 
            type="success" 
            :disabled="isExporting"
            :loading="isExporting"
            @click="handleExport"
          >
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
        </template>
      </SearchForm>

      <VxeTable
        ref="tableRef"
        :data="tableData"
        :loading="loading"
        :columns="tableColumns"
        :height="600"
        @sort-change="handleSortChange"
      />
      <Pagination
        v-model="pagination"
        :auto-load="true"
        :on-page-change="loadData"
      />

      <input
        ref="fileInputRef"
        type="file"
        accept=".csv,.xlsx"
        style="display: none"
        @change="handleFileChange"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Upload, Download } from '@element-plus/icons-vue'
import SearchForm from '../../components/SearchForm.vue'
import Pagination from '../../components/Pagination.vue'
import VxeTable from '../../components/VxeTable.vue'
import { useListPage } from '../../composables/useListPage'
import { getYiOrderList, exportYiOrder, importYiOrder } from '../../api/yiOrder'
import { repairEnums } from '../../enums/repair'

const { t } = useI18n()
const tableRef = ref(null)

const isImporting = ref(false)
const isExporting = ref(false)
const fileInputRef = ref(null)

const initialSearchForm = {
  yi_no: '',
  jiu_order_number: '',
  contract_no: '',
  dmo: '',
  estate: '',
  worker: '',
  project_type: '',
  is_complete: ''
}

const {
  pagination,
  tableData,
  loading,
  searchForm,
  loadData,
  handleSearch,
  handleReset,
  handleSortChange
} = useListPage({
  fetchApi: getYiOrderList,
  initialSearchForm,
  fieldMapping: {},
  defaultSort: 'id:desc',
  tableRef: computed(() => tableRef.value?.tableRef)
})

loadData()

const searchFields = computed(() => [
  {
    prop: 'yi_no',
    label: t('yiOrder.yi_no'),
    type: 'input',
    width: '180px'
  },
  {
    prop: 'contract_no',
    label: t('yiOrder.contract_no'),
    type: 'select',
    options: repairEnums.contractNoOptions,
    width: '150px'
  },
  {
    prop: 'dmo',
    label: t('yiOrder.dmo'),
    type: 'select',
    options: repairEnums.dmoOptionsYi,
    width: '150px'
  },
  {
    prop: 'estate',
    label: t('yiOrder.estate'),
    type: 'select',
    options: repairEnums.estateOptions,
    width: '180px'
  },
  {
    prop: 'worker',
    label: t('yiOrder.worker'),
    type: 'input',
    width: '120px'
  },
  {
    prop: 'is_complete',
    label: t('yiOrder.is_complete'),
    type: 'select',
    options: repairEnums.isCompleteOptions,
    width: '120px'
  },
  {
    prop: 'project_type',
    label: t('yiOrder.project_type'),
    type: 'select',
    options: repairEnums.projectTypeOptionsYi,
    width: '120px',
    advanced: true
  }
])

const tableColumns = computed(() => [
  {
    field: 'id',
    title: 'ID',
    width: 60,
    sortable: true
  },
  {
    field: 'yi_no',
    title: t('yiOrder.yi_no'),
    width: 140
  },
  {
    field: 'jiu_order_number',
    title: t('yiOrder.jiu_order_number'),
    width: 140
  },
  {
    field: 'contract_no',
    title: t('yiOrder.contract_no'),
    width: 140
  },
  {
    field: 'dmo',
    title: t('yiOrder.dmo'),
    width: 100
  },
  {
    field: 'estate',
    title: t('yiOrder.estate'),
    width: 140
  },
  {
    field: 'worker',
    title: t('yiOrder.worker'),
    width: 100
  },
  {
    field: 'project_type',
    title: t('yiOrder.project_type'),
    width: 120
  },
  {
    field: 'is_complete',
    title: t('yiOrder.is_complete'),
    width: 100,
    formatter: ({ row }) => row.is_complete ? t('common.yes') : t('common.no')
  },
  {
    field: 'created_at',
    title: t('common.created_at'),
    width: 160,
    sortable: true
  }
])

const handleImport = () => {
  fileInputRef.value?.click()
}

const handleExport = async () => {
  isExporting.value = true
  try {
    const params = { ...searchForm.value }
    const response = await exportYiOrder(params)
    ElMessage.success(response.message || t('common.export_success'))
  } catch (error) {
    ElMessage.error(t('common.export_fail'))
  } finally {
    isExporting.value = false
  }
}

const handleFileChange = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  isImporting.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const response = await importYiOrder(formData)
    ElMessage.success(response.message || t('common.import_success'))
    loadData()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || t('common.import_fail'))
  } finally {
    isImporting.value = false
    event.target.value = ''
  }
}
</script>
