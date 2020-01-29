<script lang="ts">

import Vue, {VNode} from 'vue'
import { VChip, VIcon, VMenu, VList, VListItem } from 'vuetify/lib'
import Component, { mixins } from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { standardFormatDate } from '../../ts/time'
import { createSingleDocFileUrl, contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import { 
    TAllFileVersionsOutput,
    allFileVersions,
    TGetVersionStorageDataOutput,
    getVersionStorageData
} from '../../ts/api/apiControlDocumentation' 
import { VersionedMetadata, FileVersion } from '../../ts/controls'

const DocProps = Vue.extend({
    props: {
        disableVersionSelect: {
            type: Boolean,
            default: false
        }
    }
})

@Component({
    components: {
        BaseResourceTable
    }
})
export default class DocFileTable extends mixins(ResourceTableProps, DocProps) {
    fileVersions : Record<number, number[]> = Object()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name'
            },
            {
                text: 'Relevant Time',
                value: 'relevantTime'
            },
            {
                text: 'Version',
                value: 'version',
            },
            {
                text: 'Filename',
                value: 'filename'
            },
            {
                text: 'Upload Time',
                value: 'uploadTime'
            },
            {
                text: 'User',
                value: 'user'
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    getFileVersions(id : number, obj : any) : number[] {
        if (!(id in this.fileVersions)) {
            allFileVersions({
                fileId: id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllFileVersionsOutput) => {
                obj.availableVersions = resp.data
                if (obj.availableVersions.length > 0) {
                    this.selectVersion(obj, obj.availableVersions[0])
                }
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
            return []
        }

        return this.fileVersions[id]
    }

    goToDocFile(item : any) {
        window.location.assign(createSingleDocFileUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id,
            !!item.version ? item.version.VersionNumber : null
        ))
    }

    selectVersion(obj : any, v : FileVersion) {
        if (obj.availableVersions.findIndex((ele : FileVersion) => ele.VersionNumber == v.VersionNumber) == -1) {
            return
        }
        
        obj.version = v
        getVersionStorageData({
            fileId: obj.id,
            orgId: PageParamsStore.state.organization!.Id,
            version: v.VersionNumber
        }).then((resp : TGetVersionStorageDataOutput) => {
            obj.uploadTime = standardFormatDate(resp.data.Storage.UploadTime)
            obj.user = createUserString(MetadataStore.getters.getUser(resp.data.Storage.UploadUserId))
            obj.filename = resp.data.Storage.StorageName
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

    transformInputResourceToTableItem(inp : any) : any {
        let obj = {
            id: inp.Id,
            name: inp.AltName,
            filename: "Loading...",
            relevantTime: standardFormatDate(inp.RelevantTime),
            availableVersions: null,
            version: null,
            uploadTime: "Loading...",
            user: "Loading...",
            value: inp
        }

        return obj
    }

    renderVersion(props : any) : VNode {
        if (!props.item.availableVersions) {
            props.item.availableVersions = this.getFileVersions(props.item.id, props.item)
        }

        if (!props.item.version && props.item.availableVersions.length > 0) {
            this.selectVersion(props.item, props.item.availableVersions[0])
        }

        let dropdownIcon = this.$createElement(
            VIcon,
            {
                props: {
                    small: true
                },
            },
            "mdi-chevron-down"
        )

        let availableVersionsNodes = !!props.item.availableVersions ? props.item.availableVersions.map((ele : FileVersion) => {
            let item = this.$createElement(
                VListItem,
                {
                    on: {
                        click: () => {
                            this.selectVersion(props.item, ele)
                        }
                    }
                },
                ele.VersionNumber.toString()
            )
            return item
        }) : []

        let menuList = this.$createElement(
            VList,
            {
                props: {
                    dense: true
                },
            },
            availableVersionsNodes
        )

        let menu = this.$createElement(
            VMenu,
            {
                props: {
                    "offset-y": true,
                    "disabled": this.disableVersionSelect
                },
                scopedSlots: {
                    activator: (chipProps: any) => {
                        let chip = this.$createElement(
                            VChip,
                            {
                                props: {
                                    pill: true
                                },
                                on : {
                                    ...chipProps.on,
                                    click: (e : MouseEvent) => {
                                        e.stopPropagation()
                                        chipProps.on.click(e)
                                    }
                                }
                            },
                            [
                                !!props.item.version ? props.item.version.VersionNumber : "Loading...",
                                dropdownIcon
                            ]
                        )
                        return chip
                    }
                }
            },
            [menuList]
        )

        return menu
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
                    resourceName: "files",
                    showExpand: true
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ({
                        File: ele.value,
                        Version: ele.version,
                    } as VersionedMetadata))),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToDocFile
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.version': this.renderVersion,
                }
            }
        )
    }
}

</script>


