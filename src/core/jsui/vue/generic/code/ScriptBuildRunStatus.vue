<template>
    <div class="statusContainer">
        <v-icon
            :color="iconColor"
            small
        >
            {{ iconName }}
        </v-icon>

        <div v-if="showTimeStamp && !notRequired && !forceFail" class="ml-2">
            <p class="ma-0 caption" v-if="!!start">
                <span class="font-weight-bold">Start: </span>
                {{ timeStartStr }}
            </p>

            <p class="ma-0 caption" v-if="!!end">
                <span class="font-weight-bold">End: </span>
                {{ timeEndStr }}
            </p>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { standardFormatTime } from '../../../ts/time'

const Props = Vue.extend({
    props: {
        success: Boolean,
        start: Date,
        end: Date,
        showTimeStamp: {
            type: Boolean,
            default: false,
        },
        notRequired: {
            type: Boolean,
            default: false,
        },
        forceFail: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class ScriptBuildRunStatus extends Props {
    get isPending() {
        return !this.end
    }

    get isSuccess() {
        return this.success && !this.forceFail
    }

    get timeStartStr() : string {
        return standardFormatTime(new Date(this.start))
    }

    get timeEndStr() : string {
        return standardFormatTime(new Date(this.end))
    }

    get iconColor() {
        if (this.notRequired) {
            return 'secondary'
        } else if (this.isPending) {
            return 'warning'
        } else if (this.isSuccess) {
            return 'success'
        } else {
            return 'error'
        }
    }

    get iconName() {
        if (this.notRequired) {
            return 'mdi-minus-circle'
        } else if (this.isPending) {
            return 'mdi-help-circle'
        } else if (this.isSuccess) {
            return 'mdi-check-circle'
        } else {
            return 'mdi-close-circle'
        }
    }
}

</script>

<style scoped>

.statusContainer {
    display: flex;
    align-items: center;
    justify-items: center;
}

</style>
