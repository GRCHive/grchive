<template>
    <div class="ma-4">
        <risk-table-with-controls></risk-table-with-controls>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import RiskTableWithControls from '../../generic/resources/RiskTableWithControls.vue'

@Component({
    components: {
        RiskTableWithControls
    }
})
export default class DashboardRiskList extends Vue {
}
</script>
