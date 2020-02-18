<template>
    <div>
        <v-form v-model="formValid" class="ma-4">
            <v-textarea v-model="commentText"
                        label="Comment"
                        filled
                        hide-details
                        class="mb-4"
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
                <single-comment-viewer
                    :params="params"
                    :key="`content-${item.Id}`"
                    :comment="item"
                    @on-delete="deleteComment(item.Id)"
                    @on-edit="editComment(item.Id, arguments[0])"
                >
                </single-comment-viewer>
                <v-divider
                    :key="`divider-${item.Id}`"
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
import { Comment } from '../../ts/comments'
import * as apiComments from '../../ts/api/apiComments'
import { contactUsUrl } from '../../ts/url'
import MetadataStore from '../../ts/metadata'
import { PageParamsStore } from '../../ts/pageParams'
import { createUserString } from '../../ts/users'
import { standardFormatTime } from '../../ts/time'
import SingleCommentViewer from './SingleCommentViewer.vue'

const Props = Vue.extend({
    props: {
        params: Object,
    }
})

@Component({
    components: {
        SingleCommentViewer
    }
})
export default class CommentManager extends Props {
    loading: boolean = true
    comments : Comment[] = []

    formValid: boolean = false
    commentText: string = ""

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

        apiComments.allComments(this.params).
            then(this.onRetrieveSuccess).catch(this.onError)
    }

    mounted() {
        this.refreshData()
    }

    submitComment() {
        apiComments.newComment({
            comment: {
                userId: PageParamsStore.state.user!.Id,
                content: this.commentText,
            },
            ...this.params
        }).then(this.onCommentSuccess).catch(this.onError)
    }

    deleteComment(commentId : number) {
        apiComments.deleteComment({
            commentId: commentId,
        }).then(() => {
            let idx : number = this.comments.findIndex((ele : Comment) => ele.Id == commentId)
            if (idx == -1) {
                return
            }
            this.comments.splice(idx, 1)
        }).catch(this.onError)
    }

    editComment(commentId: number, text : string) {
        apiComments.updateComment({
            commentId: commentId,
            content: text
        }).then((resp : apiComments.TUpdateCommentOutput) => {
            let idx : number = this.comments.findIndex((ele : Comment) => ele.Id == resp.data.Id)
            if (idx == -1) {
                return
            }
            Vue.set(this.comments, idx, resp.data)
        }).catch(this.onError)
    }
}

</script>
