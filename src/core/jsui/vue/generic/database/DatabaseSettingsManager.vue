<template>
    <div>
        <div v-if="!!settings">
            <v-checkbox
                v-model="settings.AutoRefreshEnabled"
                hide-details
                label="Auto-Refresh?"
                :readonly="!canEdit"
            >
            </v-checkbox>

            <create-scheduled-event-form
                v-model="refreshSchedule"
                no-name
                v-if="settings.AutoRefreshEnabled"
                :readonly="!canEdit"
                force-repeat
            >
            </create-scheduled-event-form>

            <v-list-item class="px-0">
                <v-list-item-action>
                    <v-btn
                        color="error"
                        v-if="canEdit"
                        @click="cancelEdit"
                    >
                        Cancel
                    </v-btn>
                </v-list-item-action>

                <v-spacer></v-spacer>

                <v-list-item-action>
                    <v-btn
                        color="success"
                        v-if="canEdit"
                        @click="saveEdit"
                        :loading="saveInProgress"
                    >
                        Save
                    </v-btn>

                    <v-btn
                        color="success"
                        v-if="!canEdit"
                        @click="canEdit = true"
                    >
                        Edit
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </div>

        <v-row align="center" justify="center" v-else>
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import {
    DatabaseSettings
} from '../../../ts/databases'
import {
    ScheduledEvent,
    createEmptyScheduledEvent,
    createScheduledEventFromRRule,
} from '../../../ts/event'
import {
    getDatabaseSettings, TGetDatabaseSettingsOutputs,
    saveDatabaseSettings
} from '../../../ts/api/apiDatabases'
import {
    contactUsUrl,
} from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import CreateScheduledEventForm from '../CreateScheduledEventForm.vue'

@Component({
    components: {
        CreateScheduledEventForm,
    }
})
export default class DatabaseSettingsManager extends Vue {
    @Prop({ type: Number, required: true})
    readonly dbId! : number

    refSettings : DatabaseSettings | null = null

    settings : DatabaseSettings | null = null
    refreshSchedule : ScheduledEvent | null = null

    canEdit : boolean = false
    saveInProgress : boolean = false

    @Watch('settings.AutoRefreshEnabled')
    resetSchedule() {
        if (this.settings!.AutoRefreshEnabled) {
            if (!!this.settings!.AutoRefreshRRule) {
                this.refreshSchedule = createScheduledEventFromRRule(this.settings!.AutoRefreshRRule)
            } else {
                this.refreshSchedule = createEmptyScheduledEvent(true)
            }
        } else {
            this.refreshSchedule = null
            this.settings!.AutoRefreshRRule = null
        }
    }

    refreshData() {
        getDatabaseSettings({
            orgId: PageParamsStore.state.organization!.Id,
            dbId: this.dbId,
        }).then((resp : TGetDatabaseSettingsOutputs) => {
            this.refSettings = resp.data
            this.cancelEdit()
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
        this.refreshData()
    }

    cancelEdit() {
        this.canEdit = false
        this.settings = JSON.parse(JSON.stringify(this.refSettings!))
        this.resetSchedule()
    }

    saveEdit() {
        this.saveInProgress = true
        saveDatabaseSettings({
            orgId: PageParamsStore.state.organization!.Id,
            dbId: this.dbId,
            autoRefreshEnabled: this.settings!.AutoRefreshEnabled,
            autoRefreshSchedule: this.refreshSchedule,
        }).then(() => {
            this.canEdit = false
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.saveInProgress = false
        })
    }
}

</script>
