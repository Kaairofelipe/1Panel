<template>
    <DrawerPro
        v-model="monitorVisible"
        :header="$t('menu.monitor')"
        @close="handleClose"
        :resource="title"
        size="large"
    >
        <el-form label-position="top" @submit.prevent>
            <el-form-item :label="$t('container.refreshTime')">
                <el-select v-model="timeInterval" @change="changeTimer">
                    <el-option label="3s" :value="3" />
                    <el-option label="5s" :value="5" />
                    <el-option label="10s" :value="10" />
                    <el-option label="30s" :value="30" />
                    <el-option label="60s" :value="60" />
                </el-select>
            </el-form-item>
        </el-form>
        <el-card>
            <v-charts
                height="200px"
                id="cpuChart"
                type="line"
                :option="chartsOption['cpuChart']"
                v-if="chartsOption['cpuChart']"
            />
        </el-card>
        <el-card class="mt-2.5">
            <v-charts
                height="200px"
                id="memoryChart"
                type="line"
                :option="chartsOption['memoryChart']"
                v-if="chartsOption['memoryChart']"
            />
        </el-card>
        <el-card class="mt-2.5">
            <v-charts
                height="200px"
                id="ioChart"
                type="line"
                :option="chartsOption['ioChart']"
                v-if="chartsOption['ioChart']"
            />
        </el-card>
        <el-card class="mt-2.5">
            <v-charts
                height="200px"
                id="networkChart"
                type="line"
                :option="chartsOption['networkChart']"
                v-if="chartsOption['networkChart']"
            />
        </el-card>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, ref } from 'vue';
import { containerStats } from '@/api/modules/container';
import { dateFormatForSecond } from '@/utils/date';
import VCharts from '@/components/v-charts/index.vue';
import i18n from '@/lang';

const title = ref();
const monitorVisible = ref(false);
const timeInterval = ref();
let timer: NodeJS.Timer | null = null;
let isInit = ref<boolean>(true);
interface DialogProps {
    containerID: string;
    container: string;
}
const dialogData = ref<DialogProps>({
    containerID: '',
    container: '',
});

const acceptParams = async (params: DialogProps): Promise<void> => {
    monitorVisible.value = true;
    dialogData.value.containerID = params.containerID;
    title.value = params.container;
    cpuData.value = [];
    memData.value = [];
    cacheData.value = [];
    ioReadData.value = [];
    ioWriteData.value = [];
    netTxData.value = [];
    netRxData.value = [];
    timeData.value = [];
    timeInterval.value = 5;
    isInit.value = true;
    loadData();
    timer = setInterval(async () => {
        if (monitorVisible.value) {
            isInit.value = false;
            loadData();
        }
    }, 1000 * timeInterval.value);
};

const cpuData = ref<Array<string>>([]);
const memData = ref<Array<string>>([]);
const cacheData = ref<Array<string>>([]);
const ioReadData = ref<Array<string>>([]);
const ioWriteData = ref<Array<string>>([]);
const netTxData = ref<Array<string>>([]);
const netRxData = ref<Array<string>>([]);
const timeData = ref<Array<string>>([]);
const chartsOption = ref({ cpuChart: null, memoryChart: null, ioChart: null, networkChart: null });

const changeTimer = () => {
    clearInterval(Number(timer));
    timer = setInterval(async () => {
        if (monitorVisible.value) {
            loadData();
        }
    }, 1000 * timeInterval.value);
};

const loadData = async () => {
    const res = await containerStats(dialogData.value.containerID);
    cpuData.value.push(res.data.cpuPercent.toFixed(2));
    if (cpuData.value.length > 20) {
        cpuData.value.splice(0, 1);
    }
    memData.value.push(res.data.memory.toFixed(2));
    if (memData.value.length > 20) {
        memData.value.splice(0, 1);
    }
    cacheData.value.push(res.data.cache.toFixed(2));
    if (cacheData.value.length > 20) {
        cacheData.value.splice(0, 1);
    }
    ioReadData.value.push(res.data.ioRead.toFixed(2));
    if (ioReadData.value.length > 20) {
        ioReadData.value.splice(0, 1);
    }
    ioWriteData.value.push(res.data.ioWrite.toFixed(2));
    if (ioWriteData.value.length > 20) {
        ioWriteData.value.splice(0, 1);
    }
    netTxData.value.push(res.data.networkTX.toFixed(2));
    if (netTxData.value.length > 20) {
        netTxData.value.splice(0, 1);
    }
    netRxData.value.push(res.data.networkRX.toFixed(2));
    if (netRxData.value.length > 20) {
        netRxData.value.splice(0, 1);
    }
    timeData.value.push(dateFormatForSecond(res.data.shotTime));
    if (timeData.value.length > 20) {
        timeData.value.splice(0, 1);
    }

    chartsOption.value['cpuChart'] = {
        title: 'CPU',
        xData: timeData.value,
        yData: [
            {
                name: 'CPU',
                data: cpuData.value,
            },
        ],
        formatStr: '%',
    };

    chartsOption.value['memoryChart'] = {
        title: i18n.global.t('monitor.memory'),
        xData: timeData.value,
        yData: [
            {
                name: i18n.global.t('monitor.memory'),
                data: memData.value,
            },
            {
                name: i18n.global.t('container.cache'),
                data: cacheData.value,
            },
        ],
        formatStr: 'MB',
    };

    chartsOption.value['ioChart'] = {
        title: i18n.global.t('monitor.disk') + ' IO',
        xData: timeData.value,
        yData: [
            {
                name: i18n.global.t('monitor.read'),
                data: ioReadData.value,
            },
            {
                name: i18n.global.t('monitor.write'),
                data: ioWriteData.value,
            },
        ],
        formatStr: 'MB',
    };

    chartsOption.value['networkChart'] = {
        title: i18n.global.t('monitor.network'),
        xData: timeData.value,
        yData: [
            {
                name: i18n.global.t('monitor.up'),
                data: netTxData.value,
            },
            {
                name: i18n.global.t('monitor.down'),
                data: netRxData.value,
            },
        ],
        formatStr: 'KB',
    };
};
const handleClose = async () => {
    monitorVisible.value = false;
    clearInterval(Number(timer));
    timer = null;
    chartsOption.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
};

onBeforeUnmount(() => {
    handleClose;
});

defineExpose({
    acceptParams,
});
</script>
