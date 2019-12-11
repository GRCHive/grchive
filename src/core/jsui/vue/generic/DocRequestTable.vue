<script lang="ts">

import Vue, {VNode} from 'vue'
import { VIcon } from 'vuetify/lib'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { createSingleDocRequestUrl, contactUsUrl } from '../../ts/url'
import { TGetDocCatOutput, getDocumentCategory } from '../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../ts/pageParams'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class DocRequestTable extends ResourceTableProps {
    catIdToName : Record<number, string> = {}
    pendingCatRequests : Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Document Category',
                value: 'docCat',
            },
            {
                text: 'Requester',
                value: 'requester',
            },
            {
                text: 'Request Time',
                value: 'requestTime',
            },
            {
                text: 'Complete',
                value: 'complete',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    loadDocumentCategoryName(id : number) {
        if (id in this.catIdToName) {
            return
        }

        if (this.pendingCatRequests.has(id)) {
            return
        }

        this.pendingCatRequests.add(id)
        new Promise<string>((resolve, reject) => {
            getDocumentCategory({
                orgId: PageParamsStore.state.organization!.Id,
                catId: id,
                lean: true,
            }).then((resp : TGetDocCatOutput) => {
                resolve(resp.data.Cat.Name)
            }).catch((err : any) => {
                reject()
            })
        }).then((val : string) => {
            Vue.set(this.catIdToName, id, val)
            this.pendingCatRequests.delete(id)
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

    goToDocRequest(item : any) {
        window.location.assign(createSingleDocRequestUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    transformInputResourceToTableItem(inp : any) : any {
        this.loadDocumentCategoryName(inp.CatId)

        return {
            id: inp.Id,
            name: inp.Name,
            docCat: !!this.catIdToName[inp.CatId] ? this.catIdToName[inp.CatId] : "Loading...", 
            requester: createUserString(MetadataStore.getters.getUser(inp.RequestedUserId)),
            requestTime: inp.RequestTime.toString(),
            complete: !!inp.CompletionTime,
            value: inp
        }
    }

    renderFulfilled(props : any) : VNode {
        return this.$createElement(
            VIcon,
            {
                props: {
                    small : true,
                    color: props.item.complete ? 'primary' : 'error'
                }
            },
            props.item.complete ? 'mdi-check' : 'mdi-close'
        )
    }

    renderExpansion(props : any) : VNode {
        return this.$createElement(
            'td',
            {
                attrs: {
                    colspan: props.headers.length
                },
            },
            props.item.value.Description
        )
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToDocRequest
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.complete': this.renderFulfilled,
                }
            }
        )
    }
}

</script>

