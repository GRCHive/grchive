export interface RawGeneralLedgerCategory {
	Id:               number
	OrgId:            number
	ParentCategoryId: number | null
	Name:             string
	Description:      string
}

export interface RawGeneralLedgerAccount {
	Id:                  number
	OrgId:               number
	ParentCategoryId:    number
	AccountId:           string
	AccountName:         string
	AccountDescription:  string
	FinanciallyRelevant: boolean
}

export interface GeneralLedgerAccount extends RawGeneralLedgerAccount {
    ParentCategory: GeneralLedgerCategory
    changed: 1
}

export interface GeneralLedgerCategory extends RawGeneralLedgerCategory {
    SubCategories: Map<number, GeneralLedgerCategory>
    SubAccounts: Map<number, GeneralLedgerAccount>
    ParentCategory: GeneralLedgerCategory | null
    changed: 1
}

export class GeneralLedger {
    topLevelCategories: Map<number, GeneralLedgerCategory>
    categories : Map<number, GeneralLedgerCategory>
    accounts: Map<number, GeneralLedgerAccount>
    changed: 1

    get listCategories() : GeneralLedgerCategory[] {
        return Array.from(this.categories.values())
    }

    get listAccounts() : GeneralLedgerAccount[] {
        return Array.from(this.accounts.values())
    }

    constructor() {
        this.topLevelCategories = new Map<number, GeneralLedgerCategory>()
        this.categories = new Map<number, GeneralLedgerCategory>()
        this.accounts = new Map<number, GeneralLedgerAccount>()
        this.changed = 1
    }

    // This function should be used when we can't assume the parent categories have been
    // added in yet.
    rebuildGL(cats : RawGeneralLedgerCategory[], accs : RawGeneralLedgerAccount[]) {
        // Need to add the categories in order. First, the categories with no parent.
        // Then the categories with those categories as the parent, etc.
        let catsToProcess = [...cats]
        let processedCatIds = new Set<number>()
        console.log(catsToProcess)
        while (catsToProcess.length != 0) {
            let catsHandled = 0

            for (let cat of catsToProcess) {
                if (cat.ParentCategoryId == null || processedCatIds.has(cat.ParentCategoryId)) {
                    this.addRawCategory(cat)
                    processedCatIds.add(cat.Id)
                    catsHandled += 1
                }
            }

            // TODO: This is an error. How should it be handled?
            if (catsHandled == 0) {
                break
            }

            catsToProcess = catsToProcess.filter((ele : RawGeneralLedgerCategory) => !processedCatIds.has(ele.Id))
        }

        for (let acc of accs) {
            this.addRawAccount(acc)
        }
    }

    createCategoryFromRaw(cat : RawGeneralLedgerCategory) : GeneralLedgerCategory {
        let newCat = <GeneralLedgerCategory>{
            ...cat,
            ParentCategory: null,
            SubCategories: new Map<number, GeneralLedgerCategory>(),
            SubAccounts: new Map<number, GeneralLedgerAccount>(),
            changed: 1
        }


        if (!!cat.ParentCategoryId) {
            let parentCat = this.categories.get(cat.ParentCategoryId)!
            newCat.ParentCategory = parentCat
        }

        return newCat
    }

    addRawCategory(cat : RawGeneralLedgerCategory) {
        let newCat = this.createCategoryFromRaw(cat)
        this.categories.set(cat.Id, newCat)
        console.log("add cat ", cat.Id)

        if (!!cat.ParentCategoryId) {
            let parentCat = this.categories.get(cat.ParentCategoryId)!
            parentCat.SubCategories.set(cat.Id, newCat)
            parentCat.changed += 1
        } else {
            this.topLevelCategories.set(cat.Id, newCat)
        }

        this.changed += 1
    }

    replaceRawCategory(cat : RawGeneralLedgerCategory) {
        if (!this.categories.has(cat.Id)) {
            return
        }
    
        let existingCat = this.categories.get(cat.Id)!

        // Remove connections from parent. Children pointer should still be OK!
        if (!!existingCat.ParentCategoryId) {
            existingCat.ParentCategory!.SubCategories.delete(existingCat.Id)
            existingCat.ParentCategory!.changed += 1
        } else {
            this.topLevelCategories.delete(existingCat.Id)
        }

        let newCat = this.createCategoryFromRaw(cat)
        existingCat.Name = newCat.Name
        existingCat.Description = newCat.Description
        existingCat.ParentCategory = newCat.ParentCategory
        existingCat.ParentCategoryId = newCat.ParentCategoryId
        existingCat.changed += 1

        // Reconnect parent.
        if (!!existingCat.ParentCategoryId) {
            existingCat.ParentCategory!.SubCategories.set(existingCat.Id, existingCat)
            existingCat.ParentCategory!.changed += 1
        } else {
            this.topLevelCategories.set(existingCat.Id, existingCat)
        }

        this.changed += 1
    }

    removeCategory(catId : number) {
        if (!this.categories.has(catId)) {
            return
        }
    
        let existingCat = this.categories.get(catId)!
        this.categories.delete(catId)
        this.topLevelCategories.delete(catId)

        if (!!existingCat.ParentCategory) {
            existingCat.ParentCategory!.SubCategories.delete(existingCat.Id)
        }

        for (let subCat of existingCat.SubCategories.values()) {
            this.removeCategory(subCat.Id)
        }

        for (let subAcc of existingCat.SubAccounts.values()) {
            this.removeAccount(subAcc.Id)
        }

        this.changed += 1
    }

    removeAccount(accountId : number) {
        if (!this.accounts.has(accountId)) {
            return
        }

        let acc = this.accounts.get(accountId)!
        acc.ParentCategory.SubAccounts.delete(accountId)
        this.accounts.delete(accountId)
        this.changed += 1
    }

    addRawAccount(acc : RawGeneralLedgerAccount) {
        console.log("parent cat ", acc.ParentCategoryId)
        let parentCat = this.categories.get(acc.ParentCategoryId)!
        let newAcc = <GeneralLedgerAccount>{
            ...acc,
            ParentCategory: parentCat,
            changed: 1
        }

        this.accounts.set(acc.Id, newAcc)
        parentCat.SubAccounts.set(acc.Id, newAcc)

        parentCat.changed += 1
        this.changed += 1
    }
}
