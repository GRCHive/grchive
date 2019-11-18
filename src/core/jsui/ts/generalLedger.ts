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
        while (catsToProcess.length != 0) {
            let catsHandled = 0

            for (let cat of catsToProcess) {
                if (cat.ParentCategoryId == null || processedCatIds.has(cat.ParentCategoryId)) {
                    this.addRawCategory(cat)
                    processedCatIds.add(cat.Id)
                }
            }

            // TODO: This is an error. How should it be handled?
            if (catsHandled == 0) {
                break
            }
        }

        for (let acc of accs) {
            this.addRawAccount(acc)
        }
    }

    addRawCategory(cat : RawGeneralLedgerCategory) {
        let newCat = <GeneralLedgerCategory>{
            ...cat,
            ParentCategory: null,
            SubCategories: new Map<number, GeneralLedgerCategory>(),
            SubAccounts: new Map<number, GeneralLedgerAccount>(),
            changed: 1
        }

        this.categories.set(cat.Id, newCat)

        if (!!cat.ParentCategoryId) {
            let parentCat = this.categories.get(cat.ParentCategoryId)!
            parentCat.SubCategories.set(cat.Id, newCat)
            newCat.ParentCategory = parentCat
            parentCat.changed += 1
        } else {
            this.topLevelCategories.set(cat.Id, newCat)
        }

        this.changed += 1
    }

    addRawAccount(acc : RawGeneralLedgerAccount) {
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
