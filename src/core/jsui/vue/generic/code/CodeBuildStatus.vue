<template>
    <div>
        <v-icon
            :color="iconColor"
            small
            v-if="hasData"
        >
            {{ iconName }}
        </v-icon>

        <v-progress-circular
            indeterminate
            size="16"
            v-else
        ></v-progress-circular>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import {
    getCodeBuildStatus, TGetCodeBuildStatusOutput,
} from '../../../ts/api/apiCode'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        commit: String
    }
})

const refreshPeriodMs = 5000

@Component
export default class CodeBuildStatus extends Props {
    hasData: boolean = false
    isPending : boolean = true
    isSuccess: boolean = false

    get iconColor() {
        if (this.isPending) {
            return 'warning'
        } else if (this.isSuccess) {
            return 'success'
        } else {
            return 'error'
        }
    }

    get iconName() {
        if (this.isPending) {
            return 'mdi-help-circle'
        } else if (this.isSuccess) {
            return 'mdi-check-circle'
        } else {
            return 'mdi-close-circle'
        }
    }

    @Watch('commit')
    refreshStatus() {
        getCodeBuildStatus({
            orgId: PageParamsStore.state.organization!.Id,
            commitHash: this.commit,
        }).then((resp : TGetCodeBuildStatusOutput) => {
            this.hasData = true
            this.isPending = resp.data.Pending
            this.isSuccess = resp.data.Success

            if (this.isPending) {
                setTimeout(this.refreshStatus, refreshPeriodMs)
            }
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
        this.refreshStatus()
    }
}

</script>
