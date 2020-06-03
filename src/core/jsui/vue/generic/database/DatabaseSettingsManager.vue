<template>
    <div>
        <div v-if="!!settings">
            <span class="title">Refresh</span>
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

            <v-divider class="my-2"></v-divider>

            <div>
                <span class="title">Notifications</span>
            </div>

            <div style="display: flex;">
                <span class="subtitle-2" style="align-self: center;">On Schema Change</span>
                <v-dialog
                    persistent
                    max-width="40%"
                    v-model="showHideAddSchemaChangeUser"
                >
                    <template v-slot:activator="{on}">
                        <v-btn
                            icon
                            color="primary"
                            :disabled="!canEdit"
                            v-on="on"
                        >
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <v-card>
                        <v-card-title>
                            Select User
                        </v-card-title>
                        <v-divider></v-divider>

                        <user-search-form-component
                            class="ma-4"
                            label="User"
                            :user.sync="userToAdd"
                        >
                        </user-search-form-component>

                        <v-card-actions>
                            <v-btn
                                color="error"
                                @click="showHideAddSchemaChangeUser = false"
                            >
                                Cancel
                            </v-btn>

                            <v-spacer></v-spacer>

                            <v-btn
                                color="success"
                                @click="addSchemaChangeNotifyUser"
                            >
                                Add
                            </v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </div>

            <user-table
                :resources="settings.OnSchemaChangeNotifyUsers"
                :use-crud-delete="canEdit"
                @delete="removeSchemaChangeNotifyUser"
            >
            </user-table>

            <v-divider class="my-2"></v-divider>

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
import UserTable from '../UserTable.vue'
import UserSearchFormComponent from '../UserSearchFormComponent.vue'

@Component({
    components: {
        CreateScheduledEventForm,
        UserTable,
        UserSearchFormComponent,
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

    showHideAddSchemaChangeUser : boolean = false
    userToAdd : User | null = null

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
            onSchemaChangeNotifyUsers: this.settings!.OnSchemaChangeNotifyUsers.map((ele : User) => ele.Id),
        }).then(() => {
            this.canEdit = false
            this.refSettings = JSON.parse(JSON.stringify(this.settings!))
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

    addSchemaChangeNotifyUser() {
        this.showHideAddSchemaChangeUser = false
        if (!!this.userToAdd) {
            let idx = this.settings!.OnSchemaChangeNotifyUsers.findIndex(
                (ele : User) => ele.Id == this.userToAdd!.Id)

            if (idx == -1) {
                this.settings!.OnSchemaChangeNotifyUsers.push(this.userToAdd)
            }
        }
        this.userToAdd = null
    }

    removeSchemaChangeNotifyUser(u : User) {
        let idx = this.settings!.OnSchemaChangeNotifyUsers.findIndex(
            (ele : User) => ele.Id == u.Id)

        if (idx == -1) {
            return
        }

        this.settings!.OnSchemaChangeNotifyUsers.splice(idx, 1)
    }
}

</script>
