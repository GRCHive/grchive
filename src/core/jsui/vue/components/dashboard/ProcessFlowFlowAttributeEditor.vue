<template>
    <div>
        <div>
            <p class="text-center font-weight-bold outline">
                Risk & Control Linkage
            </p>
            <v-divider></v-divider>
            <process-flow-risk-control-linkage></process-flow-risk-control-linkage>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ProcessFlowRiskControlLinkage from './ProcessFlowRiskControlLinkage.vue'

@Component({
    components: {
        ProcessFlowRiskControlLinkage
    }
})
export default class ProcessFlowFlowAttributeEditor extends Vue {

}

</script>
