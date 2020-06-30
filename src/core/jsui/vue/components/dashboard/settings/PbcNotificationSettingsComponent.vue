<template>
    <div>
        <v-list-item>
            <v-list-item-content>
                <v-list-item-title class="title">
                    PBC Notification Settings
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>
        <v-divider></v-divider>

        <v-list-item>
            <v-list-item-content>
                <v-list-item-title class="subtitle-1">
                    Notification Cadence
                </v-list-item-title>
            </v-list-item-content>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-dialog persistent v-model="showHideNew" max-width="40%">
                    <template v-slot:activator="{on}">
                        <v-btn
                            color="primary"
                            v-on="on"
                        >
                            Add
                        </v-btn>
                    </template>

                    <v-card>
                        <v-card-title>
                            New PBC Notification
                        </v-card-title>
                        <v-divider></v-divider>

                        <v-form v-model="newIsValid" class="ma-4">
                            <v-text-field
                                type="number"
                                label="Days Before Due"
                                :rules="[ rules.numeric, rules.geq(0) ]"
                                v-model="newDaysBefore"
                            >
                            </v-text-field>
                        </v-form>
                        
                        <v-card-actions>
                            <v-btn
                                color="error"
                                @click="cancelNew"
                                :loading="opInProgress"
                            >
                                Cancel
                            </v-btn>

                            <v-spacer></v-spacer>

                            <v-btn
                                color="success"
                                @click="saveNew"
                                :loading="opInProgress"
                                :disabled="!newIsValid"
                            >
                                Save
                            </v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </v-list-item-action>

        </v-list-item>

        <template v-if="!loadingCadence">
            <v-stepper v-model="step" alt-labels non-linear class="no-box-shadow">
                <v-stepper-header class="no-box-shadow">
                    <template v-for="(n, idx) in sortedNotifications">
                        <v-stepper-step
                            :key="`${n.Id}-step`"
                            :step="n.DaysBeforeDue"
                            :editable="!opInProgress"
                        >
                            <span v-if="n.DaysBeforeDue == 0">
                                Due Date
                            </span>

                            <span v-else>
                                T-{{ n.DaysBeforeDue }} Days
                            </span>
                        </v-stepper-step>

                        <v-divider
                            v-if="idx !== (sortedNotifications.length - 1)"
                            :key="`${n.Id}-divider`"
                        >
                        </v-divider>
                    </template>
                </v-stepper-header>
            </v-stepper>

            <div class="ma-4" v-if="!!editableNotification">
                <v-form v-model="editableIsValid">
                    <v-text-field
                        type="number"
                        label="Days Before Due"
                        :rules="[ rules.numeric, rules.geq(0) ]"
                        :value="editableNotification.DaysBeforeDue"
                        @input="editableNotification.DaysBeforeDue = Number(arguments[0])"
                    >
                    </v-text-field>

                    <div class="d-flex">
                        <v-checkbox
                            label="Send to Requester"
                            class="send-to-checkbox"
                            v-model="editableNotification.SendToRequester"
                            hide-details
                        >
                        </v-checkbox>

                        <v-checkbox
                            label="Send to Assignee"
                            class="send-to-checkbox"
                            v-model="editableNotification.SendToAssignee"
                            hide-details
                        >
                        </v-checkbox>
                    </div>
                </v-form>

                <v-list-item class="pa-0">
                    <v-list-item-content>
                        <span class="subtitle-2">Additional Users to Notify</span>
                    </v-list-item-content>

                    <v-spacer></v-spacer>

                    <v-list-item-action>
                        <v-dialog v-model="showHideAddUser" persistent max-width="40%">
                            <template v-slot:activator="{on}">
                                <v-btn icon color="primary" v-on="on">
                                    <v-icon>mdi-plus</v-icon>
                                </v-btn>
                            </template>

                            <v-card>
                                <v-card-title>
                                    Add User
                                </v-card-title>

                                <user-search-form-component
                                    :user.sync="newUser"
                                    label="User"
                                >
                                </user-search-form-component>

                                <v-card-actions>
                                    <v-btn
                                        @click="cancelUser"
                                        color="error"
                                    >
                                        Cancel
                                    </v-btn>

                                    <v-spacer></v-spacer>

                                    <v-btn
                                        @click="addUser"
                                        :disabled="!newUser"
                                        color="success"
                                    >
                                        Add
                                    </v-btn>
                                </v-card-actions>
                            </v-card>
                        </v-dialog>
                    </v-list-item-action>
                </v-list-item>

                <user-table
                    :resources="extraUsersToNotify"
                >
                </user-table>
            </div>

            <v-list-item>
                <v-list-item-action>
                    <v-dialog persistent max-width="40%" v-model="showHideDelete">
                        <template v-slot:activator="{on}">
                            <v-btn
                                color="error"
                                :loading="opInProgress"
                                v-on="on"
                            >
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            v-if="!!editableNotification"
                            item-name="PBC notifications"
                            :items-to-delete="[`Notification ${editableNotification.DaysBeforeDue} days before due.`]"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete"
                            :use-global-deletion="false"
                            :delete-in-progress="opInProgress"
                        >
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

                <v-spacer></v-spacer>

                <v-list-item-action>
                    <v-btn
                        color="success"
                        @click="applyToAll"
                        :loading="opInProgress"
                        :disabled="!editableIsValid"
                    >
                        Apply to All
                    </v-btn>
                </v-list-item-action>

                <v-list-item-action>
                    <v-btn
                        color="primary"
                        @click="apply"
                        :loading="opInProgress"
                        :disabled="!editableIsValid"
                    >
                        Apply
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </template>

        <v-row justify="center" v-else>
            <v-progress-circular size="64" indeterminate></v-progress-circular>
        </v-row>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { PbcNotificationCadenceSettings } from '../../../../ts/settings/pbcNotifications'
