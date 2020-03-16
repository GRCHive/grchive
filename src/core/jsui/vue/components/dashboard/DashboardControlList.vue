<template>
    <div class="ma-4">
        <control-table-with-controls></control-table-with-controls>
    </div >
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ControlTableWithControls from '../../generic/resources/ControlTableWithControls.vue'

@Component({
    components: {
        ControlTableWithControls
    }
})
export default class DashboardControlList extends Vue {
}

</script>
