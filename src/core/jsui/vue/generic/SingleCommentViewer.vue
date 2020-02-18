<template>
    <v-list-item>
        <v-dialog persistent max-width="40%" v-model="showDelete">
            <generic-delete-confirmation-form
                item-name="comments"
                :items-to-delete="[comment.Content]"
                @do-cancel="showDelete = false"
                @do-delete="confirmDelete"
                :use-global-deletion="false">
            </generic-delete-confirmation-form>
        </v-dialog>

        <v-list-item-content class="content">
            <v-list-item-title>
                {{ userIdToName(comment.UserId) }}
                <span class="caption">
                     {{ standardFormatTime(comment.PostTime) }}
                </span>
            </v-list-item-title>

            <v-list-item-subtitle>
                <pre v-if="!editMode">{{ comment.Content }}</pre>
                <div v-else>
                    <v-textarea v-model="editableText"
                                label="Comment"
                                filled
                                hide-details
                    ></v-textarea> 

                    <v-list-item class="px-0">
                        <v-list-item-action>
                            <v-btn
                                color="error"
                                @click="cancelEdit"
                            >
                                Cancel
                            </v-btn>
                        </v-list-item-action>

                        <v-spacer></v-spacer>

                        <v-list-item-action>
                            <v-btn
                                color="success"
                                @click="saveEdit"
                                :loading="saveInProgress"
                            >
                                Save
                            </v-btn>
                        </v-list-item-action>
                    </v-list-item>
                </div>
            </v-list-item-subtitle>

        </v-list-item-content>

        <v-spacer></v-spacer>

        <v-list-item-action>
            <v-menu offset-y :disabled="comment.UserId != currentUserId">
                <template v-slot:activator="{on}">
                    <v-btn icon v-on="on">
                        <v-icon>
                            mdi-dots-vertical
                        </v-icon>
                    </v-btn>
                </template>

                <v-list dense>
                    <v-list-item :disabled="editMode" @click="editMode = true">
                        <v-list-item-content>
                            <v-list-item-title>
                                Edit
                            </v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>

                    <v-list-item @click="showDelete = true">
                        <v-list-item-content>
                            <v-list-item-title>
                                Delete
                            </v-list-item-title>
                        </v-list-item-content>
                    </v-list-item>
                </v-list>
            </v-menu>
        </v-list-item-action>
    </v-list-item>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { Comment } from '../../ts/comments'
import * as apiComments from '../../ts/api/apiComments'
import { contactUsUrl } from '../../ts/url'
import MetadataStore from '../../ts/metadata'
import { PageParamsStore } from '../../ts/pageParams'
import { createUserString } from '../../ts/users'
import { standardFormatTime } from '../../ts/time'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'

const Props = Vue.extend({
    props: {
        params: Object,
        comment: {
            type: Object,
            default: () => Object() as Comment
        }
    }
})

@Component({
    components: {
        GenericDeleteConfirmationForm
    }
})
export default class SingleCommentViewer extends Props {
    standardFormatTime : (arg0: Date) => string = standardFormatTime
    editMode: boolean = false
    showDelete: boolean = false
    saveInProgress: boolean = false

    editableText : string = ""

    get currentUserId() : number {
        return PageParamsStore.state.user!.Id
    }

    get userIdToName() : (a0 : number) => string {
        if (!MetadataStore.state.usersInitialized) {
            return (a0: number) => "Loading..."
        }

        return (userId: number) => createUserString(MetadataStore.getters.getUser(userId))
    }

    confirmDelete() {
        this.showDelete = false
        this.$emit('on-delete')
    }

    mounted() {
        this.cancelEdit()
    }

    @Watch('comment')
    cancelEdit() {
        this.editMode = false
        this.saveInProgress = false
        this.editableText = this.comment.Content
    }

    saveEdit() {
        this.saveInProgress = true
        this.$emit('on-edit', this.editableText)
    }
}

</script>

<style scoped>

.content {
    flex: 12 1 !important;
}

</style>
