<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import { createUserString } from '../../ts/users'
import { VChip, VIcon, VMenu, VList, VListItem } from 'vuetify/lib'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import { ShellTypes, ShellScriptVersion } from '../../ts/shell'
import {
    allShellScriptVersions, TAllShellScriptVersionsOutput
} from '../../ts/api/apiShell'
import { contactUsUrl, createSingleShellUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ShellTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Type',
                value: 'type',
            },
            {
                text: 'Version',
                value: 'version',
            },
            {
                text: 'Upload Time',
                value: 'uploadTime',
                sort: (a : Date | null, b : Date | null) => {
                    if (a == null || b == null) {
                        return 0
                    }
                    return a.getTime() - b.getTime()
                }
            },
            {
                text: 'Upload User',
                value: 'uploadUser',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToShellScript(item : any) {
        window.location.assign(createSingleShellUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id) + `?version=${item.currentVersionNumber}`)
    }

    transformInputResourceToTableItem(inp : any) : any {
        let obj = {
            id: inp.Id,
            name: inp.Name,
            type: ShellTypes[inp.TypeId],
            versions: null,
            uploadTime: null,
            currentVersion: null,
            currentVersionNumber: -1,
            uploadUser: "Loading...",
            value: inp,
        }

        return obj
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

    selectVersion(obj : any, v : ShellScriptVersion, versionNum : number) {
        obj.currentVersion = v
        obj.currentVersionNumber = versionNum
        obj.uploadTime = standardFormatTime(v.UploadTime)
        obj.uploadUser = createUserString(MetadataStore.getters.getUser(v.UploadUserId))
    }

    shellVersions: Record<number, ShellScriptVersion[]> = Object()
    getShellScriptVersions(shellId : number, obj : any) : ShellScriptVersion[] {
        if (!(shellId in this.shellVersions)) {
            allShellScriptVersions({
                orgId: PageParamsStore.state.organization!.Id,
                shellId: shellId,
            }).then((resp : TAllShellScriptVersionsOutput) => {
                obj.versions = resp.data
                if (obj.versions.length > 0) {
                    this.selectVersion(obj, obj.versions[0], obj.versions.length)
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
        return this.shellVersions[shellId]
    }

    renderVersion(props : any) : VNode {
        if (!props.item.versions) {
            props.item.versions = this.getShellScriptVersions(props.item.id, props.item)
        }

        if (!props.item.currentVersion && props.item.versions.length > 0) {
            this.selectVersion(props.item, props.item.versions[0], props.item.versions.length)
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

        let availableVersionsNodes = !!props.item.versions ? props.item.versions.map((ele : ShellScriptVersion, idx : number) => {
            let versionNum = props.item.versions.length - idx
            let item = this.$createElement(
                VListItem,
                {
                    on: {
                        click: () => {
                            this.selectVersion(props.item, ele, versionNum)
                        }
                    }
                },
                `${versionNum}`
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
                    "disabled": false,
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
                                !!props.item.currentVersion ? props.item.currentVersionNumber : "Loading...",
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

    renderUploadTime(props : any) : VNode {
        return this.$createElement(
            'span',
            !!props.item.uploadTime ? standardFormatTime(props.item.uploadTime) : "Loading..."
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
                    resourceName: "shell scripts",
                    showExpand: true,
                    value: this.value.map((ele : any) => ele.File),
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToShellScript,
                },
                scopedSlots: {
                    'expanded-item': this.renderExpansion,
                    'item.version': this.renderVersion,
                    'item.uploadTime': this.renderUploadTime,
                }
            }
        )
    }

}

</script>