import {
    getOrgPbcNotificationSettings, TGetPbcNotificationSettingsOutput,
    deletePbcNotificationSetting,
    newPbcNotificationSetting, TNewPbcNotificationSettingOutput,
    editPbcNotificationSetting,
} from '../../../../ts/api/settings/apiPbcNotificationSettings'
import { PageParamsStore } from '../../../../ts/pageParams'
import { contactUsUrl } from '../../../../ts/url'
import * as rules from '../../../../ts/formRules'
import UserTable from '../../../generic/UserTable.vue'
import MetadataStore from '../../../../ts/metadata'
import GenericDeleteConfirmationForm from '../GenericDeleteConfirmationForm.vue'
import UserSearchFormComponent from '../../../generic/UserSearchFormComponent.vue'

@Component({
    components: {
        UserTable,
        GenericDeleteConfirmationForm,
        UserSearchFormComponent,
    }
})
export default class PbcNotificationSettingsComponent extends Vue {
    notifications : PbcNotificationCadenceSettings[] | null = null
    step : number = 1
    rules : any = rules

    showHideNew : boolean = false
    showHideDelete: boolean = false
    showHideAddUser: boolean = false

    opInProgress: boolean = false
    newIsValid : boolean = false
    editableIsValid : boolean = false

    newDaysBefore : number = 5
    newUser : User | null = null

    editableNotification : PbcNotificationCadenceSettings | null = null
    
    get sortedNotifications() : PbcNotificationCadenceSettings[] {
        if (!this.notifications) {
            return []
        }

        return this.notifications.sort((a : PbcNotificationCadenceSettings, b : PbcNotificationCadenceSettings) => {
            return b.DaysBeforeDue - a.DaysBeforeDue
        })
    }

    get loadingCadence() : boolean { 
        return !this.notifications
    }

    @Watch('step')
    refreshEditable() {
        this.editableNotification = null

        if (!this.notifications) {
            return
        }

        let idx : number = this.notifications.findIndex((ele : PbcNotificationCadenceSettings) => ele.DaysBeforeDue == this.step)
        if (idx == -1) {
            return
        }

        this.editableNotification = JSON.parse(JSON.stringify(this.notifications[idx]))
    }

