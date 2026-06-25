<template>
  <div class="jiu-order-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('menu.jiuOrder') }}</span>
        </div>
      </template>

      <SearchForm
        :model="searchForm"
        :fields="searchFields"
        :initial-values="initialSearchForm"
        i18n-prefix="jiuOrder"
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
import { getJiuOrderList, exportJiuOrder, importJiuOrder } from '../../api/jiuOrder'
import { repairEnums } from '../../enums/repair'

const { t } = useI18n()
const tableRef = ref(null)

const isImporting = ref(false)
const isExporting = ref(false)
const fileInputRef = ref(null)

const initialSearchForm = {
  jiu_no: '',
  contract_no: '',
  dmo: '',
  estate: '',
  worker: '',
  type: '',
  order_status: '',
  complete_status: '',
  project_type: '',
  date_start: '',
  date_end: ''
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
  fetchApi: getJiuOrderList,
  initialSearchForm,
  fieldMapping: {},
  defaultSort: 'id:desc',
  tableRef: computed(() => tableRef.value?.tableRef)
})

loadData()

const searchFields = computed(() => [
  {
    prop: 'jiu_no',
    label: t('jiuOrder.jiu_no'),
    type: 'input',
    width: '180px'
  },
  {
    prop: 'contract_no',
    label: t('jiuOrder.contract_no'),
    type: 'select',
    options: repairEnums.contractNoOptions,
    width: '150px'
  },
  {
    prop: 'dmo',
    label: t('jiuOrder.dmo'),
    type: 'select',
    options: repairEnums.dmoOptionsJiu,
    width: '150px'
  },
  {
    prop: 'estate',
    label: t('jiuOrder.estate'),
    type: 'select',
    options: repairEnums.estateOptions,
    width: '180px'
  },
  {
    prop: 'worker',
    label: t('jiuOrder.worker'),
    type: 'input',
    width: '120px'
  },
  {
    prop: 'type',
    label: t('jiuOrder.type'),
    type: 'select',
    options: repairEnums.jiuTypeOptions,
    width: '120px'
  },
  {
    prop: 'fee',
    label: t('jiuOrder.fee'),
    type: 'input',
    width: '100px',
    advanced: true
  },
  {
    prop: 'order_status',
    label: t('jiuOrder.order_status'),
    type: 'select',
    options: repairEnums.orderStatusOptions,
    width: '120px'
  },
  {
    prop: 'complete_status',
    label: t('jiuOrder.complete_status'),
    type: 'select',
    options: repairEnums.completeStatusOptions,
    width: '120px',
    advanced: true
  },
  {
    prop: 'project_type',
    label: t('jiuOrder.project_type'),
    type: 'select',
    options: repairEnums.projectTypeOptionsJiu,
    width: '120px',
    advanced: true
  },
  {
    prop: 'date_start',
    label: t('jiuOrder.date_start'),
    type: 'date',
    width: '150px',
    advanced: true
  },
  {
    prop: 'date_end',
    label: t('jiuOrder.date_end'),
    type: 'date',
    width: '150px',
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
    field: 'jiu_no',
    title: t('jiuOrder.jiu_no'),
    width: 140
  },
  {
    field: 'contract_no',
    title: t('jiuOrder.contract_no'),
    width: 140
  },
  {
    field: 'dmo',
    title: t('jiuOrder.dmo'),
    width: 100
  },
  {
    field: 'estate',
    title: t('jiuOrder.estate'),
    width: 140
  },
  {
    field: 'worker',
    title: t('jiuOrder.worker'),
    width: 100
  },
  {
    field: 'type',
    title: t('jiuOrder.type'),
    width: 120
  },
  {
    field: 'fee',
    title: t('jiuOrder.fee'),
    width: 100,
    formatter: ({ row }) => row.fee?.toFixed(2) || '0.00'
  },
  {
    field: 'order_status',
    title: t('jiuOrder.order_status'),
    width: 120
  },
  {
    field: 'complete_status',
    title: t('jiuOrder.complete_status'),
    width: 120
  },
  {
    field: 'project_type',
    title: t('jiuOrder.project_type'),
    width: 120
  },
  {
    field: 'date_start',
    title: t('jiuOrder.date_start'),
    width: 120
  },
  {
    field: 'date_end',
    title: t('jiuOrder.date_end'),
    width: 120
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
    const response = await exportJiuOrder(params)
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
    const response = await importJiuOrder(formData)
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
