<template>
    <v-menu
        offset-y
        v-model="showPicker"
        :close-on-content-click="false"
        content-class="dt-menu"
        :disabled="readonly"
        min-width="0px"
    >
        <template v-slot:activator="{ on }">
            <v-text-field
                :value="timeStr"
                :label="label"
                prepend-icon="mdi-calendar"
                readonly
                v-on="on"
            >
                <template v-slot:append v-if="clearable && !readonly">
                    <v-btn icon @click="onInput(null)">
                        <v-icon>
                            mdi-close
                        </v-icon>
                    </v-btn>
                </template>
            </v-text-field>
        </template>

        <date-time-picker
            :value="value"
            @input="onInput"
            class="dt-bg"
        >
        </date-time-picker>
    </v-menu>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import DateTimePicker from './DateTimePicker.vue'
import { standardFormatTime } from '../../ts/time'

const Props = Vue.extend({
    props: {
        value: {
            type: Date,
            default: () => new Date(),
        },
        label: {
            type: String,
            default: "Time"
        },
        readonly: {
            type: Boolean,
            default: false
        },
        clearable: {
            type: Boolean,
            default: false
        },
    }
})

@Component({
    components: {
        DateTimePicker
    }
})
export default class DateTimePickerFormComponent extends Props {
    showPicker: boolean = false

    get timeStr() : string {
        if (!this.value) {
            return ""
        }
        return standardFormatTime(new Date(this.value))
    }

    onInput(d : Date) {
        this.$emit('input', d)
    }

    @Watch('showPicker')
    resetDate() {
        if (this.showPicker && !this.value) {
            this.onInput(new Date())
        }
    }
}

</script>

<style scoped>
.dt-menu {
    box-shadow: none !important;
}

.dt-bg {
    background-color: transparent;
}
</style>
