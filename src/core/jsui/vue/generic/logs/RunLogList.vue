<template>
    <script-run-table
        :resources="allRuns"
    >
    </script-run-table>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ScriptRun } from '../../../ts/code'
import {
    allCodeRuns, TAllCodeRunsOutput,TAllCodeRunsInput
} from '../../../ts/api/apiCode'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import ScriptRunTable from '../ScriptRunTable.vue'

const Props = Vue.extend({
    props: {
        scriptId: {
            type: Number,
            default: -1,
        }
    }
})

@Component({
    components: {
        ScriptRunTable,
    }
})
export default class RunLogList extends Props {
    allRuns : ScriptRun[] = []

    refreshLogs() {
        let params : TAllCodeRunsInput = {
            orgId: PageParamsStore.state.organization!.Id,
        }

        if (this.scriptId != -1) {
            params.scriptId = this.scriptId
        }

        allCodeRuns(params).then((resp : TAllCodeRunsOutput) => {
            this.allRuns = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshLogs()
    }
}

</script>