    get extraUsersToNotify() : User[] {
        if (!this.editableNotification) {
            return []
        }

        if (!this.editableNotification.AdditionalUsers) {
            return []
        }

        return this.editableNotification.AdditionalUsers.map((id : number) => MetadataStore.getters.getUser(id))
            .filter((ele : User | null) => !!ele)
    }

    refreshCadence() {
        getOrgPbcNotificationSettings({
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetPbcNotificationSettingsOutput) => {
            this.notifications = resp.data
            if (this.sortedNotifications.length > 0) {
                this.step = this.sortedNotifications[0].DaysBeforeDue
            }
            this.refreshEditable()
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
        this.refreshCadence()
    }

    apply() {
        this.opInProgress = true
        editPbcNotificationSetting({
            orgId: PageParamsStore.state.organization!.Id,
            setting: this.editableNotification!,
            applyAll: false
        }).then(() => {
            let idx = this.notifications!.findIndex((ele : PbcNotificationCadenceSettings) => ele.Id == this.editableNotification!.Id)
            if (idx == -1) {
                return
            }

            Vue.set(this.notifications!, idx, JSON.parse(JSON.stringify(this.editableNotification!)))
            Vue.nextTick(() => {
                this.step = this.editableNotification!.DaysBeforeDue
                this.refreshEditable()
            })
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);

        }).finally(() => this.opInProgress = false)
    }

    applyToAll() {
        this.opInProgress = true
        editPbcNotificationSetting({
            orgId: PageParamsStore.state.organization!.Id,
            setting: this.editableNotification!,
            applyAll: true
        }).then(() => {
            for (let n of this.notifications!) {
                if (n.Id == this.editableNotification!.Id) {
                    n.DaysBeforeDue = this.editableNotification!.DaysBeforeDue
                }
                n.SendToAssignee = this.editableNotification!.SendToAssignee
                n.SendToRequester = this.editableNotification!.SendToRequester
                n.AdditionalUsers = this.editableNotification!.AdditionalUsers.slice()
            }

            Vue.nextTick(() => {
                this.step = this.editableNotification!.DaysBeforeDue
                this.refreshEditable()
            })
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => this.opInProgress = false)
    }

    onDelete() {
        if (!this.editableNotification) {
            return
        }

        this.opInProgress = true
        deletePbcNotificationSetting({
            orgId: PageParamsStore.state.organization!.Id,
            notificationId: this.editableNotification.Id,
        }).then(() => {
            let idx = this.notifications!.findIndex((ele : PbcNotificationCadenceSettings) => ele.Id == this.editableNotification!.Id)
            if (idx == -1) {
                return
            }
            this.notifications!.splice(idx, 1)
            if (this.sortedNotifications.length > 0) {
                this.step = this.sortedNotifications[0].DaysBeforeDue
            }
            this.showHideDelete = false
            this.refreshEditable()
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => this.opInProgress = false)
    }

    cancelNew() {
        this.showHideNew = false
        this.newDaysBefore = 5
    }

    saveNew() {
        this.opInProgress = true

        newPbcNotificationSetting({
            orgId: PageParamsStore.state.organization!.Id,
            daysBefore: Number(this.newDaysBefore),
        }).then((resp : TNewPbcNotificationSettingOutput) => {
            this.notifications!.unshift(resp.data)
            Vue.nextTick(() => {
                this.step = resp.data.DaysBeforeDue
            })
            this.cancelNew()
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => this.opInProgress = false)
    }

    cancelUser() {
        this.newUser = null
        this.showHideAddUser = false
    }

    addUser() {
        if (!this.newUser) {
            return
        }
        this.editableNotification!.AdditionalUsers!.push(this.newUser!.Id)
        this.cancelUser()
    }
}

</script>

<style scoped>

.send-to-checkbox {
    flex-grow: 1;
}

.no-box-shadow {
    box-shadow: none;
}

</style>
