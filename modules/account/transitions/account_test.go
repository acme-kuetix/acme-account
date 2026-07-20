package transitions

import (
	"strconv"
	"sync/atomic"
	"testing"

	coreutils "github.com/acme-kuetix/acme-std-core/modules/core/utils/transitions"
	persistenceStore "github.com/acme-kuetix/acme-std-persistence/modules/persistence/store/transitions"
)

// testIDCounter generates unique IDs for test seeding. Production code
// uses sequence/sequence.NextByCode(code: "account") via WSL.
var testIDCounter int64

func nextTestID() string {
	return strconv.FormatInt(atomic.AddInt64(&testIDCounter, 1), 10)
}

func newAccount() *accountTransitions {
	return NewAccountTransitions().(*accountTransitions)
}

// resetAll clears both the in-memory ID counter/code index AND the
// persistence/store data.
func resetAll() {
	persistenceStore.ResetStore()
	atomic.StoreInt64(&testIDCounter, 0)
}

// seedAccount mirrors the account-create WSL workflow: sequence.NextByCode,
// persistence/store.Create. Returns the new account id.
// Used by struct-method tests that need pre-existing accounts in the store.
func seedAccount(code, name, typ, parentCode, debitCredit string, active bool) string {
	id := nextTestID()
	now := "now"
	doc := map[string]interface{}{
		"id":          id,
		"code":        code,
		"name":        name,
		"type":        typ,
		"parentCode":  parentCode,
		"debitCredit": debitCredit,
		"active":      active,
		"createdAt":   now,
		"updatedAt":   now,
	}
	store := persistenceStore.NewStoreTransitionsConcrete()
	store.Create(accountsCollection, id, doc)
	return id
}

// accountByCode scans persistence/store for the doc with the given code.
// Returns nil if not found.
func accountByCode(code string) map[string]interface{} {
	store := persistenceStore.NewStoreTransitionsConcrete()
	listRes := store.List(accountsCollection, map[string]interface{}{})
	if !listRes.Success {
		return nil
	}
	for _, d := range coreutils.DocsFromListVal(listRes.Response) {
		if c, _ := d["code"].(string); c == code {
			return d
		}
	}
	return nil
}

// ─── BuildUpdateDoc — migrated to WSL via persistence/store/store.BuildHierarchicalUpdate ───
// The Go BuildUpdateDoc method was deleted; the update validation + merge now
// runs in WSL via persistence/store/store.BuildHierarchicalUpdate (generic
// primitive reusable for accounts, stock locations, etc.). Coverage lives
// in acme-std-persistence TestStoreBuildHierarchicalUpdate* + acme-app-erp
// TestAccountUpdateViaWSL integration tests.

func TestBuildUpdateDocMigratedToWSL(t *testing.T) {
	t.Skip("BuildUpdateDoc migrated to WSL — see persistence/store/store.BuildHierarchicalUpdate")
}

func TestBuildUpdateDocTypeChangeWithChildrenFailsMigratedToWSL(t *testing.T) {
	t.Skip("BuildUpdateDoc migrated to WSL — see persistence/store/store.BuildHierarchicalUpdate")
}

// ─── GetAccountTree (migrated to WSL via collections/ops/ops.RecurseForest) ───

func TestGetAccountTree(t *testing.T) {
	t.Skip("GetAccountTree migrated to WSL — tree now built by collections/ops/ops.RecurseForest")
}

func TestGetAccountTreeWithRootId(t *testing.T) {
	t.Skip("GetAccountTree migrated to WSL — tree now built by collections/ops/ops.RecurseForest")
}

// ─── NextAccountId (now via sequence/sequence.NextByCode in WSL) ───
// The Go NextAccountId method was deleted; ID generation now flows
// through sequence/sequence.NextByCode(code: "account", scope: "") in
// the account-create WSL workflow. The acme-app-erp integration test
// TestAccount_CreateViaWSL verifies the end-to-end flow. This test
// only verifies the test-only ID counter used for seeding.
func TestNextAccountId(t *testing.T) {
	id1 := nextTestID()
	id2 := nextTestID()
	if id1 == id2 {
		t.Errorf("expected distinct IDs, got %q twice", id1)
	}
}
