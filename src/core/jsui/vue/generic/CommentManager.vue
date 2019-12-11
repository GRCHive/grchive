<template>
    <div>
        <v-form v-model="formValid" class="ma-4">
            <v-textarea v-model="commentText"
                        label="Comment"
                        filled
            ></v-textarea> 

            <v-row justify="end">
                <v-btn color="primary" @click="submitComment">
                    Submit
                </v-btn>
            </v-row>
        </v-form>
        <v-divider></v-divider>

        <v-list two-line v-if="!loading">
            <template v-for="(item, index) in comments">
                <v-list-item
                    :key="item.Id"
                >
                    <v-list-item-content>
                        <v-list-item-title>
                            {{ userIdToName(item.UserId) }}
                            <span class="caption">
                                 {{ item.PostTime.toString() }}
                            </span>
                        </v-list-item-title>

                        <v-list-item-subtitle>
                            <pre>{{ item.Content }}</pre>
                        </v-list-item-subtitle>

                    </v-list-item-content>
                </v-list-item>
                <v-divider
                    :key="index"
                    v-if="index != comments.length - 1"
                ></v-divider>
            </template>
        </v-list>

        <v-row align="center" justify="center" class="py-4" v-else>
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Comment, CommentResource } from '../../ts/comments'
import * as apiComments from '../../ts/api/apiComments'
import { contactUsUrl } from '../../ts/url'
import MetadataStore from '../../ts/metadata'
import { PageParamsStore } from '../../ts/pageParams'
import { createUserString } from '../../ts/users'

const Props = Vue.extend({
    props: {
        params: Object,
        type: Number
    }
})

@Component
export default class CommentManager extends Props {
    loading: boolean = true
    comments : Comment[] = []

    formValid: boolean = false
    commentText: string = ""

    get userIdToName() : (a0 : number) => string {
        if (!MetadataStore.state.usersInitialized) {
            return (a0: number) => "Loading..."
        }

        return (userId: number) => createUserString(MetadataStore.getters.getUser(userId))
    }

    onRetrieveSuccess(resp : apiComments.TGetAllCommentsOutput) {
        this.comments = resp.data
        this.loading = false
    }

    onCommentSuccess(resp : apiComments.TNewCommentOutput) {
        this.commentText = ""
        this.comments.unshift(resp.data)
    }

    onError() {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops! Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    refreshData() {
        this.loading = true

        switch (this.type) {
        case CommentResource.DocumentRequest:
            apiComments.allDocumentRequestComments(this.params as apiComments.TDocRequestNewCommentInput).
                then(this.onRetrieveSuccess).catch(this.onError)
            break
        default:
            break
        }
    }

    mounted() {
        this.refreshData()
    }

    submitComment() {
        switch (this.type) {
        case CommentResource.DocumentRequest:
            apiComments.newDocumentRequestComment({
                comment: {
                    userId: PageParamsStore.state.user!.Id,
                    content: this.commentText,
                },
                ...this.params
            } as apiComments.TDocRequestNewCommentInput).then(this.onCommentSuccess).catch(this.onError)
            break
        default:
            break
        }
    }
}

</script>
