<template>
    <div>
        <v-icon
            small
            :color="iconColor"
        >
            {{ iconStr }}
        </v-icon>
        <span>{{ successfulRuns }} / {{ totalRuns }}</span>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { ShellScriptRunPerServer } from '../../../ts/shell'

@Component
export default class ShellRunStatus extends Vue{
    @Prop(Array) readonly serverRuns!: ShellScriptRunPerServer[]

    get iconStr() : string {
        if (this.inProgress) {
            return 'mdi-help-circle'
        } else if (this.success) {
            return 'mdi-check-circle'
        } else {
            return 'mdi-close-circle'
        }
    }

    get iconColor() : string {
        if (this.inProgress) {
            return "warning"
        } else if (this.success) {
            return "success"
        } else {
            return "error"
        }
    }

    get totalRuns() : number {
        return this.serverRuns.length
    }

    get successfulRuns() : number {
        return this.serverRuns.filter((ele : ShellScriptRunPerServer) => ele.Success).length
    }

    get inProgressRuns() : number {
        return this.serverRuns.filter((ele : ShellScriptRunPerServer) => !ele.EndTime).length
    }

    get success() : boolean {
        return this.successfulRuns == this.totalRuns
    }

    get inProgress() : boolean {
        return this.inProgressRuns > 0
    }
}

</script>
