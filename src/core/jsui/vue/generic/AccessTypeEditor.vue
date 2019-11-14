<template>
    <v-row class="ma-0">
        <v-col cols="6">
            <span class="body-1 font-weight-bold">{{ label }}</span>
        </v-col>

        <v-col cols="2">
            <v-checkbox :input-value="canView"
                        label="View"
                        hide-details
                        class="ma-0 pa-0"
                        @change="updateValue(AccessType.View)"
            ></v-checkbox>
        </v-col>

        <v-col cols="2">
            <v-checkbox :input-value="canEdit"
                        label="Edit"
                        hide-details
                        class="ma-0 pa-0"
                        @change="updateValue(AccessType.Edit)"
            ></v-checkbox>
        </v-col>

        <v-col cols="2">
            <v-checkbox :input-value="canManage"
                        label="Manage"
                        hide-details
                        class="ma-0 pa-0"
                        @change="updateValue(AccessType.Manage)"
            ></v-checkbox>
        </v-col>
    </v-row>
</template>

<script lang="ts">

import Vue from 'vue'
import { AccessType } from '../../ts/roles' 
export default Vue.extend({
    props: {
        label: String,
        value: Number,
    },
    data: () => ({
        AccessType,
    }),
    computed: {
        canView() : boolean {
            return (this.value & AccessType.View) != AccessType.NoAccess
        },
        canEdit() : boolean {
            return (this.value & AccessType.Edit) != AccessType.NoAccess
        },
        canManage() : boolean {
            return (this.value & AccessType.Manage) != AccessType.NoAccess
        },
        currentValue() : AccessType {
            let access = AccessType.NoAccess
            if (this.canView) {
                access = access | AccessType.View
            }

            if (this.canEdit) {
                access = access | AccessType.Edit
            }

            if (this.canManage) {
                access = access | AccessType.Manage
            }

            return access
        }
    },
    methods: {
        updateValue(access : AccessType) {
            let val = this.currentValue
            val ^= access
            this.$emit('input', val)
        }
    },
})
</script>
